package gorm

import (
	"errors"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/micro/go-micro/v2/config"
	"github.com/xxxmicro/base/database/gorm/opentracing"
	"time"
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

	addAutoCallbacks(db)

	opentracing.AddGormCallbacks(db)

	return db, nil
}

func addAutoCallbacks(db *gorm.DB) {
	// 替换替换默认的钩子
	db.Callback().Create().Replace("gorm:update_time_stamp", updateTimeForCreateCallback)
	db.Callback().Update().Replace("gorm:update_time_stamp", updateTimeForUpdateCallback)
}

func updateTimeForCreateCallback(scope *gorm.Scope) {
	if !scope.HasError() {
		nowTime := time.Now()

		// 通过 scope.Fields() 获取所有字段，判断当前是否包含所需字段
		if createTimeField, ok := scope.FieldByName("Ctime"); ok {
			if createTimeField.IsBlank {	// 可判断该字段的值是否为空
				createTimeField.Set(nowTime)
			}
		}

		if modifyTimeField, ok := scope.FieldByName("Mtime"); ok {
			if modifyTimeField.IsBlank {
				modifyTimeField.Set(nowTime)
			}
		}
	}
}

func updateTimeForUpdateCallback(scope *gorm.Scope) {
	// scope.Get(...) 根据入参获取设置了字面值的参数，例如本文中是 gorm:update_column ，它会去查找含这个字面值的字段属性
	if _, ok := scope.Get("gorm:mtime"); !ok {
		// scope.SetColumn(...) 假设没有指定 update_column 的字段，我们默认在更新回调设置 ModifiedOn 的值
		scope.SetColumn("mtime", time.Now())
	}
}


