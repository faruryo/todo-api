package models

import "time"

type TobanMember struct {
	ID uint `json:"id"`

	TobanID  uint `json:"tobanID"`
	Sequence uint `json:"sequence"`
	MemberID uint `json:"memberID"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type CreateTobanMemberInput struct {
	TobanID  uint `json:"tobanID"`
	Sequence uint `json:"sequence"`
	MemberID uint `json:"memberID"`
}
