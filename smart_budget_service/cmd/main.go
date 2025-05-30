package main

import (
	"context"
	"log"
	"net"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"

	grpcadapter "github.com/tomiristapen/banking_service/smart_budget_service/adapter/grpc"
	"github.com/tomiristapen/banking_service/smart_budget_service/domain/model"
	"github.com/tomiristapen/banking_service/smart_budget_service/infrastructure/db"
	"github.com/tomiristapen/banking_service/smart_budget_service/infrastructure/mq"
	smartbudgetpb "github.com/tomiristapen/banking_service/smart_budget_service/proto"
	"github.com/tomiristapen/banking_service/smart_budget_service/usecase"
)

func main() {
	_ = godotenv.Load()

	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		mongoURI = "mongodb://localhost:27017"
	}
	log.Printf("[SmartBudget] MONGO_URI: %s", mongoURI)
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	if err := client.Ping(context.TODO(), nil); err != nil {
		log.Fatalf("[SmartBudget] MongoDB ping failed: %v", err)
	}
	log.Println("[SmartBudget] Successfully connected and pinged MongoDB")
	col := client.Database("smart_budgeting_service").Collection("budget_expenses")
	repo := db.NewBudgetMongoRepo(col)
	uc := usecase.NewBudgetUsecase(repo)

	
	natsURL := os.Getenv("NATS_URL")
	if natsURL == "" {
		natsURL = "nats://localhost:4222"
	}
	subscriber := mq.NewMQSubscriber(natsURL, repo)
	go subscriber.SubscribePayments()

	
	txHandler := func(ctx context.Context, event *mq.TransactionEvent) {
		exp := &model.CategoryExpense{
			PaymentID: event.ID,
			UserID:    event.FromUserID,
			Category:  "transfer", 
			Amount:    event.Amount,
			Service:   "transaction_service",
			Status:    event.Status,
			CreatedAt: event.CreatedAt,
		}
		log.Printf("[SmartBudget][Transaction] AddExpense: %+v", exp)
		if err := repo.AddExpense(ctx, exp); err != nil {
			log.Printf("[SmartBudget][Transaction] Failed to add expense: %v", err)
		}
	}
	txSubscriber := mq.NewTransactionSubscriber(natsURL, txHandler)
	go txSubscriber.SubscribeTransactions()

	
	port := os.Getenv("PORT")
	if port == "" {
		port = "50053"
	}
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	handler := grpcadapter.NewBudgetHandler(uc)
	server := grpc.NewServer()
	smartbudgetpb.RegisterSmartBudgetServiceServer(server, handler)
	log.Printf("✅ SmartBudgetService is running on port %s...", port)
	if err := server.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
