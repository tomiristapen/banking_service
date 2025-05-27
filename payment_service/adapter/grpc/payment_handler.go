package grpc

import (
	"context"

	"github.com/tomiristapen/banking_service/payment_service/usecase"
	pb "github.com/tomiristapen/banking_service/payment_service/proto"
)

type PaymentHandler struct {
	pb.UnimplementedPaymentServiceServer
	uc *usecase.PaymentUsecase
}

func NewPaymentHandler(uc *usecase.PaymentUsecase) *PaymentHandler {
	return &PaymentHandler{uc: uc}
}

func (h *PaymentHandler) PayForService(ctx context.Context, req *pb.PaymentRequest) (*pb.PaymentResponse, error) {
	payment, err := h.uc.PayForService(ctx, req.UserId, req.Service, req.Amount, req.Category)
	if err != nil {
		return &pb.PaymentResponse{
			PaymentId: "",
			Status:    "failed",
			Message:   err.Error(),
		}, nil
	}

	return &pb.PaymentResponse{
		PaymentId: payment.ID,
		Status:    payment.Status,
		Message:   "Payment successful",
	}, nil
}

func (h *PaymentHandler) ListAvailableServices(ctx context.Context, _ *pb.ListServicesRequest) (*pb.ListServicesResponse, error) {
	services, err := h.uc.ListAvailableServices(ctx)
	if err != nil {
		return nil, err
	}
	return &pb.ListServicesResponse{Services: services}, nil
}

func (h *PaymentHandler) GetPaymentStatus(ctx context.Context, req *pb.PaymentStatusRequest) (*pb.PaymentStatusResponse, error) {
	p, err := h.uc.GetPaymentStatus(ctx, req.PaymentId)
	if err != nil {
		return &pb.PaymentStatusResponse{
			Status: "not_found",
		}, nil
	}
	return &pb.PaymentStatusResponse{
		Status: p.Status,
	}, nil
}
