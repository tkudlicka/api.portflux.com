package ports

import (
	"context"

	"github.com/sergicanet9/scv-go-tools/v3/repository"
	"github.com/tkudlicka/portflux-api/core/models"
)

// CryptoCurrencyRepositoy interface
type CryptoCurrencyRepository interface {
	repository.Repository
	CreateMany(ctx context.Context, entities []interface{}) ([]string, error)
}

// CryptoCurrencyService interface
type CryptoCurrencyService interface {
	Create(ctx context.Context, cryptoCurrency models.CreateCryptoCurrencyReq) (models.CreationResp, error)
	CreateMany(ctx context.Context, cryptoCurrency []models.CreateCryptoCurrencyReq) (models.MultiCreationResp, error)
	GetAll(ctx context.Context) ([]models.CryptoCurrencyResp, error)
	GetByID(ctx context.Context, ID string) (models.CryptoCurrencyResp, error)
	Update(ctx context.Context, ID string, cryptoCurrencyModels models.UpdateCryptoCurrencyReq) error
	Delete(ctx context.Context, ID string) error
}
