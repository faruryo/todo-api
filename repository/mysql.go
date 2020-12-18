package repository

import (
	"fmt"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DBConnectionOptions is a structure that summarizes the options required when connecting to a DB.
type DBConnectionOptions struct {
	User     string
	Password string
	Database string
	Host     string
}

// GetDb creates a gorm.DB instance.
func GetDb(opt DBConnectionOptions, logLevel LogLevel) (*gorm.DB, error) {
	dsn := assembleDSN(opt)
	return getDbByDialector(mysql.Open(dsn), logLevel)
}

func assembleDSN(opt DBConnectionOptions) string {
	return fmt.Sprintf("%s:%s@tcp(%s:3306)/%s?charset=utf8mb4&parseTime=true", opt.User, opt.Password, opt.Host, opt.Database)
}

func getDbByDialector(dialector gorm.Dialector, logLevel LogLevel) (*gorm.DB, error) {
	db, err := gorm.Open(
		dialector,
		&gorm.Config{
			DisableAutomaticPing:   true,
			SkipDefaultTransaction: true,
			Logger:                 logger.Default.LogMode(logLevel.logLevel()),
		},
	)
	if err != nil {
		log.Printf("failed to connect database: %v", err)
		log.Printf("dialector: %v", dialector)
		return nil, err
	}

	return db, nil
}

type LogLevel int

const (
	Silent LogLevel = iota
	Error  LogLevel = iota
	Warn   LogLevel = iota
	Info   LogLevel = iota
)

func (c LogLevel) logLevel() logger.LogLevel {
	switch c {
	case Silent:
		return logger.Silent
	case Error:
		return logger.Error
	case Warn:
		return logger.Warn
	case Info:
		return logger.Info
	default:
		return logger.Info
	}
}
