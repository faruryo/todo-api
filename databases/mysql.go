package databases

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func GetDbByDsn(dsn string, logLevel LogLevel) (*gorm.DB, error) {
	return GetDbByDialector(mysql.Open(dsn), logLevel)
}

func GetDbByDialector(dialector gorm.Dialector, logLevel LogLevel) (*gorm.DB, error) {
	return gorm.Open(
		dialector,
		&gorm.Config{
			DisableAutomaticPing:   true,
			SkipDefaultTransaction: true,
			Logger:                 logger.Default.LogMode(logLevel.logLevel()),
		},
	)
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
