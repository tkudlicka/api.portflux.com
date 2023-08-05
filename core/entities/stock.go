package entities

import (
	"time"
)

// EntityNameStock contains the name of the entity
const EntityNameStock = "stock"

// Stock struct
type Stock struct {
	StockID      string    `bson:"_stockid,omitempty"`
	HoldingID    string    `bson:"holdingid"`
	Extid        string    `bson:"extid"`
	TickerSymbol string    `bson:"ticker_symbol"`
	CompanyName  string    `bson:"company_name"`
	Slug         string    `bson:"slug"`
	CreatedAt    time.Time `bson:"created_at"`
	UpdatedAt    time.Time `bson:"updated_at"`
}
