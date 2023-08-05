package ports

import (
	"context"

	"github.com/sergicanet9/scv-go-tools/v3/repository"
	"github.com/tkudlicka/portflux-api/core/models"
)

// UserRepositoy interface
type UserRepository interface {
	repository.Repository
	CreateMany(ctx context.Context, entities []interface{}) ([]string, error)
}

// UserService interface
type UserService interface {
	Login(ctx context.Context, credentials models.LoginUserReq) (models.LoginUserResp, error)
	Create(ctx context.Context, user models.CreateUserReq) (models.CreationResp, error)
	CreateMany(ctx context.Context, users []models.CreateUserReq) (models.MultiCreationResp, error)
	GetAll(ctx context.Context) ([]models.UserResp, error)
	GetByEmail(ctx context.Context, email string) (models.UserResp, error)
	GetByID(ctx context.Context, ID string) (models.UserResp, error)
	Update(ctx context.Context, ID string, user models.UpdateUserReq) error
	Delete(ctx context.Context, ID string) error
}
