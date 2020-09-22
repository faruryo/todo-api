package repositories

import (
	"context"

	"github.com/faruryo/toban-api/models"
	"github.com/jinzhu/gorm"
)

type TobanRepository interface {
	Get(ctx context.Context, id int) (*models.Toban, error)
	GetAll(ctx context.Context) ([]*models.Toban, error)
	Create(ctx context.Context, Toban *models.Toban) (*models.Toban, error)
	Update(ctx context.Context, Toban *models.Toban) (*models.Toban, error)
	Delete(ctx context.Context, id int) (*models.Toban, error)
}

func NewTobanRepository(db *gorm.DB) TobanRepository {
	db.AutoMigrate(&models.Toban{})

	return &tobanRepository{
		db: db,
	}
}

// Interface実装チェック
var _ TobanRepository = (*tobanRepository)(nil)

type tobanRepository struct {
	db *gorm.DB
}

func (r tobanRepository) Get(ctx context.Context, id int) (*models.Toban, error) {
	panic("not implemented") // TODO: Implement
}

func (r tobanRepository) GetAll(ctx context.Context) ([]*models.Toban, error) {
	panic("not implemented") // TODO: Implement
}

func (r tobanRepository) Create(ctx context.Context, Toban *models.Toban) (*models.Toban, error) {
	panic("not implemented") // TODO: Implement
}

func (r tobanRepository) Update(ctx context.Context, Toban *models.Toban) (*models.Toban, error) {
	panic("not implemented") // TODO: Implement
}

func (r tobanRepository) Delete(ctx context.Context, id int) (*models.Toban, error) {
	panic("not implemented") // TODO: Implement
}
