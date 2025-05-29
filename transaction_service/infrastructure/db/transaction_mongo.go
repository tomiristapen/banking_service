package db

import (
	"context"
	"time"

	"transaction_service/domain/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type TransactionMongo struct {
	collection *mongo.Collection
}

func NewTransactionMongo(db *mongo.Database) *TransactionMongo {
	return &TransactionMongo{
		collection: db.Collection("transactions"),
	}
}

func (r *TransactionMongo) Create(tx *model.Transaction) error {
	tx.CreatedAt = time.Now().Format(time.RFC3339)
	res, err := r.collection.InsertOne(context.Background(), tx)
	if err != nil {
		return err
	}
	if oid, ok := res.InsertedID.(interface{ Hex() string }); ok {
		tx.ID = oid.Hex()
	}
	return nil
}

func (r *TransactionMongo) GetHistory(userID string) ([]*model.Transaction, error) {
	filter := bson.M{"$or": []bson.M{{"from_user_id": userID}, {"to_user_id": userID}}}
	cur, err := r.collection.Find(context.Background(), filter)
	if err != nil {
		return nil, err
	}
	defer cur.Close(context.Background())
	var txs []*model.Transaction
	for cur.Next(context.Background()) {
		var tx model.Transaction
		if err := cur.Decode(&tx); err != nil {
			continue
		}
		txs = append(txs, &tx)
	}
	return txs, nil
}

func (r *TransactionMongo) GetStatus(transferID string) (string, error) {
	objID, err := primitive.ObjectIDFromHex(transferID)
	if err != nil {
		return "", err
	}
	var tx model.Transaction
	err = r.collection.FindOne(context.Background(), bson.M{"_id": objID}).Decode(&tx)
	if err != nil {
		return "", err
	}
	return tx.Status, nil
}
