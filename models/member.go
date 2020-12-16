package models

import "time"

type Member struct {
	ID uint `json:"id"`

	SlackID string `json:"slackID"`

	Name string `json:"name"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type CreateMemberInput struct {
	SlackID string `json:"slackID"`
	Name    string `json:"name"`
}

type UpdateMemberInput struct {
	ID uint `json:"id"`

	SlackID *string `json:"slackID"`
	Name    *string `json:"name"`
}
