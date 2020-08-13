package mongo

import(
	"fmt"
	"context"
	"reflect"
	"gopkg.in/mgo.v2"
	"github.com/xxxmicro/base/database/mongo"
	"github.com/xxxmicro/base/domain/repository"
	"github.com/xxxmicro/base/domain/model"
	breflect "github.com/xxxmicro/base/reflect"
)

type baseRepository struct {
	db *mongo.DB
}

func NewBaseRepository(db *mongo.DB) repository.BaseRepository {
	return &baseRepository{ db }
}

func (r *baseRepository) Create(c context.Context, m model.Model) error {
	modelName := reflect.TypeOf(m).Elem().Name()
	collection := TheNamingStrategy.Table(modelName)
	fmt.Printf("xxx: %s %s\n", modelName, collection)
	return Execute(r.db.Session, r.db.Name, collection, func(c *mgo.Collection) error {
		return c.Insert(m)
	})
}

func (r *baseRepository) Update(c context.Context, m model.Model) error {
	collection := TheNamingStrategy.Table(reflect.TypeOf(m).Elem().Name())
	return Execute(r.db.Session, r.db.Name, collection, func(c *mgo.Collection) error {
		return c.Update(m.Unique(), m)
	})
}

func (r *baseRepository) FindOne(c context.Context, m model.Model) error {
	collection := TheNamingStrategy.Table(reflect.TypeOf(m).Elem().Name())
	return Execute(r.db.Session, r.db.Name, collection, func(c *mgo.Collection) error {
		return c.Find(m.Unique()).One(m)
	})
}

func (r *baseRepository) Delete(c context.Context, m model.Model) error {
	collection := TheNamingStrategy.Table(reflect.TypeOf(m).Elem().Name())
	return Execute(r.db.Session, r.db.Name, collection, func(c *mgo.Collection) error {
		return c.Remove(m.Unique())
	})
}

func (r *baseRepository) Page(c context.Context, query *model.PageQuery, m model.Model, resultPtr interface{}) (total int, pageCount int, err error){
	ms, err := breflect.GetStructInfo(m, nil)
	if err != nil {
		return
	}
	collection := TheNamingStrategy.Table(ms.Name)
	
	filters, err := buildQuery(ms, query.Filters)
	if err != nil {
		return
	}
	fmt.Print(filters)

	sorts, err := buildSort(ms, query.Sort)
	if err != nil {
		return
	}
	fmt.Print(sorts)
	
	err = Execute(r.db.Session, r.db.Name, collection, func(c *mgo.Collection) error {
		total, err := c.Find(filters).Count()
		if err != nil {
			return err
		}

		pageSize := query.PageSize
		if pageSize > 1000 {
			pageSize = 1000
		} else if pageSize <= 0 {
			pageSize = 20
		}

		pageNo := query.PageNo

		offset := (pageNo - 1) * pageSize

		pageCount = total / pageSize
		if total % pageSize != 0 {
			pageCount++
		}

		return c.Find(filters).Skip(offset).Limit(pageSize).Sort(sorts...).All(resultPtr)
	})

	return
}

func (r *baseRepository) Cursor(c context.Context, query *model.CursorQuery, m model.Model, resultPtr interface{}) (extra *model.CursorExtra, err error) {
	return
}

func (r *baseRepository) EnsureIndexes(m Indexed) (err error) {
	collection := TheNamingStrategy.Table(reflect.TypeOf(m).Elem().Name())

	err = Execute(r.db.Session, r.db.Name, collection, func(c *mgo.Collection) error {
		for _, i := range m.Indexes() {
			err = c.EnsureIndex(i)
			if err != nil {
				return err
			}
		}
		return nil
	})
	return
}
