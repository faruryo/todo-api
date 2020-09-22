package models

import "time"

type Toban struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`

	Interval     string `json:"interval"`
	DeadlineHour int    `json:"deadlineHour"`
	DeadlineDay  string `json:"deadlineDay"`
	DeadlineWeek int    `json:"deadlineWeek"`

	Enabled bool `json:"enabled"`

	TobanMemberSequence int `json:"tobanMemberSequence"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type CreateTobanInput struct {
	Name        string `json:"name"`
	Description string `json:"description"`

	Interval     string `json:"interval"`
	DeadlineHour int    `json:"deadlineHour"`
	DeadlineDay  string `json:"deadlineDay"`
	DeadlineWeek int    `json:"deadlineWeek"`

	Enabled bool `json:"enabled"`

	TobanMemberSequence int `json:"tobanMemberSequence"`
}
