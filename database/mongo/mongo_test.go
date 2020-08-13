package mongo_test

import(
	"testing"
	"gopkg.in/mgo.v2/bson"
	"github.com/xxxmicro/base/database/mongo"
	"github.com/micro/go-micro/v2/config"
	"github.com/micro/go-micro/v2/config/source/memory"
)

type User struct {
	ID bson.ObjectId `bson:"_id,omitempty"`
	Name string `bson:"name"`
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

func TestMongo(t *testing.T) {
	config, err := getConfig()
	if err != nil {
		t.Fatal(err)
		return
	}

	globalSession, err := mongo.NewMongoProvider(config)
	if err != nil {
		t.Fatal(err)
		return 
	}

	user := User{
		ID: bson.NewObjectId(),
		Name: "alice",
	}

	{
		session := globalSession.Clone()
		defer session.Close()

		err := session.DB("uim").C("users").Insert(user)
		if err != nil {
			t.Fatal(err)
			return
		}
	}
	
	{
		session := globalSession.Clone()
		defer session.Close()

		findUser := &User{}
		condition := bson.M{
			"_id": user.ID,
		}
	
		err := session.DB("uim").C("users").Find(condition).One(findUser)
		if err != nil {
			t.Fatal(err)
			return
		}

		t.Log(findUser)
	}

	return
}
