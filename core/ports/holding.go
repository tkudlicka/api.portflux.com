package ports

import (
	"context"

	"github.com/sergicanet9/scv-go-tools/v3/repository"
	"github.com/tkudlicka/portflux-api/core/models"
)

// HoldingRepositoy interface
type HoldingRepository interface {
	repository.Repository
	CreateMany(ctx context.Context, entities []interface{}) ([]string, error)
}

// HodingService interface
type HoldingService interface {
	Create(ctx context.Context, holding models.CreateHoldingReq) (models.CreationResp, error)
	CreateMany(ctx context.Context, holdings []models.CreateHoldingReq) (models.MultiCreationResp, error)
	GetAll(ctx context.Context) ([]models.HoldingResp, error)
	GetByID(ctx context.Context, ID string) (models.HoldingResp, error)
	Update(ctx context.Context, ID string, holding models.UpdateHoldingReq) error
	Delete(ctx context.Context, ID string) error
}
