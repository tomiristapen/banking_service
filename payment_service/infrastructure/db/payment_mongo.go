package db

import (
	"context"
	"errors"
	"time"

	"github.com/tomiristapen/banking_service/payment_service/domain/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// PaymentMongoRepo implements the PaymentRepository interface
type PaymentMongoRepo struct {
	paymentCol     *mongo.Collection
	serviceListCol *mongo.Collection
}

func NewPaymentMongoRepo(paymentCol, serviceListCol *mongo.Collection) *PaymentMongoRepo {
	return &PaymentMongoRepo{
		paymentCol:     paymentCol,
		serviceListCol: serviceListCol,
	}
}

// internal Mongo document struct
type paymentDoc struct {
	ID        string    `bson:"_id"`
	UserID    string    `bson:"user_id"`
	Amount    float64   `bson:"amount"`
	Service   string    `bson:"service"`
	Category  string    `bson:"category"`
	Type      string    `bson:"type"`
	Status    string    `bson:"status"`
	CreatedAt time.Time `bson:"created_at"`
}

// helpers for mapping between domain and db

func toDoc(p *model.Payment) *paymentDoc {
	return &paymentDoc{
		ID:        p.ID,
		UserID:    p.UserID,
		Amount:    p.Amount,
		Service:   p.Service,
		Category:  p.Category,
		Type:      p.Type,
		Status:    p.Status,
		CreatedAt: p.CreatedAt,
	}
}

func fromDoc(d *paymentDoc) *model.Payment {
	return &model.Payment{
		ID:        d.ID,
		UserID:    d.UserID,
		Amount:    d.Amount,
		Service:   d.Service,
		Category:  d.Category,
		Type:      d.Type,
		Status:    d.Status,
		CreatedAt: d.CreatedAt,
	}
}

// Pay inserts a payment record into Mongo
func (r *PaymentMongoRepo) Pay(ctx context.Context, p *model.Payment) error {
	doc := toDoc(p)
	_, err := r.paymentCol.InsertOne(ctx, doc)
	return err
}

// GetStatus finds a payment by ID
func (r *PaymentMongoRepo) GetStatus(ctx context.Context, id string) (*model.Payment, error) {
	var doc paymentDoc
	err := r.paymentCol.FindOne(ctx, bson.M{"_id": id}).Decode(&doc)
	if err != nil {
		return nil, err
	}
	return fromDoc(&doc), nil
}

// ListAvailableServices returns all service names from the "services" collection
func (r *PaymentMongoRepo) ListAvailableServices(ctx context.Context) ([]string, error) {
	cursor, err := r.serviceListCol.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	var results []struct {
		Name string `bson:"name"`
	}
	if err := cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	if len(results) == 0 {
		return nil, errors.New("no services found")
	}

	var services []string
	for _, r := range results {
		services = append(services, r.Name)
	}
	return services, nil
}
