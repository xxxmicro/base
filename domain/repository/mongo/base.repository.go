package mongo

import(
	"fmt"
	"context"
	"reflect"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
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
	ms, err := breflect.GetStructInfo(m, nil)
	if err != nil {
		return err
	}
	collection := TheNamingStrategy.Table(ms.Name)
	
	fmt.Printf("xxx: %s %s\n", ms.Name, collection)
	return Execute(r.db.Session, r.db.Name, collection, func(c *mgo.Collection) error {
		return c.Insert(m)
	})
}

func (r *baseRepository) Update(c context.Context, m model.Model) error {
	ms, err := breflect.GetStructInfo(m, nil)
	if err != nil {
		return err
	}
	collection := TheNamingStrategy.Table(ms.Name)

	return Execute(r.db.Session, r.db.Name, collection, func(c *mgo.Collection) error {
		return c.Update(m.Unique(), m)
	})
}

func (r *baseRepository) FindOne(c context.Context, m model.Model) error {
	ms, err := breflect.GetStructInfo(m, nil)
	if err != nil {
		return err
	}
	collection := TheNamingStrategy.Table(ms.Name)
	
	return Execute(r.db.Session, r.db.Name, collection, func(c *mgo.Collection) error {
		return c.Find(m.Unique()).One(m)
	})
}

func (r *baseRepository) Delete(c context.Context, m model.Model) error {
	ms, err := breflect.GetStructInfo(m, nil)
	if err != nil {
		return err
	}
	collection := TheNamingStrategy.Table(ms.Name)
	
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

	sorts, err := buildSort(ms, query.Sort)
	if err != nil {
		return
	}
	
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
	ms, err := breflect.GetStructInfo(m, nil)
	if err != nil {
		return
	}
	collection := TheNamingStrategy.Table(ms.Name)
	
	filters, err := buildQuery(ms, query.Filters)
	if err != nil {
		return
	}

	cursorFilter, sort, reverse, err := mongoCursorFilter(ms, query)
	if err != nil {
		return
	}	

	filters = bson.M{"$and": []bson.M{cursorFilter, filters}}

	size := query.Size
	if size > 1000 {
		size = 1000
	} else if size <= 0 {
		size = 20
	}

	err = Execute(r.db.Session, r.db.Name, collection, func(c *mgo.Collection) error {
		// 多取一个，用于判断是否有更多数据
		return c.Find(filters).Limit(size).Sort(sort).All(resultPtr)
	})

	if err != nil {
		return
	}
	
	var minCursor interface{} = nil
	var maxCursor interface{} = nil

	count := breflect.SlicePtrLen(resultPtr)
	if count > 0 {
		if reverse {
			breflect.SlicePtrReverse(resultPtr)
		}
		
		minCursorModel := breflect.SlicePtrIndexOf(resultPtr, 0)
		minCursor, err = breflect.GetStructField(minCursorModel, query.CursorSort.Property)
		if err != nil {
			return
		}

		maxCursorModel := breflect.SlicePtrIndexOf(resultPtr, count - 1)
		maxCursor, err = breflect.GetStructField(maxCursorModel, query.CursorSort.Property)
		if err != nil {
			return
		}
	}

	extra = &model.CursorExtra{
		Direction: query.Direction,
		Size:      size,
		HasMore:   count == size,
		MinCursor: minCursor,
		MaxCursor: maxCursor,
	}

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
