package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/faruryo/toban-api/graph/generated"
	"github.com/faruryo/toban-api/models"
)

func (r *mutationResolver) TodoCreate(ctx context.Context, todo models.TodoInput) (*models.Todo, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) TodoComplete(ctx context.Context, id int, updatedBy int) (*models.Todo, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) TodoDelete(ctx context.Context, id int, updatedBy int) (*models.Todo, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) UserCreate(ctx context.Context, user models.UserInput) (*models.User, error) {
	panic(fmt.Errorf("not implemented"))
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

type mutationResolver struct{ *Resolver }
