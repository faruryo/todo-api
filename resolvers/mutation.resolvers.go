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

func (r *mutationResolver) CreateTobanWariate(ctx context.Context, tw models.CreateTobanWariateInput) (*models.TobanWariate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) CreateToban(ctx context.Context, toban models.CreateTobanInput) (*models.Toban, error) {
	t := &models.Toban{
		Name:        toban.Name,
		Description: toban.Description,

		Interval:     toban.Interval,
		DeadlineHour: toban.DeadlineHour,
		DeadlineDay:  toban.DeadlineDay,
		DeadlineWeek: toban.DeadlineWeek,

		Enabled: true,

		TobanMemberSequence: 0,

		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	return r.TobanRepository.Create(ctx, t)
}

func (r *mutationResolver) DeleteToban(ctx context.Context, id int) (bool, error) {
	return r.TobanRepository.Delete(ctx, id)
}

func (r *mutationResolver) CreateTobanMember(ctx context.Context, tm models.CreateTobanMemberInput) (*models.TobanMember, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) CreateMember(ctx context.Context, member models.CreateMemberInput) (*models.Member, error) {
	id := len(r.members)
	m := &models.Member{
		ID:      id,
		SlackID: member.SlackID,

		Name: member.Name,

		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if r.members == nil {
		r.members = map[int]*models.Member{}
	}
	r.members[id] = m
	return m, nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

type mutationResolver struct{ *Resolver }
