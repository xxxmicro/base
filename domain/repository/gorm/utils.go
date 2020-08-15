package gorm

import(
	"fmt"
	"strings"
	"errors"
	"time"
	_gorm "github.com/jinzhu/gorm"
	"github.com/xxxmicro/base/domain/model"
	"github.com/xxxmicro/base/types/smarttime"
)

var (
	fieldsCache = make(map[string]map[string]*_gorm.StructField)
)

// 约定 name为小写
func FindField(name string, ms *_gorm.ModelStruct, dbHandler *_gorm.DB) (*_gorm.StructField, bool) {
	name = strings.ToLower(name)

	tableName := ms.TableName(dbHandler)
	fieldsMap := fieldsCache[tableName]
	if fieldsMap == nil {
		fieldsMap = make(map[string]*_gorm.StructField)

		for _, field := range ms.StructFields {
			fieldName := strings.ToLower(field.Name)
			fieldsMap[fieldName] = field
		}

		fieldsCache[tableName] = fieldsMap
	}
	field, ok := fieldsMap[name]
	return field, ok
}

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
			/* TODO 暂时默认就是 AND
			subFilters := v.([]interface{})
			for _, item := range subFilters {
				db = buildQuery(db, subFilter, ms)
			}*/
		}
	case model.FilterType_OR:
		{
			/* TODO  暂时不支持 page 中支持 or
			for _, item := range subFilters {
				db := buildQuery(db, subFilter, ms)
				orCond = orCond.Or(subCond)
			}
			*/
		}
	default:
		{
			field, ok := FindField(key, ms, db)
			if !ok {
				err := errors.New(fmt.Sprintf("ERR_DB_UNKNOWN_FIELD %s", key))
				return nil, err
			}

			vMap, ok := value.(map[string]interface{})
			if !ok {
				switch field.Struct.Type.String() {
				case "time.Time", "*time.Time":
					v, err := smarttime.Parse(value)
					if err == nil {
						value = time.Time(v)
					}	
				}
			
				return db.Where(fmt.Sprintf("%s = ?", key), value), nil
			}

			for vKey, vValue := range vMap {
				switch field.Struct.Type.String() {
				case "time.Time", "*time.Time":
					v, err := smarttime.Parse(vValue)
					if err == nil {
						vValue = v
					}
				}

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
		if _, ok := FindField(sortKey, ms, dbHandler); !ok {
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
