package repository

import (
	"context"
	"database/sql/driver"
	"os"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
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
func getRepoAndMock(t *testing.T) (Repository, sqlmock.Sqlmock) {
	t.Helper()
	// DBモック用意
	db, mock, err := getDBMock()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	// Test開始
	repo := NewRepositoryNoMigrate(db)
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	return repo, mock
}
