package mongo

import(
	"errors"
	"fmt"
	"gopkg.in/mgo.v2"
)

type DBFunc func(*mgo.Collection) error

func Execute(globalSession *mgo.Session, database string, collection string, fn DBFunc) error {
	session := globalSession.Clone()
	defer session.Close()
	db := session.DB(database)
	c := db.C(collection)
	if c == nil {
		return errors.New(fmt.Sprintf("ERR_DB_UNKNOWN_TABLE %s", collection))
	}

	if err := fn(c); err != nil {
		return err
	}

	return nil
}



