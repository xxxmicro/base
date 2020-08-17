package mongo

import(
	"fmt"
	"testing"
	"encoding/json"
	"time"
	"context"
	"gopkg.in/mgo.v2/bson"
	"github.com/xxxmicro/base/domain/model"
	"github.com/xxxmicro/base/database/mongo"
	"github.com/micro/go-micro/v2/config"
	"github.com/micro/go-micro/v2/logger"
	"github.com/micro/go-micro/v2/config/source/memory"
	"github.com/stretchr/testify/assert"
)

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
	assert := assert.New(t)

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
		now := time.Now()
		user1.Ctime = now
		user2.Mtime = now
		err := userRepo.Create(context.Background(), user1)
		assert.NoError(err)
		if err != nil {
			log.Fatal("插入记录失败")
			return
		}

		time.Sleep(time.Second * 3)

		now = time.Now()
		user2.Ctime = now
		user2.Mtime = now
		err = userRepo.Create(context.Background(), user2)
		assert.NoError(err)
		if err != nil {
			log.Fatal("插入记录失败")
			return
		}

		log.Info("插入记录成功")
	}

	user3 := &User{
		ID: bson.NewObjectId(),
		Name: "关羽",
		Age: 38,
	}
	{
		change, err := userRepo.Upsert(context.Background(), user3)
		assert.NoError(err)
		t.Logf("change: %v", change)
	}

	{
		data := map[string]interface{}{
			"name": "赵云",
		}
		err := userRepo.Update(context.Background(), &User{}, data)
		if !assert.Error(err) {
			t.Fatal(err)
		}

		err = userRepo.Update(context.Background(), user1, data)
		if !assert.NoError(err) {
			t.Fatal(err)
		}
		log.Info("选择更新成功")
	}

	{
		findUser := &User{ID: user2.ID }
		err := userRepo.FindOne(context.Background(), findUser)
		assert.NoError(err)
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
		total, pageCount, err := userRepo.Page(context.Background(), &User{}, pageQuery, &items)
		if assert.Error(err) {
			t.Fatal(err)
		}

		if assert.Equal(1, total) {
			log.Info("翻页查询总数正确")
		} else {
			log.Info(fmt.Sprintf("翻页查询总数错误, 期望1, 返回%d", total))
		}

		if assert.Equal( 1, pageCount) {
			log.Info("翻页查询页数正确")
		} else {
			log.Info(fmt.Sprintf("翻页查询页数错误, 期望1, 返回%d", pageCount))
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

		b, _ := json.Marshal(extra)
		s := string(b)
		t.Log(s)
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
		assert.NoError(err)
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
		assert.Equal(0, total)
		assert.Equal(0, pageCount)
	
		log.Info("翻页核对成功")
	}
}
