package resolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/faruryo/toban-api/graph/generated"
	"github.com/faruryo/toban-api/models"
)

func (r *tobanMemberResolver) TobanID(ctx context.Context, obj *models.TobanMember) (*models.Toban, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *tobanMemberResolver) MemberID(ctx context.Context, obj *models.TobanMember) (*models.Member, error) {
	member, ok := r.members[obj.MemberID]
	if !ok {
		return nil, fmt.Errorf("member %d does not exist", obj.MemberID)
	}

	return member, nil
}

// TobanMember returns generated.TobanMemberResolver implementation.
func (r *Resolver) TobanMember() generated.TobanMemberResolver { return &tobanMemberResolver{r} }

type tobanMemberResolver struct{ *Resolver }
