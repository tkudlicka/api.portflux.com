package models

import (
	"fmt"
	"strings"
	"time"

	"github.com/sergicanet9/scv-go-tools/v3/wrappers"
)

// BrokerResp broker response struct
type BrokerResp struct {
	BrokerID    string    `json:"brokerid"`
	Extid       string    `json:"extid"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Slug        string    `json:"slug"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// CreateBrokerReq broker request struct
type CreateBrokerReq struct {
	BrokerID    string    `json:"-"`
	Extid       string    `json:"extid"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Slug        string    `json:"-"`
	CreatedAt   time.Time `json:"-"`
	UpdatedAt   time.Time `json:"-"`
}

func (req CreateBrokerReq) Validate() error {
	var msgs []string

	if req.Extid == "" {
		msgs = append(msgs, "external ID cannot be empty")
	}
	if req.Name == "" {
		msgs = append(msgs, "name cannot be empty")
	}

	if req.Description == "" {
		msgs = append(msgs, "description cannot be empty")
	}

	if len(msgs) > 0 {
		return wrappers.NewValidationErr(fmt.Errorf(strings.Join(msgs, " | ")))
	}

	return nil
}

// UpdateBrokerReq update broker request struct
type UpdateBrokerReq struct {
	BrokerID    string    `json:"-"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"-"`
	UpdatedAt   time.Time `json:"-"`
}

func (req UpdateBrokerReq) Validate() error {
	var msgs []string

	if req.Name == "" {
		msgs = append(msgs, "name cannot be empty")
	}

	if req.Description == "" {
		msgs = append(msgs, "description cannot be empty")
	}

	if len(msgs) > 0 {
		return wrappers.NewValidationErr(fmt.Errorf(strings.Join(msgs, " | ")))
	}

	return nil
}
