package gorm

import(
	"time"
	"errors"
	 _ "github.com/go-sql-driver/mysql"
	"github.com/micro/go-micro/v2/config"
	"github.com/jinzhu/gorm"
)

func NewDbProvider(config config.Config) (*gorm.DB, error) {
	driver := config.Get("db", "driver").String("")
	connectionString := config.Get("db", "connection_string").String("")

	if len(driver) == 0 {
		return nil, errors.New("driver is empty")
	}

	if len(connectionString) == 0 {
		return nil, errors.New("connection_string is empty")
	}

	db, err := gorm.Open(driver, connectionString)
	if err != nil {
		return nil, err
	}

	// defer db.Close()
	db.LogMode(true)
	db.DB().SetMaxIdleConns(10)
	db.DB().SetConnMaxLifetime(3 * time.Minute)

	AddGormCallbacks(db)

	return db, nil
}
