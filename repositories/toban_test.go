package repositories

import (
	"context"
	"database/sql/driver"
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

func getDBMock() (*gorm.DB, sqlmock.Sqlmock, error) {
	db, mock, err := sqlmock.New()
	if err != nil {
		return nil, nil, err
	}
	gormDB, err := gorm.Open(
		mysql.Dialector{Config: &mysql.Config{
			DriverName: "mysql", Conn: db, SkipInitializeWithVersion: true,
		}}, &gorm.Config{})
	if err != nil {
		return nil, nil, err
	}
	return gormDB, mock, nil
}

func TestNewTodoRepository(t *testing.T) {
	// DBモック用意
	db, mock, err := getDBMock()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	sql := regexp.QuoteMeta("CREATE TABLE `tobans`")
	mock.ExpectExec(sql).WillReturnResult(sqlmock.NewResult(1, 1))
	NewTobanRepository(db)
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestNewTodoRepositoryNoMigrate(t *testing.T) {
	// DBモック用意
	db, mock, err := getDBMock()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	NewTobanRepositoryNoMigrate(db)
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestGet(t *testing.T) {
	ctx := context.Background()

	// DBモック用意
	db, mock, err := getDBMock()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	input := &models.Toban{
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
	rows := sqlmock.NewRows([]string{"id", "name", "description", "interval", "deadline_hour", "deadline_week_day", "deadline_week", "enabled", "toban_member_sequence", "created_at", "updated_at"}).
		AddRow(input.ID, input.Name, input.Description, input.Interval, input.DeadlineHour, input.DeadlineWeekDay, input.DeadlineWeek, input.Enabled, input.TobanMemberSequence, input.CreatedAt, input.UpdatedAt)
	sql := regexp.QuoteMeta("SELECT * FROM `tobans`")
	mock.ExpectQuery(sql).WithArgs(input.ID).WillReturnRows(rows)

	repo := NewTobanRepositoryNoMigrate(db)

	output, err := repo.Get(ctx, input.ID)
	if err != nil {
		t.Fatal(err)
	}
	if diff := cmp.Diff(input, output); diff != "" {
		t.Errorf("differs: (-got +want)\n%s", diff)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestGet_ErrNoSuchEntity(t *testing.T) {
	ctx := context.Background()

	// DBモック用意
	db, mock, err := getDBMock()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	input := &models.Toban{
		ID: 1,
	}
	rows := sqlmock.NewRows([]string{"id", "name", "description", "interval", "deadline_hour", "deadline_week_day", "deadline_week", "enabled", "toban_member_sequence", "created_at", "updated_at"})
	sql := regexp.QuoteMeta("SELECT * FROM `tobans`")
	mock.ExpectQuery(sql).WithArgs(input.ID).WillReturnRows(rows)

	repo := NewTobanRepositoryNoMigrate(db)

	if _, err := repo.Get(ctx, input.ID); err != ErrNoSuchEntity {
		t.Fatalf("It doesn't return an error when no such entity. %v", err)
	}
}

func TestCreate(t *testing.T) {
	ctx := context.Background()

	// DBモック用意
	db, mock, err := getDBMock()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

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
	sql := regexp.QuoteMeta("INSERT INTO `tobans` (`name`,`description`,`interval`,`deadline_hour`,`deadline_week_day`,`deadline_week`,`enabled`,`toban_member_sequence`,`created_at`,`updated_at`) VALUES (?,?,?,?,?,?,?,?,?,?)")
	mock.ExpectExec(sql).WithArgs(input.Name, input.Description, input.Interval, input.DeadlineHour, input.DeadlineWeekDay, input.DeadlineWeek, input.Enabled, input.TobanMemberSequence, AnyTime{}, AnyTime{}).WillReturnResult(sqlmock.NewResult(1, 1))

	repo := NewTobanRepositoryNoMigrate(db)
	_, err = repo.Create(ctx, input)
	if err != nil {
		t.Fatal(err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestCreate_IdMustBeZero(t *testing.T) {
	ctx := context.Background()

	// DBモック用意
	db, _, err := getDBMock()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	input := &models.Toban{
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
	repo := NewTobanRepositoryNoMigrate(db)

	if _, err := repo.Create(ctx, input); err != ErrBadRequestIdMustBeZero {
		t.Fatalf("It doesn't return an error when the ID is 0. %v", err)
	}
}

func TestUpdate(t *testing.T) {
	ctx := context.Background()

	// DBモック用意
	db, mock, err := getDBMock()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	input := &models.Toban{
		ID:                  1,
		Name:                "掃除機",
		Description:         "desc",
		Interval:            "DAILY",
		DeadlineHour:        23,
		DeadlineWeekDay:     "SUNDAY",
		DeadlineWeek:        0,
		Enabled:             true,
		TobanMemberSequence: 0,
	}
	sql := regexp.QuoteMeta("UPDATE `tobans` SET `name`=?,`description`=?,`interval`=?,`deadline_hour`=?,`deadline_week_day`=?,`deadline_week`=?,`enabled`=?,`toban_member_sequence`=?,`created_at`=?,`updated_at`=? WHERE `id` = ?")
	mock.ExpectExec(sql).WithArgs(input.Name, input.Description, input.Interval, input.DeadlineHour, input.DeadlineWeekDay, input.DeadlineWeek, input.Enabled, input.TobanMemberSequence, AnyTime{}, AnyTime{}, input.ID).WillReturnResult(sqlmock.NewResult(1, 1))

	repo := NewTobanRepositoryNoMigrate(db)
	_, err = repo.Update(ctx, input)
	if err != nil {
		t.Fatal(err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
