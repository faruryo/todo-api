package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/faruryo/toban-api/graph/generated"
	"github.com/faruryo/toban-api/models"
)

func (r *queryResolver) Todo(ctx context.Context, id int) (*models.Todo, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Todos(ctx context.Context, limit *int, offset *int) ([]*models.Todo, error) {
	return r.todos, nil
}

func (r *queryResolver) User(ctx context.Context, id int) (*models.User, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Users(ctx context.Context, limit *int, offset *int) ([]*models.User, error) {
	panic(fmt.Errorf("not implemented"))
}

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type queryResolver struct{ *Resolver }
