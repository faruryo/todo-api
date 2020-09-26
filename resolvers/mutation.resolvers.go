package resolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
	"time"

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

	return r.TobanRepository.Create(ctx, t)
}

func (r *mutationResolver) DeleteToban(ctx context.Context, id uint) (bool, error) {
	return r.TobanRepository.Delete(ctx, id)
}

func (r *mutationResolver) UpdateToban(ctx context.Context, input models.UpdateTobanInput) (*models.Toban, error) {
	t, err := r.TobanRepository.Get(ctx, input.ID)
	if err != nil {
		return nil, err
	}

	if input.Name != nil {
		t.Name = *input.Name
	}
	if input.Description != nil {
		t.Description = *input.Description
	}
	if input.Interval != nil {
		t.Interval = *input.Interval
	}
	if input.DeadlineHour != nil {
		t.DeadlineHour = *input.DeadlineHour
	}
	if input.DeadlineWeekDay != nil {
		t.DeadlineWeekDay = *input.DeadlineWeekDay
	}
	if input.DeadlineWeek != nil {
		t.DeadlineWeek = *input.DeadlineWeek
	}
	if input.Enabled != nil {
		t.Enabled = *input.Enabled
	}
	if input.TobanMemberSequence != nil {
		t.TobanMemberSequence = *input.TobanMemberSequence
	}
	t.CreatedAt = time.Time{}
	t.UpdatedAt = time.Time{}
	return r.TobanRepository.Update(ctx, t)
}

func (r *mutationResolver) CreateTobanMember(ctx context.Context, input models.CreateTobanMemberInput) (*models.TobanMember, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) CreateMember(ctx context.Context, input models.CreateMemberInput) (*models.Member, error) {
	id := uint(len(r.members))
	m := &models.Member{
		ID:      id,
		SlackID: input.SlackID,

		Name: input.Name,

		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if r.members == nil {
		r.members = map[uint]*models.Member{}
	}
	r.members[id] = m
	return m, nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

type mutationResolver struct{ *Resolver }
