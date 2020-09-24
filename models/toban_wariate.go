package models

import "time"

type TobanWariate struct {
	ID uint `json:"id"`

	TobanID       uint `json:"tobanID"`
	TobanSequence uint `json:"sequence"`
	MemberID      uint `json:"memberID"`

	IsDone bool      `json:"isDone"`
	DoneAt time.Time `json:"doneAt"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type CreateTobanWariateInput struct {
	TobanID       uint `json:"tobanID"`
	TobanSequence uint `json:"tobanSequence"`
	MemberID      uint `json:"memberID"`
}
