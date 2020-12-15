package resolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/faruryo/toban-api/graph/generated"
	"github.com/faruryo/toban-api/models"
)

func (r *queryResolver) TobanWariate(ctx context.Context, id uint) (*models.TobanWariate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) TobanWariates(ctx context.Context) ([]*models.TobanWariate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Toban(ctx context.Context, id uint) (*models.Toban, error) {
	return r.Repository.GetTobanByID(ctx, id)
}

func (r *queryResolver) Tobans(ctx context.Context) ([]*models.Toban, error) {
	return r.Repository.GetAllTobans(ctx)
}

func (r *queryResolver) TobanMember(ctx context.Context, id uint) (*models.TobanMember, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) TobanMembers(ctx context.Context) ([]*models.TobanMember, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Member(ctx context.Context, id uint) (*models.Member, error) {
	return r.Repository.GetMemberByID(ctx, id)
}

func (r *queryResolver) Members(ctx context.Context) ([]*models.Member, error) {
	return r.Repository.GetAllMembers(ctx)
}

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type queryResolver struct{ *Resolver }
