package mongo

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/tomiristapen/banking_service/user_service/domain/model"
	"github.com/tomiristapen/banking_service/user_service/domain/repository"
)

type mongoUserRepository struct {
	collection *mongo.Collection
}

func NewUserRepository(db *mongo.Database) repository.UserRepository {
	return &mongoUserRepository{
		collection: db.Collection("users"),
	}
}

func (r *mongoUserRepository) CreateUser(ctx context.Context, user *model.User) (*model.User, error) {
	doc := bson.M{
		"_id":         user.ID,
		"name":        user.Name,
		"email":       user.Email,
		"password":    user.Password,
		"is_verified": user.IsVerified,
		"created_at":  user.CreatedAt,
		"balance":     user.Balance, // ✅ balance қосылды
	}

	_, err := r.collection.InsertOne(ctx, doc)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *mongoUserRepository) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	var result bson.M
	err := r.collection.FindOne(ctx, bson.M{"email": email}).Decode(&result)
	if err != nil {
		return nil, err
	}

	return mapToUser(result), nil
}

func (r *mongoUserRepository) FindByID(ctx context.Context, id string) (*model.User, error) {
	var result bson.M
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&result)
	if err != nil {
		return nil, errors.New("user not found")
	}

	return mapToUser(result), nil
}

func (r *mongoUserRepository) UpdateVerificationStatus(ctx context.Context, userID string, verified bool) error {
	filter := bson.M{"_id": userID}
	update := bson.M{"$set": bson.M{"is_verified": verified}}

	_, err := r.collection.UpdateOne(ctx, filter, update)
	return err
}

func (r *mongoUserRepository) UpdateUser(ctx context.Context, user *model.User) error {
	filter := bson.M{"_id": user.ID}
	update := bson.M{
		"$set": bson.M{
			"name":        user.Name,
			"email":       user.Email,
			"password":    user.Password,
			"is_verified": user.IsVerified,
			"balance":     user.Balance, // ✅ balance жаңарту
		},
	}
	_, err := r.collection.UpdateOne(ctx, filter, update)
	return err
}

func mapToUser(data bson.M) *model.User {
	return &model.User{
		ID:         data["_id"].(string),
		Name:       data["name"].(string),
		Email:      data["email"].(string),
		Password:   data["password"].(string),
		IsVerified: data["is_verified"].(bool),
		CreatedAt:  data["created_at"].(primitive.DateTime).Time(),
		Balance:    data["balance"].(float64), // ✅ balance оқу
	}
}
