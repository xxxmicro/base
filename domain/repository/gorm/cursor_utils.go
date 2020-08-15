package gorm

import(
	"fmt"
	"time"
	"errors"
	_gorm "github.com/jinzhu/gorm"
	"github.com/xxxmicro/base/domain/model"
	"github.com/xxxmicro/base/types/smarttime"
)

func gormCursorFilter(queryHandler *_gorm.DB, ms *_gorm.ModelStruct, query *model.CursorQuery) (*_gorm.DB, bool, error) {
	var orderBy string
	var reverse bool

	sortKey := query.CursorSort.Property
	
	field, ok := FindField(sortKey, ms, queryHandler);
	if !ok {
		err := errors.New(fmt.Sprintf("ERR_DB_UNKNOWN_FIELD %s", sortKey))
		return nil, reverse, err
	}

	value := query.Cursor

	switch field.Struct.Type.String() {
	case "time.Time", "*time.Time":
		v, err := smarttime.Parse(value)
		if err == nil {
			value = time.Time(v)
		}
	}

	switch query.CursorSort.Type {
	case model.SortType_DSC:
		{
			if query.Direction == 0 {
				// 游标前
				orderBy = fmt.Sprintf("%s %s", sortKey, "ASC")
				reverse = true
				if query.Cursor != nil {
					queryHandler = queryHandler.Where(fmt.Sprintf("%s > ?", sortKey), value)
				}
			} else {
				// 游标后
				orderBy = fmt.Sprintf("%s %s", sortKey, "DESC")
				reverse = false
				if query.Cursor != nil {
					queryHandler = queryHandler.Where(fmt.Sprintf("%s > ?", sortKey), value)
				}
			}
		}
	default: // SortType_ASC
		{
			if query.Direction == 0 {
				// 游标前
				orderBy = fmt.Sprintf("%s %s", sortKey, "DESC")
				reverse = true
				if query.Cursor != nil {
					queryHandler = queryHandler.Where(fmt.Sprintf("%s > ?", sortKey), value)
				}
			} else {
				// 游标后
				orderBy = fmt.Sprintf("%s %s", sortKey, "ASC")
				reverse = false
				if query.Cursor != nil {
					queryHandler = queryHandler.Where(fmt.Sprintf("%s > ?", sortKey), value)
				}
			}
		}
	}

	queryHandler = queryHandler.Order(orderBy)

	return queryHandler, reverse, nil
}
