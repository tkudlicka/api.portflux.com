package entities

import (
	"time"
)

// EntityNameBroker contains the name of the entity
const EntityNameBroker = "broker"

// Broker struct
type Broker struct {
	BrokerID    string    `bson:"_brokerid,omitempty"`
	Extid       string    `bson:"extid"`
	Name        string    `bson:"name"`
	Description string    `bson:"description"`
	Slug        string    `bson:"slug"`
	CreatedAt   time.Time `bson:"created_at"`
	UpdatedAt   time.Time `bson:"updated_at"`
}
