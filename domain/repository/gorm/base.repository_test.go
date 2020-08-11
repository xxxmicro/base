package gorm

import(
	"testing"
	"encoding/json"
	"time"
	"context"
	"github.com/satori/go.uuid"
	"github.com/xxxmicro/base/database/gorm"
	"github.com/xxxmicro/base/domain/model"
	_gorm "github.com/jinzhu/gorm"
	"github.com/micro/go-micro/v2/config"
	"github.com/micro/go-micro/v2/config/source/memory"
)

type User struct {
	ID	string				`json:"id"`
	Name string				`json:"name"`
	Age int 					`json:"age"`
	Ctime time.Time 	`json:"ctime"`
	Mtime time.Time 	`json:"mtime"`
	Dtime time.Time 	`json:"dtime"`
}

func (user *User) GetTable() string {
	return "users"
}

func (user *User) GenerateID() {
	user.ID = uuid.NewV4().String()
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
		return nil, err
	}
	
	return db, nil
}

func TestBuildQuery(t *testing.T) {
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

	m := &User{}

	userRepo := NewBaseRepository(db)

	pageQuery := &model.PageQuery{
		Filters: map[string]interface{}{
			"name": "terry",
			"age": map[string]interface{}{
				"GT": 22,
			},
		},
		PageSize: 10,
		PageNo: 1,
	}

	page, err := userRepo.Page(context.Background(), pageQuery, m)
	if err != nil {
		t.Fatal(err)
		return
	}

	b, _ := json.Marshal(page)
	s := string(b)

	t.Log(s)
}