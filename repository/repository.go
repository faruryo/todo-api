package repository

import (
	"context"

	"github.com/faruryo/toban-api/models"
	"gorm.io/gorm"
)

type Repository interface {
	GetTobanByID(ctx context.Context, id uint) (*models.Toban, error)
	GetAllTobans(ctx context.Context) ([]*models.Toban, error)
	CreateToban(ctx context.Context, toban *models.Toban) (*models.Toban, error)
	UpdateToban(ctx context.Context, toban *models.UpdateTobanInput) (*models.Toban, error)
	DeleteTobanByID(ctx context.Context, id uint) (bool, error)
}

func NewRepository(db *gorm.DB) Repository {
	db.AutoMigrate(&models.Toban{})

	return NewRepositoryNoMigrate(db)
}

func NewRepositoryNoMigrate(db *gorm.DB) Repository {
	return &repository{
		db: db,
	}
}

// Interface implementation check
var _ Repository = (*repository)(nil)

type repository struct {
	db *gorm.DB
}
