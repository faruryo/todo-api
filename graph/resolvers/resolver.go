package resolvers

//go:generate go run github.com/99designs/gqlgen

import (
	"github.com/faruryo/toban-api/repository"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

// Resolver データストアを持っているResolver構造体
type Resolver struct {
	Repository repository.Repository
}
