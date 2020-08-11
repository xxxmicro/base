package gorm

import(
	"errors"
	"fmt"
	_gorm "github.com/jinzhu/gorm"
	"github.com/xxxmicro/base/domain/model"
)

// ErrFilter
var (
	ErrFilter          		= errors.New("过滤参数错误")
	ErrFilterValueType 		= errors.New("过滤值类型错误")
	ErrFilterValueSize 		= errors.New("过滤值大小错误")
	ErrFilterOperate   		= errors.New("过滤操作错误")
)


func buildQuery(db *_gorm.DB, ms *_gorm.ModelStruct, filters map[string]interface{}) (*_gorm.DB, error) {
	if filters == nil || len(filters) == 0 {
		return db, nil
	}

	var err error
	for key, value := range filters {
		db, err = gormFilter(db, ms, key, value)
		if err != nil {
			return nil, err
		}
	}

	return db, nil
}

func gormFilter(db *_gorm.DB, ms *_gorm.ModelStruct, key string, value interface{}) (*_gorm.DB, error) {
	filterType := model.FilterType(key)

	switch filterType {
	case model.FilterType_AND:
		{
			/* TODO
			subFilters := v.([]interface{})
			for _, item := range subFilters {
				db = buildQuery(db, subFilter, ms)
			}*/
		}
	case model.FilterType_OR:
		{
			/* TODO
			for _, item := range subFilters {
				db := buildQuery(db, subFilter, ms)
				orCond = orCond.Or(subCond)
			}
			*/
		}
	default:
		{
			if _, ok := FindColumn(key, ms, db); !ok {
				err := errors.New(fmt.Sprintf("ERR_DB_UNKNOWN_FIELD %s", key))
				return nil, err
			}

			vMap, ok := value.(map[string]interface{})
			if !ok {
				return db.Where(fmt.Sprintf("%s = ?", key), value), nil
			}

			for vKey, vValue := range vMap {
				filterType = model.FilterType(vKey)
				switch filterType {
				case model.FilterType_EQ:
					return db.Where(fmt.Sprintf("%s = ?", key), vValue), nil
				case model.FilterType_NE:
					return db.Where(fmt.Sprintf("%s != ?", key), vValue), nil
				case model.FilterType_GT:
					return db.Where(fmt.Sprintf("%s > ?", key), vValue), nil
				case model.FilterType_GTE:
					return db.Where(fmt.Sprintf("%s >= ?", key), vValue), nil	
				case model.FilterType_LT:
					return db.Where(fmt.Sprintf("%s < ?", key), vValue), nil	
				case model.FilterType_LTE:
					return db.Where(fmt.Sprintf("%s <= ?", key), vValue), nil	
				case model.FilterType_LIKE:
					return db.Where(fmt.Sprintf("%s LIKE ?", key), vValue), nil
				case model.FilterType_MATCH:
					return db.Where(fmt.Sprintf("%s LIKE ?", key), vValue), nil
				case model.FilterType_NOT_LIKE:
					return db.Not(fmt.Sprintf("%s LIKE ?", key), vValue), nil
				case model.FilterType_IN:
					return gormFilterIn(db, key, vValue)
				case model.FilterType_NOT_IN:
					return gormFilterNotIn(db, key, vValue)
				case model.FilterType_BETWEEN:
					return gormFilterBetween(db, key, vValue)
				case model.FilterType_IS_NULL:
					return db.Where(fmt.Sprintf("%s IS NULL", key)), nil
				case model.FilterType_NOT_NULL:
					return db.Where(fmt.Sprintf("%s IS NOT NULL", key)), nil
				}
			}
		}
	}

	return db, nil
}

func gormFilterIn(db *_gorm.DB, key string, value interface{}) (*_gorm.DB, error) {
	values, ok := value.([]interface{})
	if !ok {
		return nil, ErrFilterValueType
	}

	return db.Where(fmt.Sprintf("%s IN (?)", key), values), nil
}

func gormFilterNotIn(db *_gorm.DB, key string, value interface{}) (*_gorm.DB, error) {
	values, ok := value.([]interface{})
	if !ok {
		return nil, ErrFilterValueType
	}

	return db.Where(fmt.Sprintf("%s NOT IN (?)", key), values), nil
}

func gormFilterBetween(db *_gorm.DB, key string, value interface{}) (*_gorm.DB, error) {
	values, ok := value.([]interface{})
	if !ok {
		return nil, ErrFilterValueType
	}
	if len(values) != 2 {
		return nil, ErrFilterValueSize
	}
	if values[0] != nil && values[1] != nil {
		return db.Where(fmt.Sprintf("%s between ? and ?", key), values[0], values[1]), nil
	} else if values[0] != nil && values[1] == nil {
		return db.Where(fmt.Sprintf("%s >= ?", key), values[0]), nil
	} else if values[0] == nil && values[1] != nil {
		return db.Where(fmt.Sprintf("%s <= ?", key), values[1]), nil
	} else {
		return db, nil
	}
}


func buildSort(dbHandler *_gorm.DB, ms *_gorm.ModelStruct, sorts []*model.SortSpec) (db *_gorm.DB, err error) {
	if sorts == nil || len(sorts) == 0 {
		db = dbHandler
		return
	}

	for _, sort := range sorts {
		sortKey := sort.Property
		if _, ok := FindColumn(sortKey, ms, dbHandler); !ok {
			err = errors.New(fmt.Sprintf("unknown field: %s", sortKey))
			return
		}

		sortDir := string(sort.Type)
		if sortDir == "DSC" {
			sortDir = "desc"
		} else {
			sortDir = "asc"
		}

		db = dbHandler.Order(fmt.Sprintf("%s %s", sortKey, sortDir))
	}

	return
}

func pageQuery(queryHandler *_gorm.DB, pageNo int, pageSize int, resultPtr interface{}) (count int, pageCount int, err error) {
	limit, offset := getLimitOffset(pageNo-1, pageSize)

	count = 0
	queryHandler.Count(&count)
	queryHandler.Limit(limit).Offset(offset).Find(resultPtr)
	if err = queryHandler.Error; err != nil {
		return
	}

	pageCount = count / pageSize
	if count % pageSize != 0 {
		pageCount++
	}

	return
}

func getLimitOffset(pageNo, pageSize int) (limit, offset int) {
	if pageNo < 0 {
		pageNo = 0
	}

	if pageSize < 1 {
		pageSize = 20
	}
	return pageSize, pageNo * pageSize
}