package model

type Transaction struct {
	ID         string  `bson:"_id,omitempty" json:"id"`
	FromUserID string  `bson:"from_user_id" json:"from_user_id"`
	ToUserID   string  `bson:"to_user_id" json:"to_user_id"`
	Amount     float64 `bson:"amount" json:"amount"`
	Currency   string  `bson:"currency" json:"currency"`
	Status     string  `bson:"status" json:"status"`
	CreatedAt  string  `bson:"created_at" json:"created_at"`
}
