package db

import (
	"context"
	"log"

	"github.com/tomiristapen/banking_service/smart_budget_service/domain/model"
	"github.com/tomiristapen/banking_service/smart_budget_service/domain/repository"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type BudgetMongoRepo struct {
	col *mongo.Collection
}

func NewBudgetMongoRepo(col *mongo.Collection) repository.BudgetRepository {
	return &BudgetMongoRepo{col: col}
}

// mongoCategoryExpense используется только для хранения в MongoDB
// и не "засоряет" domain/model тегами
type mongoCategoryExpense struct {
	PaymentID string  `bson:"payment_id"`
	UserID    string  `bson:"user_id"`
	Category  string  `bson:"category"`
	Amount    float64 `bson:"amount"`
	Service   string  `bson:"service"`
	Status    string  `bson:"status"`
	CreatedAt string  `bson:"created_at"`
}

func toMongo(exp *model.CategoryExpense) *mongoCategoryExpense {
	return &mongoCategoryExpense{
		PaymentID: exp.PaymentID,
		UserID:    exp.UserID,
		Category:  exp.Category,
		Amount:    exp.Amount,
		Service:   exp.Service,
		Status:    exp.Status,
		CreatedAt: exp.CreatedAt,
	}
}

func fromMongo(m *mongoCategoryExpense) *model.CategoryExpense {
	return &model.CategoryExpense{
		PaymentID: m.PaymentID,
		UserID:    m.UserID,
		Category:  m.Category,
		Amount:    m.Amount,
		Service:   m.Service,
		Status:    m.Status,
		CreatedAt: m.CreatedAt,
	}
}

func (r *BudgetMongoRepo) AddExpense(ctx context.Context, expense *model.CategoryExpense) error {
	log.Printf("[SmartBudget][Mongo][DEBUG] AddExpense: type=%T, value=%+v", expense, expense)
	mongoExp := toMongo(expense)
	res, err := r.col.InsertOne(ctx, mongoExp)
	if err != nil {
		log.Printf("[SmartBudget][Mongo][ERROR] InsertOne failed: %v", err)
		return err
	}
	log.Printf("[SmartBudget][Mongo][DEBUG] InsertOne result: insertedID=%v", res.InsertedID)
	return nil
}

func (r *BudgetMongoRepo) GetBudgetSummary(ctx context.Context, userID string) ([]*model.BudgetSummary, error) {
	pipeline := mongo.Pipeline{
		bson.D{{Key: "$match", Value: bson.D{{Key: "user_id", Value: userID}}}},
		bson.D{{Key: "$group", Value: bson.D{{Key: "_id", Value: "$category"}, {Key: "total", Value: bson.D{{Key: "$sum", Value: "$amount"}}}}}},
	}
	cur, err := r.col.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	var res []*model.BudgetSummary
	for cur.Next(ctx) {
		var doc struct {
			ID    string  `bson:"_id"`
			Total float64 `bson:"total"`
		}
		if err := cur.Decode(&doc); err != nil {
			continue
		}
		res = append(res, &model.BudgetSummary{UserID: userID, Category: doc.ID, Total: doc.Total})
	}
	return res, nil
}

func (r *BudgetMongoRepo) GetCategoryExpenses(ctx context.Context, userID, category string) ([]*model.CategoryExpense, error) {
	filter := bson.M{"user_id": userID, "category": category}
	cur, err := r.col.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	var res []*model.CategoryExpense
	for cur.Next(ctx) {
		var mExp mongoCategoryExpense
		if err := cur.Decode(&mExp); err != nil {
			continue
		}
		res = append(res, fromMongo(&mExp))
	}
	return res, nil
}
