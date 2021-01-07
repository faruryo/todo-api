package resolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/faruryo/toban-api/graph/generated"
	"github.com/faruryo/toban-api/models"
)

func (r *mutationResolver) CreateTobanWariate(ctx context.Context, input models.CreateTobanWariateInput) (*models.TobanWariate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) CreateToban(ctx context.Context, input models.CreateTobanInput) (*models.Toban, error) {
	t := &models.Toban{
		Name:        input.Name,
		Description: input.Description,

		Interval:        input.Interval,
		DeadlineHour:    input.DeadlineHour,
		DeadlineWeekDay: input.DeadlineWeekDay,
		DeadlineWeek:    input.DeadlineWeek,

		Enabled: true,

		TobanMemberSequence: 0,
	}

	return r.Repository.CreateToban(ctx, t)
}

func (r *mutationResolver) DeleteToban(ctx context.Context, id uint) (bool, error) {
	return r.Repository.DeleteTobanByID(ctx, id)
}

func (r *mutationResolver) UpdateToban(ctx context.Context, input models.UpdateTobanInput) (*models.Toban, error) {
	return r.Repository.UpdateToban(ctx, &input)
}

func (r *mutationResolver) CreateTobanMember(ctx context.Context, input models.CreateTobanMemberInput) (*models.TobanMember, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) CreateMember(ctx context.Context, input models.CreateMemberInput) (*models.Member, error) {
	m := &models.Member{
		SlackID: input.SlackID,
		Name:    input.Name,
	}

	return r.Repository.CreateMember(ctx, m)
}

func (r *mutationResolver) DeleteMember(ctx context.Context, id uint) (bool, error) {
	return r.Repository.DeleteMemberByID(ctx, id)
}

func (r *mutationResolver) UpdateMember(ctx context.Context, input models.UpdateMemberInput) (*models.Member, error) {
	return r.Repository.UpdateMember(ctx, &input)
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

type mutationResolver struct{ *Resolver }
