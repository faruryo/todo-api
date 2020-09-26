package repositories

import (
	"context"
	"errors"

	"github.com/faruryo/toban-api/models"
	"gorm.io/gorm"
)

type TobanRepository interface {
	Get(ctx context.Context, id uint) (*models.Toban, error)
	GetAll(ctx context.Context) ([]*models.Toban, error)
	Create(ctx context.Context, toban *models.Toban) (*models.Toban, error)
	Update(ctx context.Context, toban *models.Toban) (*models.Toban, error)
	Delete(ctx context.Context, id uint) (bool, error)
}

func NewTobanRepository(db *gorm.DB) TobanRepository {
	db.AutoMigrate(&models.Toban{})

	return NewTobanRepositoryNoMigrate(db)
}

func NewTobanRepositoryNoMigrate(db *gorm.DB) TobanRepository {
	return &tobanRepository{
		db: db,
	}
}

// Interface実装チェック
var _ TobanRepository = (*tobanRepository)(nil)

type tobanRepository struct {
	db *gorm.DB
}

func (r tobanRepository) Get(ctx context.Context, id uint) (*models.Toban, error) {
	var toban models.Toban
	err := r.db.First(&toban, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		// return &toban, nil
		return nil, ErrNoSuchEntity
	}
	if err != nil {
		return nil, err
	}

	return &toban, nil
}

func (r tobanRepository) GetAll(ctx context.Context) ([]*models.Toban, error) {
	var tobans []*models.Toban
	if err := r.db.Find(&tobans).Error; err != nil {
		return nil, err
	}

	return tobans, nil
}

func (r tobanRepository) Create(ctx context.Context, toban *models.Toban) (*models.Toban, error) {
	if toban.ID != 0 {
		return nil, ErrBadRequestIdMustBeZero
	}
	if !toban.CreatedAt.IsZero() {
		return nil, ErrBadRequestUpdateCreatedAt
	}
	if !toban.UpdatedAt.IsZero() {
		return nil, ErrBadRequestUpdateUpdatedAt
	}

	if err := r.db.Create(toban).Error; err != nil {
		return nil, err
	}

	return toban, nil
}

func (r tobanRepository) Update(ctx context.Context, toban *models.Toban) (*models.Toban, error) {
	if toban.ID == 0 {
		return nil, ErrBadRequestIdMustNotBeZero
	}
	if !toban.CreatedAt.IsZero() {
		return nil, ErrBadRequestUpdateCreatedAt
	}
	if !toban.UpdatedAt.IsZero() {
		return nil, ErrBadRequestUpdateUpdatedAt
	}
	// TODO:トランザクションで現在値からを元にアップデートするように変更する

	if err := r.db.Save(&toban).Error; err != nil {
		return nil, err
	}

	return toban, nil
}

func (r tobanRepository) Delete(ctx context.Context, id uint) (bool, error) {
	if id == 0 {
		return false, ErrBadRequestIdMustNotBeZero
	}

	var toban models.Toban
	if err := r.db.Delete(toban, id).Error; err != nil {
		return false, err
	}

	return true, nil
}
