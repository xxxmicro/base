package gorm

import(
	"fmt"
	"errors"
	_gorm "github.com/jinzhu/gorm"
	"github.com/xxxmicro/base/domain/model"
)

func GormCursorFilter(queryHandler *_gorm.DB, ms *_gorm.ModelStruct, query *model.CursorQuery) (*_gorm.DB, bool, error) {
	var orderBy string
	var reverse bool

	sortKey := query.CursorSort.Property
	if _, ok := FindColumn(sortKey, ms, queryHandler); !ok {
		err := errors.New(fmt.Sprintf("ERR_DB_UNKNOWN_FIELD %s", sortKey))
		return nil, reverse, err
	}

	switch query.CursorSort.Type {
	case model.SortType_DSC:
		{
			if query.Direction == 0 {
				// 游标前
				orderBy = fmt.Sprintf("%s %s", sortKey, "ASC")
				reverse = true
				if query.Cursor != nil {
					queryHandler = queryHandler.Where(fmt.Sprintf("%s > ?", sortKey), query.Cursor)
				}
			} else {
				// 游标后
				orderBy = fmt.Sprintf("%s %s", sortKey, "DESC")
				reverse = false
				if query.Cursor != nil {
					queryHandler = queryHandler.Where(fmt.Sprintf("%s > ?", sortKey), query.Cursor)
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
					queryHandler = queryHandler.Where(fmt.Sprintf("%s > ?", sortKey), query.Cursor)
				}
			} else {
				// 游标后
				orderBy = fmt.Sprintf("%s %s", sortKey, "ASC")
				reverse = false
				if query.Cursor != nil {
					queryHandler = queryHandler.Where(fmt.Sprintf("%s > ?", sortKey), query.Cursor)
				}
			}
		}
	}

	queryHandler = queryHandler.Order(orderBy)

	return queryHandler, reverse, nil
}
