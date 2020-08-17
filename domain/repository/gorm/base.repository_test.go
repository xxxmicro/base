package gorm

import (
	"context"
	"encoding/json"
	"fmt"
	_gorm "github.com/jinzhu/gorm"
	"github.com/micro/go-micro/v2/config"
	"github.com/micro/go-micro/v2/config/source/memory"
	"github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"github.com/xxxmicro/base/database/gorm"
	"github.com/xxxmicro/base/domain/model"
	"testing"
	"time"
)

type User struct {
	ID	string				`json:"id" gorm:"primary_key"`
	Name string				`json:"name"`
	Age int 				`json:"age"`
	Ctime time.Time 		`json:"ctime" gorm:"update_time_stamp"`
	Mtime time.Time 		`json:"mtime" gorm:"update_time_stamp"`
	Dtime *time.Time 		`json:"dtime"`
}

func (u *User) BeforeCreate(scope *_gorm.Scope) error {
	scope.SetColumn("id", uuid.NewV4().String())
	return nil
}

func (u *User) Unique() interface{} {
	return map[string]interface{}{
		"id": u.ID,
	}
}

func getConfig() (config.Config, error) {
	config, err := config.NewConfig()
	if err != nil {
		return nil, err
	}

	data := []byte(`{
		"db": {
			"driver": "mysql",
			"connection_string": "root:root@tcp(localhost:3306)/uim?charset=utf8mb4&parseTime=True&loc=Local"
		}
	}`)
	source := memory.NewSource(memory.WithJSON(data))

	err = config.Load(source)
	if err != nil {
		return nil, err 
	}

	return config, nil
}

func getDB(config config.Config) (*_gorm.DB, error) {
	db, err := gorm.NewDbProvider(config)
	if err != nil {
		log.Panic("数据库连接失败")
		return nil, err
	}
	
	return db, nil
}

func TestCrud(t *testing.T) {
	assert := assert.New(t)

	config, err := getConfig()
	if err != nil {
		t.Fatal(err)
	}

	db, err := getDB(config)
	if err != nil {
		t.Fatal(err)
	}

	db.AutoMigrate(&User{})
	log.Info("创建数据表完毕")


	userRepo := NewBaseRepository(db)

	user1 := &User{
		Name: "吕布",
		Age: 28,
	}
	
	user2 := &User{
		Name: "貂蝉",
		Age: 21,
	}

	{
		err := userRepo.Create(context.Background(), user1)
		if assert.Error(err) {
			t.Fatal(err)
		}

		time.Sleep(time.Second * 3)

		err = userRepo.Create(context.Background(), user2)
		if assert.Error(err) {
			t.Fatal(err)
		}

		log.Info("插入记录成功")
	}

	user3 := &User{
		Name: "关羽",
		Age: 38,
	}
	{
		change, err := userRepo.Upsert(context.Background(), user3)
		assert.Error(err)
		t.Logf("change: %v", change)
	}

	{
		user1.Name = "赵云"
		err := userRepo.Update(context.Background(), user1, user1)
		if assert.Error(err) {
			t.Fatal(err)
		}

		data := map[string]interface{}{
			"name": "孙悟空",
			"age": 0,
		}

		err = userRepo.Update(context.Background(), &User{}, data)
		if assert.NoError(err) {
			t.Fatal(err)
		}

		// 如果这么更新， age将不会被设置, &User{ Name: "sunwukong", Age: 0}
		err = userRepo.Update(context.Background(), user2, data)
		if assert.Error(err) {
			t.Fatal(err)
		}
		log.Info("选择更新成功")
	}

	{
		findUser := &User{ ID: user1.ID }
		err := userRepo.FindOne(context.Background(), findUser)
		if assert.Error(err) {
			t.Fatal(err)
		}
		log.Info("找到对应记录")
	}

	{
		pageQuery := &model.PageQuery{
			Filters: map[string]interface{}{
				"name": "赵云",
				"age": map[string]interface{}{
					"GT": 22,
				},
			},
			PageSize: 10,
			PageNo: 1,
		}
	
		items := make([]*User, 0)
		total, pageCount, err := userRepo.Page(context.Background(), &User{}, pageQuery, &items)
		if assert.Error(err) {
			t.Fatal(err)
		}

		if assert.Equal(1, total) {
			log.Info("翻页查询正确")
		} else {
			log.Info(fmt.Sprintf("翻页查询错误, 期望1条记录，实际返回%d条", total))
		}

		if assert.Equal(1, pageCount) {
			log.Info("翻页查询正确")
		} else {
			log.Info(fmt.Sprintf("翻页查询错误 期望1页, 实际返回%d页", pageCount))
		}

		b, _ := json.Marshal(items)
		s := string(b)
		t.Log(s)
	}

	{
		h, _ := time.ParseDuration("1s")
		t1 := user1.Ctime.Add(h)
		cursor := t1.UnixNano() / 1e6

		cursorQuery := &model.CursorQuery{
			Filters: map[string]interface{}{
			},
			CursorSort: &model.SortSpec{
				Property: "ctime",
			},
			Cursor: cursor,
			Size: 10,
		}

		items := make([]*User, 0)
		extra, err := userRepo.Cursor(context.Background(), cursorQuery, &User{}, &items)
		if assert.Error(err) {
			t.Fatal(err)
		}

		if assert.Equal(2, len(items)) {
			log.Info("游标查询正确")
		} else {
			log.Info(fmt.Sprintf("游标查询错误 期望1条, 实际返回%d条", len(items)))
		}

		b, _ := json.Marshal(items)
		s := string(b)
		t.Log(s)
		t.Log(extra)
		log.Info("游标查询成功")
	}

	{
		err := userRepo.Delete(context.Background(), &User{ID: user1.ID})
		assert.NoError(err)
		log.Info("删除记录成功")

		items := make([]*User, 0)
		total, pageCount, err := userRepo.Page(context.Background(), &User{}, &model.PageQuery{
			Filters: map[string]interface{}{},
			PageSize: 10,
			PageNo: 1,
		}, &items)
		if assert.Error(err) {
			t.Fatal(err)
		}

		assert.Equal(1, total)
		assert.Equal(1, len(items))

		err = userRepo.Delete(context.Background(), &User{ID: user2.ID})
		err = userRepo.Delete(context.Background(), &User{ID: user3.ID})
		log.Info("删除记录成功")

		items = make([]*User, 0)
		total, pageCount, err = userRepo.Page(context.Background(), &User{}, &model.PageQuery{
			Filters: map[string]interface{}{},
			PageSize: 10,
			PageNo: 1,
		}, &items)
		assert.NoError(err)
		assert.Equal(t, 0, total)
		assert.Equal(t, 0, pageCount)
	
		log.Info("翻页核对成功")
	}
}
