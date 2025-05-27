package usecase

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/patrickmn/go-cache"
	log "github.com/sirupsen/logrus"
	"github.com/tomiristapen/banking_service/user_service/domain/model"
	"github.com/tomiristapen/banking_service/user_service/domain/repository"
	"github.com/tomiristapen/banking_service/user_service/infrastructure/smtp"
	"github.com/tomiristapen/banking_service/user_service/metrics"
	"github.com/tomiristapen/banking_service/user_service/utils"
)

type UserUseCase struct {
	userRepo repository.UserRepository
	mailer   smtp.Mailer
	codes    map[string]string
	mu       sync.Mutex

	// TTL-кеш для профилей и баланса
	profileCache *cache.Cache // userID -> *model.User
}

func NewUserUseCase(userRepo repository.UserRepository, mailer smtp.Mailer) *UserUseCase {
	return &UserUseCase{
		userRepo:     userRepo,
		mailer:       mailer,
		codes:        make(map[string]string),
		profileCache: cache.New(5*time.Minute, 10*time.Minute),
	}
}

func (u *UserUseCase) Register(ctx context.Context, name, email, password string) (*model.User, error) {
	existingUser, _ := u.userRepo.FindByEmail(ctx, email)
	if existingUser != nil {
		return nil, errors.New("user already exists")
	}

	hashed, err := utils.HashPassword(password)
	if err != nil {
		return nil, err
	}

	user := &model.User{
		ID:         utils.GenerateID(),
		Name:       name,
		Email:      email,
		Password:   hashed,
		CreatedAt:  time.Now(),
		IsVerified: false,
		Balance:    10000, // ✅ стартовый баланс
	}

	savedUser, err := u.userRepo.CreateUser(ctx, user)
	if err != nil {
		return nil, err
	}

	code := utils.GenerateVerificationCode()

	u.mu.Lock()
	u.codes[savedUser.ID] = code
	u.mu.Unlock()

	metrics.RegisterSuccess.Inc()
	log.WithFields(log.Fields{
		"user_id":           savedUser.ID,
		"email":             savedUser.Email,
		"verification_code": code, // <-- логируем verification_code
	}).Info("User registered")

	_ = u.mailer.SendVerificationEmail(savedUser.Email, code)

	return savedUser, nil
}

func (u *UserUseCase) Login(ctx context.Context, email, password string) (string, error) {
	user, err := u.userRepo.FindByEmail(ctx, email)
	if err != nil || user == nil {
		metrics.LoginFailed.Inc()
		log.WithField("email", email).Warn("Login failed: user not found")
		return "", errors.New("user not found")
	}

	if err := utils.CheckPassword(password, user.Password); err != nil {
		metrics.LoginFailed.Inc()
		log.WithField("email", email).Warn("Login failed: invalid password")
		return "", errors.New("invalid password")
	}

	if !user.IsVerified {
		metrics.LoginFailed.Inc()
		log.WithField("email", email).Warn("Login failed: email not verified")
		return "", errors.New("email not verified")
	}

	token, err := utils.GenerateJWT(user.ID, user.Email)
	if err != nil {
		log.WithError(err).Error("Failed to generate token")
		return "", err
	}

	metrics.LoginSuccess.Inc()
	log.WithField("email", email).Info("Login successful")

	return token, nil
}

func (u *UserUseCase) VerifyEmail(ctx context.Context, userID, code string) (bool, error) {
	u.mu.Lock()
	expectedCode, ok := u.codes[userID]
	log.WithFields(log.Fields{
		"user_id":       userID,
		"input_code":    code,
		"expected_code": expectedCode,
		"found":         ok,
	}).Info("[DEBUG] VerifyEmail input and codes map state")
	u.mu.Unlock()

	if !ok || expectedCode != code {
		return false, errors.New("invalid or expired verification code")
	}

	user, err := u.userRepo.FindByID(ctx, userID)
	if err != nil {
		return false, errors.New("user not found")
	}

	user.IsVerified = true
	if err := u.userRepo.UpdateVerificationStatus(ctx, userID, true); err != nil {
		return false, err
	}

	u.mu.Lock()
	delete(u.codes, userID)
	u.mu.Unlock()

	return true, nil
}

func (u *UserUseCase) GetProfile(ctx context.Context, userID string) (*model.User, error) {
	if cached, found := u.profileCache.Get(userID); found {
		log.Infof("[CACHE] GetProfile: userID=%s (from cache)", userID)
		return cached.(*model.User), nil
	}
	log.Infof("[DB] GetProfile: userID=%s (from DB)", userID)
	user, err := u.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	u.profileCache.Set(userID, user, cache.DefaultExpiration)
	return user, nil
}

// ✅ Balance functional methods

func (u *UserUseCase) GetBalance(ctx context.Context, userID string) (float64, error) {
	if cached, found := u.profileCache.Get(userID); found {
		log.Infof("[CACHE] GetBalance: userID=%s (from cache)", userID)
		return cached.(*model.User).Balance, nil
	}
	log.Infof("[DB] GetBalance: userID=%s (from DB)", userID)
	user, err := u.userRepo.FindByID(ctx, userID)
	if err != nil {
		return 0, errors.New("user not found")
	}
	u.profileCache.Set(userID, user, cache.DefaultExpiration)
	return user.Balance, nil
}

func (u *UserUseCase) DecreaseBalance(ctx context.Context, userID string, amount float64) error {
	user, err := u.userRepo.FindByID(ctx, userID)
	if err != nil {
		return errors.New("user not found")
	}
	if user.Balance < amount {
		return errors.New("insufficient funds")
	}
	user.Balance -= amount
	err = u.userRepo.UpdateUser(ctx, user)
	if err == nil {
		u.profileCache.Set(userID, user, cache.DefaultExpiration)
		log.Infof("[CACHE] DecreaseBalance: userID=%s (cache updated)", userID)
	}
	return err
}
