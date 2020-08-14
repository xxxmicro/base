package mongo

import(
	"fmt"
	"errors"
	"gopkg.in/mgo.v2/bson"
	breflect "github.com/xxxmicro/base/reflect"
	"github.com/xxxmicro/base/domain/model"
)

func mongoCursorFilter(ms *breflect.StructInfo, cursorQuery *model.CursorQuery) (filter bson.M, sort string, reverse bool, err error) {
	prop := cursorQuery.CursorSort.Property
	if _, ok := ms.FieldsMap[prop]; !ok {
		err = errors.New(fmt.Sprintf("ERR_DB_UNKNOWN_FIELD %s", prop))
		return
	}

	switch cursorQuery.CursorSort.Type {
	case model.SortType_DSC:
		{
			if cursorQuery.Direction == 0 {
				// 游标前
				sort = prop
				reverse = true
				if cursorQuery.Cursor != nil {
					filter = bson.M{prop: bson.M{"$gt": cursorQuery.Cursor}}
				}
			} else {
				// 游标后
				sort = fmt.Sprintf("-%s", prop)
				reverse = false
				if cursorQuery.Cursor != nil {
					filter = bson.M{prop: bson.M{"$lt": cursorQuery.Cursor}}
				}
			}
		}
	default: // SortType_ASC
		{
			if cursorQuery.Direction == 0 {
				// 游标前
				sort = fmt.Sprintf("-%s", prop)
				reverse = true
				if cursorQuery.Cursor != nil {
					filter = bson.M{prop: bson.M{"$lt": cursorQuery.Cursor}}
				}
			} else {
				// 游标后
				sort = prop
				reverse = false
				if cursorQuery.Cursor != nil {
					filter = bson.M{prop: bson.M{"$gt": cursorQuery.Cursor}}
				}
			}
		}
	}

	return
}
