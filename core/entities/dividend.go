package entities

import (
	"time"
)

// EntityNameDividend contains the name of the entity
const EntityNameDividend = "dividend"

// Dividend struct
type Dividend struct {
	DividendID       string    `bson:"dividendid,omitempty"`
	StockID          string    `bson:"stockid,omitempty"`
	DividendPerShare int       `bson:"dividend_per_share"`
	DividendDate     time.Time `bson:"dividend_date"`
	CreatedAt        time.Time `bson:"created_at"`
	UpdatedAt        time.Time `bson:"updated_at"`
}
