package gorm

import(
	"time"
	"errors"
	_ "github.com/go-sql-driver/mysql"
	_gorm "github.com/jinzhu/gorm"
	"github.com/micro/go-micro/v2/config"
)

func NewDbProvider(config config.Config) (*_gorm.DB, error) {
	driver := config.Get("db", "driver").String()
	connectionString := config.Get("db", "connectionString").String()

	if len(driver) == 0 {
		return nil, errors.New("driver is empty")
	}

	if len(connectionString) == 0 {
		return nil, errors.New("connectionString is empty")
	}

	db, err := _gorm.Open(driver, connectionString)
	if err != nil {
		panic(err)
	}

	// defer db.Close()
	db.LogMode(true)
	db.DB().SetMaxIdleConns(10)
	db.DB().SetConnMaxLifetime(3 * time.Minute)

	AddGormCallbacks(db)

	return db, nil
}