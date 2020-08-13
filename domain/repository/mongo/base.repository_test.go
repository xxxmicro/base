package mongo

import(
	"testing"
	"encoding/json"
	"time"
	"context"
	"gopkg.in/mgo.v2/bson"
	"github.com/xxxmicro/base/log"
	"github.com/xxxmicro/base/domain/model"
	"github.com/xxxmicro/base/database/mongo"
	"github.com/micro/go-micro/v2/config"
	"github.com/micro/go-micro/v2/config/source/memory"
	"github.com/stretchr/testify/assert"
)

func init() {
	log.Init("dev")
}

type User struct {
	ID	bson.ObjectId		`bson:"_id"`
	Name string				`bson:"name"`
	Age int 				`bson:"age"`
	Ctime time.Time 		`bson:"ctime"`
	Mtime time.Time 		`bson:"mtime"`
	Dtime time.Time 		`bson:"dtime"`
}

func (u *User) Unique() interface{} {
	return bson.M{
		"_id": u.ID,
	}
}



func getConfig() (config.Config, error) {
	config, err := config.NewConfig()
	if err != nil {
		return nil, err
	}

	data := []byte(`{
		"mongo": {
			"addrs": [
				"localhost:27017"
			],
			"database": "uim"
		}
	}`)
	source := memory.NewSource(memory.WithJSON(data))

	err = config.Load(source)
	if err != nil {
		return nil, err 
	}

	return config, nil
}

func getDB(config config.Config) (*mongo.DB, error) {
	db, err := mongo.NewMongoProvider(config)
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

	// db.AutoMigrate(&User{})
	// log.Info("创建数据表完毕")


	userRepo := NewBaseRepository(db)

	user1 := &User{
		ID: bson.NewObjectId(),
		Name: "吕布",
		Age: 28,
	}
	
	user2 := &User{
		ID: bson.NewObjectId(),
		Name: "貂蝉",
		Age: 21,
	}

	{
		err := userRepo.Create(context.Background(), user1)
		assert.NoError(t, err)
		if err != nil {
			log.Fatal("插入记录失败")
			return
		}

		err = userRepo.Create(context.Background(), user2)
		assert.NoError(t, err)
		if err != nil {
			log.Fatal("插入记录失败")
			return
		}

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
		err := userRepo.FindOne(context.Background(), findUser)
		assert.NoError(t, err)
		if err != nil {
			log.Info("查找记录失败")
			t.Fatal(err)
			return
		}
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
	
		items := make([]*User, 0)
		total, pageCount, err := userRepo.Page(context.Background(), pageQuery, &User{}, &items)
		assert.NoError(t, err)
		assert.Equal(t, 1, total)
		assert.Equal(t, 1, pageCount)

		b, _ := json.Marshal(items)
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
		err := userRepo.Delete(context.Background(), &User{ID: user1.ID})
		assert.NoError(t, err)
		if err != nil {
			log.Fatal("删除记录失败")			
			return
		}
		log.Info("删除记录成功")

		items := make([]*User, 0)
		total, pageCount, err := userRepo.Page(context.Background(), &model.PageQuery{
			Filters: map[string]interface{}{},
			PageSize: 10,
			PageNo: 1,
		}, &User{}, &items)
		assert.NoError(t, err)
		assert.Equal(t, 1, total)
		assert.Equal(t, 1, len(items))

		err = userRepo.Delete(context.Background(), &User{ID: user2.ID})
		assert.NoError(t, err)
		if err != nil {
			log.Fatal("删除记录失败")			
			return
		}
		log.Info("删除记录成功")

		items = make([]*User, 0)
		total, pageCount, err = userRepo.Page(context.Background(), &model.PageQuery{
			Filters: map[string]interface{}{},
			PageSize: 10,
			PageNo: 1,
		}, &User{}, &items)
		assert.NoError(t, err)
		assert.Equal(t, 0, total)
		assert.Equal(t, 0, pageCount)
	
		log.Info("翻页核对成功")
	}
}
