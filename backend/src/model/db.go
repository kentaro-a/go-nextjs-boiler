package model

import (
	"app/config"
	"fmt"
	"log"
	"os"

	"github.com/pkg/errors"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func GetDSN() string {
	c := config.Get()
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", c.DB.User, c.DB.Password, c.DB.Host, c.DB.Port, c.DB.DBName)
}

func NewDB() (*gorm.DB, error) {
	logger_config := logger.Config{
		IgnoreRecordNotFoundError: true,
		Colorful:                  true,
		LogLevel:                  logger.Silent,
	}
	newLogger := logger.New(log.New(os.Stdout, "\r\n", log.LstdFlags), logger_config)

	_db, err := gorm.Open(mysql.Open(GetDSN()), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return _db, nil
}

// func (_db *DB) NewTransaction() *DB {
// 	tx := &DB{_db.Begin()}
// 	return tx
// }

// func (_db *DB) Close() error {
// 	sqldb, err := _db.DB.DB()
// 	if err != nil {
// 		return errors.WithStack(err)
// 	}
// 	return sqldb.Close()
// }
