package ports

import (
	"context"

	"github.com/sergicanet9/scv-go-tools/v3/repository"
	"github.com/tkudlicka/portflux-api/core/models"
)

// TransactionRepositoy interface
type TransactionRepository interface {
	repository.Repository
	CreateMany(ctx context.Context, entities []interface{}) ([]string, error)
}

// TransactionService interface
type TransactionService interface {
	Create(ctx context.Context, user models.CreateUserReq) (models.CreationResp, error)
	CreateMany(ctx context.Context, users []models.CreateUserReq) (models.MultiCreationResp, error)
	GetAll(ctx context.Context) ([]models.UserResp, error)
	GetByID(ctx context.Context, ID string) (models.UserResp, error)
	Update(ctx context.Context, ID string, user models.UpdateUserReq) error
	Delete(ctx context.Context, ID string) error
}
