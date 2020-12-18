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

func TestGetMemberByID(t *testing.T) {
	repo, mock := getRepoAndMock(t)

	dbOutput := &models.Member{
		ID:      1,
		SlackID: "slack01",
		Name:    "slack.01",
	}

	// Prepare sqlmock
	rows := sqlmock.NewRows([]string{"id", "slack_id", "name", "created_at", "updated_at"}).
		AddRow(dbOutput.ID, dbOutput.SlackID, dbOutput.Name, dbOutput.CreatedAt, dbOutput.UpdatedAt)
	sql := regexp.QuoteMeta("SELECT * FROM `members`")
	mock.ExpectQuery(sql).WithArgs(dbOutput.ID).WillReturnRows(rows)

	// Start Test
	output, err := repo.GetMemberByID(context.Background(), dbOutput.ID)
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

func TestGetMemberByID_Error(t *testing.T) {
	repo, mock := getRepoAndMock(t)

	input := &models.Member{
		ID: 1,
	}

	// Prepare sqlmock
	rows := sqlmock.NewRows([]string{"id", "slack_id", "name", "created_at", "updated_at"})
	sql := regexp.QuoteMeta("SELECT * FROM `members`")
	mock.ExpectQuery(sql).WithArgs(input.ID).WillReturnRows(rows)

	// Start Test
	if _, err := repo.GetMemberByID(context.Background(), input.ID); err != ErrNoSuchEntity {
		t.Fatalf("it doesn't return an error when no such entity. %v", err)
	}
}

func TestGetAllMembers(t *testing.T) {
	repo, mock := getRepoAndMock(t)

	dbOutputs := []*models.Member{
		{
			ID:        1,
			SlackID:   "slack01",
			Name:      "slack.01",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			ID:        5,
			SlackID:   "slack05",
			Name:      "slack.05",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	// Prepare sqlmock
	rows := sqlmock.NewRows([]string{"id", "slack_id", "name", "created_at", "updated_at"})
	for _, dbOutput := range dbOutputs {
		rows.AddRow(dbOutput.ID, dbOutput.SlackID, dbOutput.Name, dbOutput.CreatedAt, dbOutput.UpdatedAt)
	}
	sql := regexp.QuoteMeta("SELECT * FROM `members`")
	mock.ExpectQuery(sql).WillReturnRows(rows)

	// Start Test
	output, err := repo.GetAllMembers(context.Background())
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

func TestCreateMember(t *testing.T) {
	repo, mock := getRepoAndMock(t)

	input := &models.Member{
		ID:      0,
		SlackID: "slack01",
		Name:    "slack.01",
	}

	// Prepare sqlmock
	sql := regexp.QuoteMeta("INSERT INTO `members` (`slack_id`,`name`,`created_at`,`updated_at`) VALUES (?,?,?,?)")
	mock.ExpectExec(sql).WithArgs(input.SlackID, input.Name, AnyTime{}, AnyTime{}).WillReturnResult(sqlmock.NewResult(1, 1))

	// Start Test
	_, err := repo.CreateMember(context.Background(), input)
	if err != nil {
		t.Fatal(err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestCreateMember_Error(t *testing.T) {
	repo, _ := getRepoAndMock(t)

	cases := []struct {
		input *models.Member
		err   error
	}{
		{
			input: &models.Member{ID: 1},
			err:   ErrBadRequestIDMustBeZero,
		},
		{
			input: &models.Member{CreatedAt: time.Now()},
			err:   ErrBadRequestUpdateCreatedAt,
		},
		{
			input: &models.Member{UpdatedAt: time.Now()},
			err:   ErrBadRequestUpdateUpdatedAt,
		},
	}

	for _, c := range cases {
		if _, err := repo.CreateMember(context.Background(), c.input); err != c.err {
			t.Errorf("Reverse(%v) => err(%v), want err(%v)", c.input, err, c.err)
		}
	}
}

func TestUpdateMember(t *testing.T) {
	repo, mock := getRepoAndMock(t)

	dbOutput := &models.Member{
		ID:        1,
		SlackID:   "slack01",
		Name:      "slack.01",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	input := &models.UpdateMemberInput{
		ID:      dbOutput.ID,
		SlackID: &dbOutput.SlackID,
		Name:    &dbOutput.Name,
	}

	// Prepare sqlmock
	mock.ExpectBegin()
	rows := sqlmock.NewRows([]string{"id", "slack_id", "name", "created_at", "updated_at"}).
		AddRow(dbOutput.ID, dbOutput.SlackID, dbOutput.Name, dbOutput.CreatedAt, dbOutput.UpdatedAt)
	sql := regexp.QuoteMeta("SELECT * FROM `members`")
	mock.ExpectQuery(sql).WithArgs(input.ID).WillReturnRows(rows)
	sql = regexp.QuoteMeta("UPDATE `members` SET `slack_id`=?,`name`=?,`created_at`=?,`updated_at`=? WHERE `id` = ?")
	mock.ExpectExec(sql).WithArgs(dbOutput.SlackID, dbOutput.Name, AnyTime{}, AnyTime{}, input.ID).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	// Start Test
	_, err := repo.UpdateMember(context.Background(), input)
	if err != nil {
		t.Fatal(err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestUpdateMember_Error(t *testing.T) {
	repo, _ := getRepoAndMock(t)

	inputInstance := &models.Member{
		Name: "TestUpdate_IdMustNotBeZero",
	}

	cases := []struct {
		input *models.UpdateMemberInput
		err   error
	}{
		{
			input: &models.UpdateMemberInput{
				Name: &inputInstance.Name,
			},
			err: ErrBadRequestIDMustNotBeZero,
		},
	}

	for _, c := range cases {
		if _, err := repo.UpdateMember(context.Background(), c.input); err != c.err {
			t.Errorf("Reverse(%v) => err(%v), want err(%v)", c.input, err, c.err)
		}
	}
}

func TestDeleteMemberByID(t *testing.T) {
	repo, mock := getRepoAndMock(t)

	var input uint = 1

	// Prepare sqlmock
	sql := regexp.QuoteMeta("DELETE FROM `members` WHERE `members`.`id` = ?")
	mock.ExpectExec(sql).WithArgs(input).WillReturnResult(sqlmock.NewResult(1, 1))

	// Start Test
	output, err := repo.DeleteMemberByID(context.Background(), input)
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
func TestDeleteMemberByIDNoSuchEntity(t *testing.T) {
	repo, mock := getRepoAndMock(t)

	var input uint = 1

	// Prepare sqlmock
	sql := regexp.QuoteMeta("DELETE FROM `members` WHERE `members`.`id` = ?")
	mock.ExpectExec(sql).WithArgs(input).WillReturnResult(sqlmock.NewResult(0, 0))

	// Start Test
	output, err := repo.DeleteMemberByID(context.Background(), input)
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

func TestDeleteMemberByID_Error(t *testing.T) {
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
		output, err := repo.DeleteMemberByID(context.Background(), c.input)
		if err != c.err {
			t.Errorf("Reverse(%v) => err(%v), want err(%v)", c.input, err, c.err)
		}
		if output != c.output {
			t.Errorf("Reverse(%v) => err(%v), want err(%v)", c.input, output, c.output)
		}
	}
}
