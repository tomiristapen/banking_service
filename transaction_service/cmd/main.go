package main

import (
	"context"
	"log"
	"net"
	"os"
	grpcadapter "transaction_service/adapter/grpc"
	"transaction_service/infrastructure/db"
	"transaction_service/infrastructure/mq"
	transactionpb "transaction_service/proto"
	"transaction_service/usecase"

	"github.com/nats-io/nats.go"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
)

func main() {
	// MongoDB setup
	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		mongoURI = "mongodb://localhost:27017"
	}
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatal(err)
	}
	dbConn := client.Database("banking")
	txRepo := db.NewTransactionMongo(dbConn)

	// NATS setup
	natsURL := os.Getenv("NATS_URL")
	if natsURL == "" {
		natsURL = nats.DefaultURL
	}
	nc, err := nats.Connect(natsURL)
	if err != nil {
		log.Fatal(err)
	}
	publisher := mq.NewNatsPublisher(nc, "transactions")

	// Usecase
	uc := usecase.NewTransactionUsecase(txRepo, publisher)

	// gRPC server
	lis, err := net.Listen("tcp", ":50053")
	if err != nil {
		log.Fatal(err)
	}
	grpcServer := grpc.NewServer()
	transactionpb.RegisterTransactionServiceServer(grpcServer, grpcadapter.NewTransactionHandler(uc))
	log.Println("TransactionService gRPC server started on :50053")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatal(err)
	}
}
