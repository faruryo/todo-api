package models

import "time"

type TobanMember struct {
	ID int `json:"id"`

	TobanID  int `json:"tobanID"`
	Sequence int `json:"sequence"`
	MemberID int `json:"memberID"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type CreateTobanMemberInput struct {
	TobanID  int `json:"tobanID"`
	Sequence int `json:"sequence"`
	MemberID int `json:"memberID"`
}
