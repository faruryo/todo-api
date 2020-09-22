package models

import "time"

type TobanWariate struct {
	ID int `json:"id"`

	TobanID       int `json:"tobanID"`
	TobanSequence int `json:"sequence"`
	MemberID      int `json:"memberID"`

	IsDone bool      `json:"isDone"`
	DoneAt time.Time `json:"doneAt"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type CreateTobanWariateInput struct {
	TobanID       int `json:"tobanID"`
	TobanSequence int `json:"tobanSequence"`
	MemberID      int `json:"memberID"`
}
