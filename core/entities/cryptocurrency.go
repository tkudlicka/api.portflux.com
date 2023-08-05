package entities

import (
	"time"
)

// EntityNameCryptoCurrency contains the name of the entity
const EntityNameCryptoCurrency = "cryptocurrency"

// Cryptocurrency struct
type CryptoCurrency struct {
	CryptoCurrencyID string    `bson:"_crypto_currencyid,omitempty"`
	HoldingID        string    `bson:"holdingid"`
	Symbol           string    `bson:"symbol"`
	Name             string    `bson:"name"`
	CreatedAt        time.Time `bson:"created_at"`
	UpdatedAt        time.Time `bson:"updated_at"`
}
