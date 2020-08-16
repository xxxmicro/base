package mongo

import (
	"context"
	"errors"
	"fmt"
	"github.com/xxxmicro/base/database/mongo"
	"github.com/xxxmicro/base/domain/model"
	"github.com/xxxmicro/base/domain/repository"
	reflect2 "github.com/xxxmicro/base/domain/repository/mongo/reflect"
	breflect "github.com/xxxmicro/base/reflect"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"reflect"
)

type baseRepository struct {
	db *mongo.DB
}

func NewBaseRepository(db *mongo.DB) repository.BaseRepository {
	return &baseRepository{ db }
}

func (r *baseRepository) Create(c context.Context, m model.Model) error {
	ms, err := reflect2.GetStructInfo(m, nil)
	if err != nil {
		return err
	}
	collection := TheNamingStrategy.Table(ms.Name)

	// TODO 找出 ctime, utime 的 tag 进行设置

	return Execute(r.db.Session, r.db.Name, collection, func(c *mgo.Collection) error {
		return c.Insert(m)
	})
}

func (r *baseRepository) Upsert(c context.Context, m model.Model) (changeInfo *repository.ChangeInfo, err error) {
	ms, err := reflect2.GetStructInfo(m, nil)
	if err != nil {
		return
	}
	collection := TheNamingStrategy.Table(ms.Name)

	Execute(r.db.Session, r.db.Name, collection, func(c *mgo.Collection) error {
		var change *mgo.ChangeInfo
		change, err = c.Upsert(m.Unique(), m)
		if err != nil {
			return err
		}

		changeInfo = &repository.ChangeInfo{
			Updated: change.Updated,
			Removed: change.Removed,
			Matched: change.Matched,
			UpsertedId: change.UpsertedId,
		}
		return nil
	})
	return
}

func (r *baseRepository) Update(c context.Context, m model.Model, change interface{}) error {
	ms, err := reflect2.GetStructInfo(m, nil)
	if err != nil {
		return err
	}
	collection := TheNamingStrategy.Table(ms.Name)

	return Execute(r.db.Session, r.db.Name, collection, func(c *mgo.Collection) error {
		return c.Update(m.Unique(), bson.M{
			"$set": change,
		})
	})
}

func (r *baseRepository) FindOne(c context.Context, m model.Model) error {
	ms, err := reflect2.GetStructInfo(m, nil)
	if err != nil {
		return err
	}
	collection := TheNamingStrategy.Table(ms.Name)
	
	return Execute(r.db.Session, r.db.Name, collection, func(c *mgo.Collection) error {
		return c.Find(m.Unique()).One(m)
	})
}

func (r *baseRepository) Delete(c context.Context, m model.Model) error {
	ms, err := reflect2.GetStructInfo(m, nil)
	if err != nil {
		return err
	}
	collection := TheNamingStrategy.Table(ms.Name)
	
	return Execute(r.db.Session, r.db.Name, collection, func(c *mgo.Collection) error {
		return c.Remove(m.Unique())
	})
}

func (r *baseRepository) Page(c context.Context, m model.Model, query *model.PageQuery, resultPtr interface{}) (total int, pageCount int, err error){
	ms, err := reflect2.GetStructInfo(m, nil)
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
		total, err = c.Find(filters).Count()
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
	ms, err := reflect2.GetStructInfo(m, nil)
	if err != nil {
		return
	}
	collection := TheNamingStrategy.Table(ms.Name)

	filters, err := buildQuery(ms, query.Filters)
	if err != nil {
		return
	}

	cursorProp, ok := ms.FieldsMap[query.CursorSort.Property]
	if !ok {
		err = errors.New(fmt.Sprintf("cursor prop(%s) not found", query.CursorSort.Property))
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


		minCursor, err = breflect.GetStructField(minCursorModel, cursorProp.Name)
		if err != nil {
			return
		}

		maxCursorModel := breflect.SlicePtrIndexOf(resultPtr, count - 1)
		maxCursor, err = breflect.GetStructField(maxCursorModel, cursorProp.Name)
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
