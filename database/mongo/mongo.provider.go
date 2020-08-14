package mongo

import(
	"time"
	"errors"
	"gopkg.in/mgo.v2"
	"github.com/micro/go-micro/v2/config"
)

// 主要为了 apollo 动态替换内部实现
type DB struct {
	Name string
	Session *mgo.Session
}

func NewMongoProvider(config config.Config) (*DB, error) {
	addrs := config.Get("mongo", "addrs").StringSlice(nil)
	if len(addrs) == 0 {
		return nil, errors.New("addrs must be set")
	}

	database := config.Get("mongo", "database").String("")
	if len(database) == 0 {
		return nil, errors.New("database must be set")
	}

	poolLimit := config.Get("mongo", "pool_limit").Int(20)
	
	timeout := config.Get("mongo", "timeout").Duration(time.Second * 5)
	mode := config.Get("mongo", "mode").Int(0)

	dialInfo := &mgo.DialInfo{
		Addrs: addrs,
		Database: database,
		PoolLimit: poolLimit,
		Timeout: timeout,
	}

	replicaSetName := config.Get("mongo", "replica_set_name").String("")
	if len(replicaSetName) > 0 {
		dialInfo.ReplicaSetName = replicaSetName
	}

	username := config.Get("mongo", "username").String("")
	if len(username) > 0 {
		dialInfo.Username = username
	}

	password := config.Get("mongo", "password").String("")
	if len(password) > 0 {
		dialInfo.Password = password
	}

	source := config.Get("mongo", "source").String("")
	if len(source) >= 0 {
		dialInfo.Source = source
	}

	globalSession, err := mgo.DialWithInfo(dialInfo)
	if err != nil {
		return nil, err
	}

	globalSession.SetSocketTimeout(time.Second * 10)
	globalSession.SetSyncTimeout(time.Second * 20)
	
	if mode <= 0 {
		mode = int(mgo.Primary)
	}
	globalSession.SetMode(mgo.Mode(mode), true)


	db := &DB{ Name: database, Session: globalSession }	

	go watchConfigChange(config, db)
	
	return db, nil
}

func watchConfigChange(config config.Config, db *DB) {
	// TODO
}
