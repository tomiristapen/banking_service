package usecase

import (
    "context"
    "errors"
    "testing"

    "github.com/tomiristapen/banking_service/user_service/domain/model"
    "github.com/tomiristapen/banking_service/user_service/utils"
)

// Фейковый репо
type fakeRepo struct {
    users map[string]*model.User
}

func (f *fakeRepo) CreateUser(ctx context.Context, user *model.User) (*model.User, error) {
    f.users[user.Email] = user
    return user, nil
}

func (f *fakeRepo) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
    if u, ok := f.users[email]; ok {
        return u, nil
    }
    return nil, errors.New("not found")
}

func (f *fakeRepo) FindByEmail(ctx context.Context, email string) (*model.User, error) {
    return f.GetUserByEmail(ctx, email)
}

func (f *fakeRepo) GetUserByID(ctx context.Context, id string) (*model.User, error) {
    for _, u := range f.users {
        if u.ID == id {
            return u, nil
        }
    }
    return nil, errors.New("not found")
}

func (f *fakeRepo) FindByID(ctx context.Context, id string) (*model.User, error) {
    return f.GetUserByID(ctx, id)
}

func (f *fakeRepo) UpdateUser(ctx context.Context, user *model.User) error {
    f.users[user.Email] = user
    return nil
}

func (f *fakeRepo) UpdateVerificationStatus(ctx context.Context, id string, status bool) error {
    for _, u := range f.users {
        if u.ID == id {
            u.IsVerified = status
            return nil
        }
    }
    return errors.New("not found")
}

// фейковый мейлер
type fakeMailer struct {
    SentTo string
    Code   string
}

func (f *fakeMailer) SendVerificationEmail(email, code string) error {
    f.SentTo = email
    f.Code = code
    return nil
}

// тест регистрации
func TestRegisterUser(t *testing.T) {
    repo := &fakeRepo{users: make(map[string]*model.User)}
    mailer := &fakeMailer{}
    uc := NewUserUseCase(repo, mailer)

    name := "Test"
    email := "test@example.com"
    password := "123456"

    user, err := uc.Register(context.Background(), name, email, password)
    if err != nil {
        t.Fatalf("expected no error, got %v", err)
    }

    if user.Email != email {
        t.Errorf("expected email %s, got %s", email, user.Email)
    }

    if mailer.SentTo != email {
        t.Errorf("expected verification email to %s, got %s", email, mailer.SentTo)
    }

    if mailer.Code == "" {
        t.Errorf("expected verification code to be sent, got empty string")
    }
}

// тест логина 
func TestLoginUser(t *testing.T) {
    repo := &fakeRepo{users: make(map[string]*model.User)}
    mailer := &fakeMailer{}
    uc := NewUserUseCase(repo, mailer)

    email := "test@example.com"
    password := "123456"
    _, _ = uc.Register(context.Background(), "Test", email, password)

    repo.users[email].IsVerified = true

    token, err := uc.Login(context.Background(), email, password)
    if err != nil {
        t.Fatalf("expected login to succeed, got error: %v", err)
    }

    if token == "" {
        t.Errorf("expected non-empty token")
    }
}

// Тест VerifyEmail
func TestVerifyEmail(t *testing.T) {
    repo := &fakeRepo{users: make(map[string]*model.User)}
    mailer := &fakeMailer{}
    uc := NewUserUseCase(repo, mailer)

    email := "test@example.com"
    password := "123456"
    user, _ := uc.Register(context.Background(), "Test", email, password)

    code := mailer.Code
    ok, err := uc.VerifyEmail(context.Background(), user.ID, code)
    if err != nil {
        t.Fatalf("expected no error, got %v", err)
    }
    if !ok {
        t.Errorf("expected verification success")
    }
    if !repo.users[email].IsVerified {
        t.Errorf("expected user to be marked as verified")
    }
}

// тест GetUserProfile
func TestGetUserProfile(t *testing.T) {
    repo := &fakeRepo{users: make(map[string]*model.User)}
    mailer := &fakeMailer{}
    uc := NewUserUseCase(repo, mailer)

    email := "test@example.com"
    password := "123456"
    user, _ := uc.Register(context.Background(), "Test", email, password)

    repo.users[email].IsVerified = true

    token, _ := uc.Login(context.Background(), email, password)
    userID, _ := utils.ParseToken(token)

    got, err := uc.GetProfile(context.Background(), userID)
    if err != nil {
        t.Fatalf("expected no error, got %v", err)
    }
    if got.ID != user.ID {
        t.Errorf("expected user ID %s, got %s", user.ID, got.ID)
    }
}

// тест повторной регистрации
func TestRegisterConflict(t *testing.T) {
    repo := &fakeRepo{users: make(map[string]*model.User)}
    mailer := &fakeMailer{}
    uc := NewUserUseCase(repo, mailer)

    email := "test@example.com"
    password := "123456"
    _, _ = uc.Register(context.Background(), "Test", email, password)

    _, err := uc.Register(context.Background(), "Test", email, password)
    if err == nil {
        t.Fatalf("expected error due to duplicate email, got none")
    }
}

// тест логина не верифицированного 
func TestLoginUnverifiedUser(t *testing.T) {
    repo := &fakeRepo{users: make(map[string]*model.User)}
    mailer := &fakeMailer{}
    uc := NewUserUseCase(repo, mailer)

    email := "test@example.com"
    password := "123456"
    _, _ = uc.Register(context.Background(), "Test", email, password)

    _, err := uc.Login(context.Background(), email, password)
    if err == nil {
        t.Fatalf("expected error due to unverified email, got none")
    }
}
