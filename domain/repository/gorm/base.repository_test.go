package gorm

import(
	"testing"
	"encoding/json"
	"time"
	"context"
	"github.com/satori/go.uuid"
	"github.com/xxxmicro/base/log"
	"github.com/xxxmicro/base/database/gorm"
	"github.com/xxxmicro/base/domain/model"
	_gorm "github.com/jinzhu/gorm"
	"github.com/micro/go-micro/v2/config"
	"github.com/micro/go-micro/v2/config/source/memory"
	"github.com/stretchr/testify/assert"
)

func init() {
	log.Init("dev")
}

type User struct {
	ID	string				`json:"id"`
	Name string				`json:"name"`
	Age int 				`json:"age"`
	Ctime time.Time 		`json:"ctime"`
	Mtime time.Time 		`json:"mtime"`
	Dtime time.Time 		`json:"dtime"`
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
	config, err := getConfig()
	if err != nil {
		t.Fatal(err)
		return
	}

	db, err := getDB(config)
	if err != nil {
		t.Fatal(err)
		return
	}

	db.AutoMigrate(&User{})
	log.Info("创建数据表完毕")



	userRepo := NewBaseRepository(db)

	user1 := User{
		ID: uuid.NewV4().String(),
		Name: "吕布",
		Age: 28,
	}
	
	user2 := User{
		ID: uuid.NewV4().String(),
		Name: "貂蝉",
		Age: 21,
	}

	{
		err := userRepo.Create(context.Background(), user1)
		assert.NoError(t, err)

		err = userRepo.Create(context.Background(), user2)
		assert.NoError(t, err)

		log.Info("插入记录成功")
	}

	{
		user1.Name = "赵云"
		err := userRepo.Update(context.Background(), user1)
		assert.NoError(t, err)
		log.Info("更新记录成功")
	}

	{
		findUser := &User{ID: user1.ID }
		err := userRepo.FindOne(context.Background(), findUser, findUser)
		assert.NoError(t, err)
		t.Log(findUser)
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
	
		page, err := userRepo.Page(context.Background(), pageQuery, &User{})
		assert.NoError(t, err)
		assert.Equal(t, 1, page.Total)

		b, _ := json.Marshal(page)
		s := string(b)
		t.Log(s)
		log.Info("翻页查询正确")
	}

	{
		cursorQuery := &model.CursorQuery{
			Filters: map[string]interface{}{
			},
			CursorSort: &model.SortSpec{
				Property: "ctime",
			},
			Cursor: nil,
			Size: 10,
		}

		items := make([]*User, 0)
		extra, err := userRepo.Cursor(context.Background(), cursorQuery, &User{}, &items)
		assert.NoError(t, err)
		b, _ := json.Marshal(items)
		s := string(b)
		t.Log(s)
		t.Log(extra)
		log.Info("游标查询成功")
	}

	{
		err := userRepo.Delete(context.Background(), User{ID: user1.ID})
		assert.NoError(t, err)
		log.Info("删除记录成功")

		page, err := userRepo.Page(context.Background(), &model.PageQuery{
			Filters: map[string]interface{}{},
			PageSize: 10,
			PageNo: 1,
		}, &User{})
		assert.NoError(t, err)
		assert.Equal(t, 1, page.Total)

		err = userRepo.Delete(context.Background(), User{ID: user2.ID})
		assert.NoError(t, err)
		log.Info("删除记录成功")

		page, err = userRepo.Page(context.Background(), &model.PageQuery{
			Filters: map[string]interface{}{},
			PageSize: 10,
			PageNo: 1,
		}, &User{})
		assert.NoError(t, err)
		assert.Equal(t, 0, page.Total)
	
		log.Info("翻页核对成功")
	}
}
