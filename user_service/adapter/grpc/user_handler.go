package grpc

import (
	"context"
	"errors"
	"os"

	userpb "github.com/tomiristapen/banking_service/user_service/proto"
	"github.com/tomiristapen/banking_service/user_service/usecase"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/golang-jwt/jwt/v4"
)

// Load JWT secret from environment variable
var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

type UserHandler struct {
	userpb.UnimplementedUserServiceServer
	uc *usecase.UserUseCase
}

func NewUserHandler(uc *usecase.UserUseCase) *UserHandler {
	return &UserHandler{
		uc: uc,
		UnimplementedUserServiceServer: userpb.UnimplementedUserServiceServer{},
	}
}

// Parse userId from JWT token
func parseUserIdFromToken(tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Check signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return jwtSecret, nil
	})

	if err != nil {
		return "", err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if userId, ok := claims["user_id"].(string); ok {
			return userId, nil
		}
		return "", errors.New("user_id not found in token")
	}
	return "", errors.New("invalid token claims")
}

func (h *UserHandler) RegisterUser(ctx context.Context, req *userpb.RegisterRequest) (*userpb.RegisterResponse, error) {
	user, err := h.uc.Register(ctx, req.Name, req.Email, req.Password)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "registration failed: %v", err)
	}
	return &userpb.RegisterResponse{
		UserId:  user.ID,
		Message: "registration successful",
	}, nil
}

func (h *UserHandler) LoginUser(ctx context.Context, req *userpb.LoginRequest) (*userpb.LoginResponse, error) {
	token, err := h.uc.Login(ctx, req.Email, req.Password)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "login failed: %v", err)
	}
	return &userpb.LoginResponse{
		Token:   token,
		Message: "login successful",
	}, nil
}

func (h *UserHandler) VerifyEmail(ctx context.Context, req *userpb.VerifyEmailRequest) (*userpb.VerifyEmailResponse, error) {
	success, err := h.uc.VerifyEmail(ctx, req.UserId, req.Code)
	if err != nil {
		return &userpb.VerifyEmailResponse{Success: false, Message: err.Error()}, nil
	}
	return &userpb.VerifyEmailResponse{Success: success, Message: "email verified"}, nil
}

// Use token from request to get user profile
func (h *UserHandler) GetUserProfile(ctx context.Context, req *userpb.ProfileRequest) (*userpb.ProfileResponse, error) {
	userId, err := parseUserIdFromToken(req.Token)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "invalid token: %v", err)
	}

	user, err := h.uc.GetProfile(ctx, userId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "get profile failed: %v", err)
	}

	return &userpb.ProfileResponse{
		Id:         user.ID,
		Name:       user.Name,
		Email:      user.Email,
		IsVerified: user.IsVerified,
		Balance:    user.Balance,
	}, nil
}

func (h *UserHandler) GetBalance(ctx context.Context, req *userpb.GetBalanceRequest) (*userpb.GetBalanceResponse, error) {
	balance, err := h.uc.GetBalance(ctx, req.UserId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "get balance failed: %v", err)
	}
	return &userpb.GetBalanceResponse{Balance: balance}, nil
}

func (h *UserHandler) DecreaseBalance(ctx context.Context, req *userpb.DecreaseBalanceRequest) (*userpb.BalanceUpdateResponse, error) {
	err := h.uc.DecreaseBalance(ctx, req.UserId, req.Amount)
	if err != nil {
		return &userpb.BalanceUpdateResponse{Success: false, Message: err.Error()}, nil
	}
	return &userpb.BalanceUpdateResponse{Success: true, Message: "balance updated successfully"}, nil
}
