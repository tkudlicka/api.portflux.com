package entities

import (
	"time"
)

// EntityNameHolding contains the name of the entity
const EntityNameHolding = "holding"

// Holding struct
type Holding struct {
	HoldingID          string    `bson:"_holdingid,omitempty"`
	PortfolioID        string    `bson:"portfolioid"`
	BrokerID           string    `bson:"brokerid"`
	Extid              string    `bson:"extid"`
	Name               string    `bson:"name"`
	Description        string    `bson:"description"`
	Slug               string    `bson:"slug"`
	TradeDate          time.Time `bson:"trade_date"`
	TradeType          string    `bson:"trade_type"`
	Quantity           int       `bson:"quantity"`
	SharePrice         int       `bson:"share_price"`
	ExchangeRate       int       `bson:"exchange_rate"`
	ExchangeCurrencyID string    `bson:"exchange_currencyid"`
	BrokerageUnitPrice int       `bson:"brokerage_unit_price"`
	BrokerageCurrency  string    `bson:"brokerage_currency"`
	CreatedAt          time.Time `bson:"created_at"`
	UpdatedAt          time.Time `bson:"updated_at"`
}
