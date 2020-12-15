package repository

import (
	"context"
	"database/sql/driver"
	"os"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/faruryo/toban-api/models"
	"github.com/google/go-cmp/cmp"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type AnyTime struct{}

func (a AnyTime) Match(v driver.Value) bool {
	_, ok := v.(time.Time)
	return ok
}

var ctx context.Context

func TestMain(m *testing.M) {

	ctx = context.Background()

	code := m.Run()

	os.Exit(code)
}

func getDBMock() (*gorm.DB, sqlmock.Sqlmock, error) {
	db, mock, err := sqlmock.New()
	if err != nil {
		return nil, nil, err
	}
	logLevel := Silent
	if testing.Verbose() {
		logLevel = Info
	}
	gormDB, err := GetDbByDialector(
		mysql.Dialector{Config: &mysql.Config{
			DriverName:                "mysql",
			Conn:                      db,
			SkipInitializeWithVersion: true,
		}},
		logLevel,
	)
	if err != nil {
		return nil, nil, err
	}
	return gormDB, mock, nil
}

func TestNewTobanRepository(t *testing.T) {
	// DBモック用意
	db, mock, err := getDBMock()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	// sqlmock準備
	sql := regexp.QuoteMeta("CREATE TABLE `tobans`")
	mock.ExpectExec(sql).WillReturnResult(sqlmock.NewResult(1, 1))

	// Test開始
	NewTobanRepository(db)
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
func getRepoAndMock(t *testing.T) (TobanRepository, sqlmock.Sqlmock) {
	t.Helper()
	// DBモック用意
	db, mock, err := getDBMock()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	// Test開始
	repo := NewTobanRepositoryNoMigrate(db)
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	return repo, mock
}

func TestNewTobanRepositoryNoMigrate(t *testing.T) {
	getRepoAndMock(t)
}

func TestGet(t *testing.T) {
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
	output, err := repo.Get(ctx, dbOutput.ID)
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

func TestGet_Error(t *testing.T) {
	repo, mock := getRepoAndMock(t)

	input := &models.Toban{
		ID: 1,
	}

	// sqlmock準備
	rows := sqlmock.NewRows([]string{"id", "name", "description", "interval", "deadline_hour", "deadline_week_day", "deadline_week", "enabled", "toban_member_sequence", "created_at", "updated_at"})
	sql := regexp.QuoteMeta("SELECT * FROM `tobans`")
	mock.ExpectQuery(sql).WithArgs(input.ID).WillReturnRows(rows)

	// Test開始
	if _, err := repo.Get(ctx, input.ID); err != ErrNoSuchEntity {
		t.Fatalf("it doesn't return an error when no such entity. %v", err)
	}
}

func TestGetAll(t *testing.T) {
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
	output, err := repo.GetAll(ctx)
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

func TestCreate(t *testing.T) {
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
	_, err := repo.Create(ctx, input)
	if err != nil {
		t.Fatal(err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestCreate_Error(t *testing.T) {
	repo, _ := getRepoAndMock(t)

	cases := []struct {
		input *models.Toban
		err   error
	}{
		{
			input: &models.Toban{ID: 1},
			err:   ErrBadRequestIdMustBeZero,
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
		if _, err := repo.Create(ctx, c.input); err != c.err {
			t.Errorf("Reverse(%v) => err(%v), want err(%v)", c.input, err, c.err)
		}
	}
}

func TestUpdate(t *testing.T) {
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
	_, err := repo.Update(ctx, input)
	if err != nil {
		t.Fatal(err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestUpdate_Error(t *testing.T) {
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
			err: ErrBadRequestIdMustNotBeZero,
		},
	}

	for _, c := range cases {
		if _, err := repo.Update(ctx, c.input); err != c.err {
			t.Errorf("Reverse(%v) => err(%v), want err(%v)", c.input, err, c.err)
		}
	}
}

func TestDelete(t *testing.T) {
	repo, mock := getRepoAndMock(t)

	var input uint = 1

	// sqlmock準備
	sql := regexp.QuoteMeta("DELETE FROM `tobans` WHERE `tobans`.`id` = ?")
	mock.ExpectExec(sql).WithArgs(input).WillReturnResult(sqlmock.NewResult(1, 1))

	// Test開始
	output, err := repo.Delete(ctx, input)
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
func TestDeleteNoSuchEntity(t *testing.T) {
	repo, mock := getRepoAndMock(t)

	var input uint = 1

	// sqlmock準備
	sql := regexp.QuoteMeta("DELETE FROM `tobans` WHERE `tobans`.`id` = ?")
	mock.ExpectExec(sql).WithArgs(input).WillReturnResult(sqlmock.NewResult(0, 0))

	// Test開始
	output, err := repo.Delete(ctx, input)
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

func TestDelete_Error(t *testing.T) {
	repo, _ := getRepoAndMock(t)

	cases := []struct {
		input  uint
		output bool
		err    error
	}{
		{
			input:  0,
			output: false,
			err:    ErrBadRequestIdMustNotBeZero,
		},
	}

	for _, c := range cases {
		output, err := repo.Delete(ctx, c.input)
		if err != c.err {
			t.Errorf("Reverse(%v) => err(%v), want err(%v)", c.input, err, c.err)
		}
		if output != c.output {
			t.Errorf("Reverse(%v) => err(%v), want err(%v)", c.input, output, c.output)
		}
	}
}
