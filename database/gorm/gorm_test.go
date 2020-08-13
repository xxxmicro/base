package gorm_test

import(
	"testing"
	"context"
	"github.com/xxxmicro/base/database/gorm"
	"github.com/xxxmicro/base/opentracing/jaeger"
	"github.com/micro/go-micro/v2/config"
	"github.com/micro/go-micro/v2/config/source/memory"
	opentracing "github.com/opentracing/opentracing-go"
)

type User struct {
	Name string `json:"name"`
}

func TestGorm(t *testing.T) {
	config, err := config.NewConfig()
	if err != nil {
		t.Fatal(err)
		return 
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
		t.Fatal(err)
		return 
	}

	tracer, err := jaeger.NewTracerProvider(config)
	if err != nil {
		t.Fatal(err)
		return 
	}
	t.Log(tracer)

	db, err := gorm.NewDbProvider(config)
	if err != nil {
		t.Fatal(err)
		return 
	}

	user := User{
		Name: "alice",
	} 


	span, ctx := opentracing.StartSpanFromContext(context.Background(), "handler")
	defer span.Finish()

	db = gorm.SetSpanToGorm(ctx, db)

	err = db.Table("users").Create(user).Error
	if err != nil {
		t.Fatal(err)
		return 
	}

	err = db.Table("users").Take(&user).Error
	if err != nil {
		t.Fatal(err)
		return 
	}

	return
}
