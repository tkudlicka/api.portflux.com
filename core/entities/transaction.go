package entities

import (
	"time"
)

// EntityNameTransaction contains the name of the entity
const EntityNameTransaction = "transaction"

// Transaction struct
type Transaction struct {
	TransactionID    string    `bson:"transactionid,omitempty"`
	StockID          string    `bson:"stockid"`
	CryptocurrencyID string    `bson:"cryptocurrencyid"`
	Quantity         string    `bson:"quantity"`
	TransactionPrice string    `bson:"transaction_price"`
	TransactionDate  time.Time `bson:"transaction_date"`
	CreatedAt        time.Time `bson:"created_at"`
	UpdatedAt        time.Time `bson:"updated_at"`
}
