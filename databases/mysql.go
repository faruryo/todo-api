package databases

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func GetDbByDsn(dsn string, logLevel LogLevel) (*gorm.DB, error) {
	return GetDbByDialector(mysql.Open(dsn), logLevel)
}

func GetDbByEnv(logLevel LogLevel) (*gorm.DB, error) {
	var (
		user     = os.Getenv("MYSQL_USER")
		password = os.Getenv("MYSQL_PASSWORD")
		database = os.Getenv("MYSQL_DATABASE")
		host     = os.Getenv("MYSQL_HOST")
	)
	dsn := fmt.Sprintf("%s:%s@tcp(%s:3306)/%s?charset=utf8mb4&parseTime=true", user, password, host, database)
	return GetDbByDialector(mysql.Open(dsn), logLevel)
}

func GetDbByDialector(dialector gorm.Dialector, logLevel LogLevel) (*gorm.DB, error) {
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
