package grpcclient

import (
	context "context"

	userpb "github.com/tomiristapen/banking_service/user_service/proto"
	"google.golang.org/grpc"
)

type UserServiceClient struct {
	client userpb.UserServiceClient
}

func NewUserServiceClient(conn *grpc.ClientConn) *UserServiceClient {
	return &UserServiceClient{client: userpb.NewUserServiceClient(conn)}
}

func (u *UserServiceClient) GetBalance(ctx context.Context, userID string) (float64, error) {
	resp, err := u.client.GetBalance(ctx, &userpb.GetBalanceRequest{UserId: userID})
	if err != nil {
		return 0, err
	}
	return resp.Balance, nil
}

func (u *UserServiceClient) DecreaseBalance(ctx context.Context, userID string, amount float64) (bool, string, error) {
	resp, err := u.client.DecreaseBalance(ctx, &userpb.DecreaseBalanceRequest{UserId: userID, Amount: amount})
	if err != nil {
		return false, "", err
	}
	return resp.Success, resp.Message, nil
}
