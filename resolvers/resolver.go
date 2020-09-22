package resolvers

//go:generate go run github.com/99designs/gqlgen

import (
	"github.com/faruryo/toban-api/models"
	"github.com/faruryo/toban-api/repositories"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

// Resolver データストアを持っているResolver構造体
type Resolver struct {
	tobans          []*models.Toban
	members         map[int]*models.Member
	TobanRepository repositories.TobanRepository
}
