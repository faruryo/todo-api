package repository

import (
	"context"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/faruryo/toban-api/models"
	"github.com/google/go-cmp/cmp"
)

func TestGetTobanByID(t *testing.T) {
	repo, mock := getRepoAndMock(t)

	dbOutput := &models.Toban{
		ID:                  1,
		Name:                "掃除機",
		Description:         "desc",
		Interval:            "DAILY",
		DeadlineHour:        23,
		DeadlineWeekDay:     "SUNDAY",
		DeadlineWeek:        0,
		Enabled:             true,
		TobanMemberSequence: 0,
		CreatedAt:           time.Now(),
		UpdatedAt:           time.Now(),
	}

	// sqlmock準備
	rows := sqlmock.NewRows([]string{"id", "name", "description", "interval", "deadline_hour", "deadline_week_day", "deadline_week", "enabled", "toban_member_sequence", "created_at", "updated_at"}).
		AddRow(dbOutput.ID, dbOutput.Name, dbOutput.Description, dbOutput.Interval, dbOutput.DeadlineHour, dbOutput.DeadlineWeekDay, dbOutput.DeadlineWeek, dbOutput.Enabled, dbOutput.TobanMemberSequence, dbOutput.CreatedAt, dbOutput.UpdatedAt)
	sql := regexp.QuoteMeta("SELECT * FROM `tobans`")
	mock.ExpectQuery(sql).WithArgs(dbOutput.ID).WillReturnRows(rows)

	// Test開始
	output, err := repo.GetTobanByID(context.Background(), dbOutput.ID)
	if err != nil {
		t.Fatal(err)
	}
	if diff := cmp.Diff(dbOutput, output); diff != "" {
		t.Errorf("input and output are different\n%s", diff)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestGetTobanByID_Error(t *testing.T) {
	repo, mock := getRepoAndMock(t)

	input := &models.Toban{
		ID: 1,
	}

	// sqlmock準備
	rows := sqlmock.NewRows([]string{"id", "name", "description", "interval", "deadline_hour", "deadline_week_day", "deadline_week", "enabled", "toban_member_sequence", "created_at", "updated_at"})
	sql := regexp.QuoteMeta("SELECT * FROM `tobans`")
	mock.ExpectQuery(sql).WithArgs(input.ID).WillReturnRows(rows)

	// Test開始
	if _, err := repo.GetTobanByID(context.Background(), input.ID); err != ErrNoSuchEntity {
		t.Fatalf("it doesn't return an error when no such entity. %v", err)
	}
}

func TestGetAllTobans(t *testing.T) {
	repo, mock := getRepoAndMock(t)

	dbOutputs := []*models.Toban{
		{
			ID:                  1,
			Name:                "掃除機",
			Description:         "掃除機をかける",
			Interval:            "WEEKLY",
			DeadlineHour:        23,
			DeadlineWeekDay:     "SUNDAY",
			DeadlineWeek:        0,
			Enabled:             true,
			TobanMemberSequence: 0,
			CreatedAt:           time.Now(),
			UpdatedAt:           time.Now(),
		},
		{
			ID:                  5,
			Name:                "家計簿",
			Description:         "家計簿をつける",
			Interval:            "WEEKLY",
			DeadlineHour:        23,
			DeadlineWeekDay:     "SUNDAY",
			DeadlineWeek:        0,
			Enabled:             true,
			TobanMemberSequence: 0,
			CreatedAt:           time.Now(),
			UpdatedAt:           time.Now(),
		},
	}

	// sqlmock準備
	rows := sqlmock.NewRows([]string{"id", "name", "description", "interval", "deadline_hour", "deadline_week_day", "deadline_week", "enabled", "toban_member_sequence", "created_at", "updated_at"})
	for _, dbOutput := range dbOutputs {
		rows.AddRow(dbOutput.ID, dbOutput.Name, dbOutput.Description, dbOutput.Interval, dbOutput.DeadlineHour, dbOutput.DeadlineWeekDay, dbOutput.DeadlineWeek, dbOutput.Enabled, dbOutput.TobanMemberSequence, dbOutput.CreatedAt, dbOutput.UpdatedAt)
	}
	sql := regexp.QuoteMeta("SELECT * FROM `tobans`")
	mock.ExpectQuery(sql).WillReturnRows(rows)

	// Test開始
	output, err := repo.GetAllTobans(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if diff := cmp.Diff(dbOutputs, output); diff != "" {
		t.Errorf("input and output are different\n%s", diff)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestCreateToban(t *testing.T) {
	repo, mock := getRepoAndMock(t)

	input := &models.Toban{
		ID:                  0,
		Name:                "掃除機",
		Description:         "desc",
		Interval:            "DAILY",
		DeadlineHour:        23,
		DeadlineWeekDay:     "SUNDAY",
		DeadlineWeek:        0,
		Enabled:             true,
		TobanMemberSequence: 0,
	}

	// sqlmock準備
	sql := regexp.QuoteMeta("INSERT INTO `tobans` (`name`,`description`,`interval`,`deadline_hour`,`deadline_week_day`,`deadline_week`,`enabled`,`toban_member_sequence`,`created_at`,`updated_at`) VALUES (?,?,?,?,?,?,?,?,?,?)")
	mock.ExpectExec(sql).WithArgs(input.Name, input.Description, input.Interval, input.DeadlineHour, input.DeadlineWeekDay, input.DeadlineWeek, input.Enabled, input.TobanMemberSequence, AnyTime{}, AnyTime{}).WillReturnResult(sqlmock.NewResult(1, 1))

	// Test開始
	_, err := repo.CreateToban(context.Background(), input)
	if err != nil {
		t.Fatal(err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestCreateToban_Error(t *testing.T) {
	repo, _ := getRepoAndMock(t)

	cases := []struct {
		input *models.Toban
		err   error
	}{
		{
			input: &models.Toban{ID: 1},
			err:   ErrBadRequestIDMustBeZero,
		},
		{
			input: &models.Toban{CreatedAt: time.Now()},
			err:   ErrBadRequestUpdateCreatedAt,
		},
		{
			input: &models.Toban{UpdatedAt: time.Now()},
			err:   ErrBadRequestUpdateUpdatedAt,
		},
	}

	for _, c := range cases {
		if _, err := repo.CreateToban(context.Background(), c.input); err != c.err {
			t.Errorf("Reverse(%v) => err(%v), want err(%v)", c.input, err, c.err)
		}
	}
}

func TestUpdateToban(t *testing.T) {
	repo, mock := getRepoAndMock(t)

	dbOutput := &models.Toban{
		ID:                  1,
		Name:                "掃除機",
		Description:         "desc",
		Interval:            "DAILY",
		DeadlineHour:        23,
		DeadlineWeekDay:     "SUNDAY",
		DeadlineWeek:        0,
		Enabled:             true,
		TobanMemberSequence: 0,
		CreatedAt:           time.Now(),
		UpdatedAt:           time.Now(),
	}
	input := &models.UpdateTobanInput{
		ID:                  dbOutput.ID,
		Name:                &dbOutput.Name,
		Description:         &dbOutput.Description,
		Interval:            &dbOutput.Interval,
		DeadlineHour:        &dbOutput.DeadlineHour,
		DeadlineWeekDay:     &dbOutput.DeadlineWeekDay,
		DeadlineWeek:        &dbOutput.DeadlineWeek,
		Enabled:             &dbOutput.Enabled,
		TobanMemberSequence: &dbOutput.TobanMemberSequence,
	}

	// sqlmock準備
	mock.ExpectBegin()
	rows := sqlmock.NewRows([]string{"id", "name", "description", "interval", "deadline_hour", "deadline_week_day", "deadline_week", "enabled", "toban_member_sequence", "created_at", "updated_at"}).
		AddRow(dbOutput.ID, dbOutput.Name, dbOutput.Description, dbOutput.Interval, dbOutput.DeadlineHour, dbOutput.DeadlineWeekDay, dbOutput.DeadlineWeek, dbOutput.Enabled, dbOutput.TobanMemberSequence, dbOutput.CreatedAt, dbOutput.UpdatedAt)
	sql := regexp.QuoteMeta("SELECT * FROM `tobans`")
	mock.ExpectQuery(sql).WithArgs(input.ID).WillReturnRows(rows)
	sql = regexp.QuoteMeta("UPDATE `tobans` SET `name`=?,`description`=?,`interval`=?,`deadline_hour`=?,`deadline_week_day`=?,`deadline_week`=?,`enabled`=?,`toban_member_sequence`=?,`created_at`=?,`updated_at`=? WHERE `id` = ?")
	mock.ExpectExec(sql).WithArgs(dbOutput.Name, dbOutput.Description, dbOutput.Interval, dbOutput.DeadlineHour, dbOutput.DeadlineWeekDay, dbOutput.DeadlineWeek, dbOutput.Enabled, dbOutput.TobanMemberSequence, AnyTime{}, AnyTime{}, input.ID).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	// Test開始
	_, err := repo.UpdateToban(context.Background(), input)
	if err != nil {
		t.Fatal(err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestUpdateToban_Error(t *testing.T) {
	repo, _ := getRepoAndMock(t)

	inputInstance := &models.Toban{
		Name: "TestUpdate_IdMustNotBeZero",
	}

	cases := []struct {
		input *models.UpdateTobanInput
		err   error
	}{
		{
			input: &models.UpdateTobanInput{
				Name: &inputInstance.Name,
			},
			err: ErrBadRequestIDMustNotBeZero,
		},
	}

	for _, c := range cases {
		if _, err := repo.UpdateToban(context.Background(), c.input); err != c.err {
			t.Errorf("Reverse(%v) => err(%v), want err(%v)", c.input, err, c.err)
		}
	}
}

func TestDeleteTobanByID(t *testing.T) {
	repo, mock := getRepoAndMock(t)

	var input uint = 1

	// sqlmock準備
	sql := regexp.QuoteMeta("DELETE FROM `tobans` WHERE `tobans`.`id` = ?")
	mock.ExpectExec(sql).WithArgs(input).WillReturnResult(sqlmock.NewResult(1, 1))

	// Test開始
	output, err := repo.DeleteTobanByID(context.Background(), input)
	if err != nil {
		t.Fatalf("Unexpected error :%v", err)
	}
	if !output {
		t.Errorf("output: %v != true", output)
	}

	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
func TestDeleteTobanByIDNoSuchEntity(t *testing.T) {
	repo, mock := getRepoAndMock(t)

	var input uint = 1

	// sqlmock準備
	sql := regexp.QuoteMeta("DELETE FROM `tobans` WHERE `tobans`.`id` = ?")
	mock.ExpectExec(sql).WithArgs(input).WillReturnResult(sqlmock.NewResult(0, 0))

	// Test開始
	output, err := repo.DeleteTobanByID(context.Background(), input)
	if err != nil {
		t.Fatalf("Unexpected error :%v", err)
	}
	if !output {
		t.Errorf("output: %v != true", output)
	}

	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestDeleteTobanByID_Error(t *testing.T) {
	repo, _ := getRepoAndMock(t)

	cases := []struct {
		input  uint
		output bool
		err    error
	}{
		{
			input:  0,
			output: false,
			err:    ErrBadRequestIDMustNotBeZero,
		},
	}

	for _, c := range cases {
		output, err := repo.DeleteTobanByID(context.Background(), c.input)
		if err != c.err {
			t.Errorf("Reverse(%v) => err(%v), want err(%v)", c.input, err, c.err)
		}
		if output != c.output {
			t.Errorf("Reverse(%v) => err(%v), want err(%v)", c.input, output, c.output)
		}
	}
}
