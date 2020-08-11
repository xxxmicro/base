package gorm

import(
	"strings"
	_gorm "github.com/jinzhu/gorm"
)

var (
	fieldsCache = make(map[string]map[string]*_gorm.StructField)
)

// 约定 name为小写
func FindColumn(name string, ms *_gorm.ModelStruct, dbHandler *_gorm.DB) (*_gorm.StructField, bool) {
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