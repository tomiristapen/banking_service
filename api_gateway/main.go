package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"

	paymentpb "github.com/tomiristapen/banking_service/api_gateway/proto/payment"
	smartbudgetpb "github.com/tomiristapen/banking_service/api_gateway/proto/smartbudget"
	transactionpb "github.com/tomiristapen/banking_service/api_gateway/proto/transaction"
	userpb "github.com/tomiristapen/banking_service/api_gateway/proto/user"
)

func main() {
	_ = godotenv.Load()

	r := mux.NewRouter()

	// gRPC-Gateway mux
	gwmux := runtime.NewServeMux()

	ctx := context.Background()

	userAddr := os.Getenv("USER_SERVICE_ADDR")
	transactionAddr := os.Getenv("TRANSACTION_SERVICE_ADDR")
	paymentAddr := os.Getenv("PAYMENT_SERVICE_ADDR")
	smartBudgetAddr := os.Getenv("SMART_BUDGET_SERVICE_ADDR")

	log.Println("User service addr:", userAddr)
	log.Println("Transaction service addr:", transactionAddr)
	log.Println("Payment service addr:", paymentAddr)
	log.Println("Smart budget service addr:", smartBudgetAddr)

	opts := []grpc.DialOption{grpc.WithInsecure()}

	if err := userpb.RegisterUserServiceHandlerFromEndpoint(ctx, gwmux, userAddr, opts); err != nil {
		log.Fatalf("Failed to register user service: %v", err)
	}
	if err := transactionpb.RegisterTransactionServiceHandlerFromEndpoint(ctx, gwmux, transactionAddr, opts); err != nil {
		log.Fatalf("Failed to register transaction service: %v", err)
	}
	if err := paymentpb.RegisterPaymentServiceHandlerFromEndpoint(ctx, gwmux, paymentAddr, opts); err != nil {
		log.Fatalf("Failed to register payment service: %v", err)
	}
	if err := smartbudgetpb.RegisterSmartBudgetServiceHandlerFromEndpoint(ctx, gwmux, smartBudgetAddr, opts); err != nil {
		log.Fatalf("Failed to register smart budget service: %v", err)
	}

	// Mount gRPC-Gateway mux at root
	r.PathPrefix("/").Handler(gwmux)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("API Gateway listening on :%s", port)
	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
