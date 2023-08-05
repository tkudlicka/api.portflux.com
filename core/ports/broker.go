package ports

import (
	"context"

	"github.com/sergicanet9/scv-go-tools/v3/repository"
	"github.com/tkudlicka/portflux-api/core/models"
)

// BrokerRepository interface
type BrokerRepository interface {
	repository.Repository
	CreateMany(ctx context.Context, entities []interface{}) ([]string, error)
	GetBySlug(ctx context.Context, Slug string) (interface{}, error)
}

// UserService interface
type BrokerService interface {
	Create(ctx context.Context, broker models.CreateBrokerReq) (models.CreationResp, error)
	CreateMany(ctx context.Context, brokers []models.CreateBrokerReq) (models.MultiCreationResp, error)
	GetAll(ctx context.Context) ([]models.BrokerResp, error)
	GetBySlug(ctx context.Context, Slug string) (models.BrokerResp, error)
	GetByID(ctx context.Context, ID string) (models.BrokerResp, error)
	Update(ctx context.Context, ID string, broker models.UpdateBrokerReq) error
	Delete(ctx context.Context, ID string) error
}
