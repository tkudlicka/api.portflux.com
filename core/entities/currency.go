package entities

import (
	"time"
)

// EntityNameCurrency contains the name of the entity
const EntityNameCurrency = "currency"

// Currency struct
type Currency struct {
	CurrencyID string    `bson:"_currencyid,omitempty"`
	Code       string    `bson:"code"`
	Name       string    `bson:"name"`
	Symbol     string    `bson:"symbol"`
	CreatedAt  time.Time `bson:"created_at"`
	UpdatedAt  time.Time `bson:"updated_at"`
}
