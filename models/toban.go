package models

import (
	"fmt"
	"io"
	"strconv"
	"time"
)

type Toban struct {
	ID uint `json:"id" gorm:"not null"`

	Name        string `json:"name" gorm:"type:VARCHAR(256);not null"`
	Description string `json:"description" gorm:"type:VARCHAR(1024);not null"`

	Interval        Interval `json:"interval" gorm:"type:ENUM('DAILY','WEEKLY','MONTHLY');not null"`
	DeadlineHour    uint     `json:"deadlineHour" gorm:"not null"`
	DeadlineWeekDay WeekDay  `json:"deadlineWeekDay" gorm:"type:ENUM('MONDAY','TUESDAY','WEDNESDAY','THURSDAY','FRIDAY','SATURDAY','SUNDAY');not null"`
	DeadlineWeek    uint     `json:"deadlineWeek" gorm:"not null"`

	Enabled bool `json:"enabled" gorm:"not null"`

	TobanMemberSequence uint `json:"tobanMemberSequence" gorm:"not null"`

	CreatedAt time.Time `json:"createdAt" gorm:"not null"`
	UpdatedAt time.Time `json:"updatedAt" gorm:"not null"`
}

type CreateTobanInput struct {
	Name        string `json:"name"`
	Description string `json:"description"`

	Interval        Interval `json:"interval"`
	DeadlineHour    uint     `json:"deadlineHour"`
	DeadlineWeekDay WeekDay  `json:"deadlineWeekDay"`
	DeadlineWeek    uint     `json:"deadlineWeek"`
}

type UpdateTobanInput struct {
	ID uint `json:"id"`

	Name        *string `json:"name"`
	Description *string `json:"description"`

	Interval        *Interval `json:"interval"`
	DeadlineHour    *uint     `json:"deadlineHour"`
	DeadlineWeekDay *WeekDay  `json:"deadlineWeekDay"`
	DeadlineWeek    *uint     `json:"deadlineWeek"`

	Enabled *bool `json:"enabled"`

	TobanMemberSequence *uint `json:"tobanMemberSequence"`
}

type Interval string

const (
	IntervalDaily   Interval = "DAILY"
	IntervalWeekly  Interval = "WEEKLY"
	IntervalMonthly Interval = "MONTHLY"
)

func (e Interval) IsValid() bool {
	switch e {
	case IntervalDaily, IntervalWeekly, IntervalMonthly:
		return true
	}
	return false
}

func (e Interval) String() string {
	return string(e)
}

func (e *Interval) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = Interval(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid Interval", str)
	}
	return nil
}

func (e Interval) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type WeekDay string

const (
	Monday    WeekDay = "MONDAY"
	Tuesday   WeekDay = "TUESDAY"
	Wednesday WeekDay = "WEDNESDAY"
	Thursday  WeekDay = "THURSDAY"
	Friday    WeekDay = "FRIDAY"
	Saturday  WeekDay = "SATURDAY"
	Sunday    WeekDay = "SUNDAY"
)

func (e WeekDay) IsValid() bool {
	switch e {
	case Monday, Tuesday, Wednesday, Thursday, Friday, Saturday, Sunday:
		return true
	}
	return false
}

func (e WeekDay) String() string {
	return string(e)
}

func (e *WeekDay) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = WeekDay(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid WeekDay", str)
	}
	return nil
}

func (e WeekDay) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}
