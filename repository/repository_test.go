package repository

import (
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestNewRepository(t *testing.T) {
	// Prepare DB mock
	db, mock, err := getDBMock()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	// Prepare sqlmock
	sql := regexp.QuoteMeta("CREATE TABLE `tobans`")
	mock.ExpectExec(sql).WillReturnResult(sqlmock.NewResult(1, 1))
	sql = regexp.QuoteMeta("CREATE TABLE `members`")
	mock.ExpectExec(sql).WillReturnResult(sqlmock.NewResult(1, 1))

	// Start Test
	_, err = NewRepository(db)
	if err != nil {
		t.Errorf("failed to create repository.: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestNewRepositoryNoMigrate(t *testing.T) {
	getRepoAndMock(t)
}
