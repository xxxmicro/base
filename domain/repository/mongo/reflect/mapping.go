package reflect

import (
	"fmt"
	"reflect"
	"strings"
)

var StructInfoMap = make(map[reflect.Type]*StructInfo)

//结构体信息
type StructInfo struct {
	FieldsMap map[string]*StructField //字段字典集合
	Name      string                  //类型名
}

//结构体字段信息
type StructField struct {
	Name           string //字段名
	FieldType      reflect.Type
	TableFieldName string //表属性名
	Primary        bool   //是否主键字段
}

//获得结构体的信息
func GetStructInfo(target interface{}, customize func(*StructField)) (*StructInfo, error) {
	v := reflect.ValueOf(target)
	if v.Kind() != reflect.Ptr {
		return nil, fmt.Errorf("not ptr param")
	}
	t := v.Elem().Type()
	//判断target的类型
	if t.Kind() != reflect.Struct {
		return nil, fmt.Errorf("not struct param")
	}
	return getReflectInfo(t, customize)
}

// 获得结构体的反射的信息
func getReflectInfo(t reflect.Type, customize func(*StructField)) (*StructInfo, error) {
	var structInfo *StructInfo

	fieldsMap := make(map[string]*StructField)
	// 从map里取结构体信息, 如果map没有则新建一个然后存map
	if value, ok := StructInfoMap[t]; ok {
		structInfo = value
	} else {
		// 遍历所有属性
		for index := 0; index < t.NumField(); index++ {
			structField := t.Field(index)
			// 数据库字段名
			tableField := strings.TrimSpace(structField.Tag.Get("bson"))
			structFieldType := structField.Type

			if len(tableField) != 0 {
				if structField.Type.Kind() == reflect.Ptr {
					structFields := parseEmbedStruct(structField)
					for _, v := range structFields {
						v.Name = structField.Name + "." + v.Name
						v.TableFieldName = tableField + "." + v.TableFieldName
						fieldsMap[v.TableFieldName] = v
					}
					continue
				}

				// 构造一个新的StructField
				sf := &StructField{
					Name:           structField.Name,
					TableFieldName: tableField,
					FieldType:      structFieldType,
				}
				// 将新的StructField放入Map
				fieldsMap[tableField] = sf

				if customize != nil {
					customize(sf)
				}
			}
		}

		//构造一个新的StructInfo
		structInfo = &StructInfo{
			Name:      t.Name(),
			FieldsMap: fieldsMap,
		}
		//将新的StructInfo放入Map当缓存用
		StructInfoMap[t] = structInfo
	}
	return structInfo, nil
}

func parseEmbedStruct(embedField reflect.StructField) []*StructField {
	t := embedField.Type.Elem()
	sfSlice := make([]*StructField, 0)
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		tableField := strings.TrimSpace(field.Tag.Get("bson"))
		structFieldType := field.Type

		if len(tableField) != 0 {
			if field.Type.Kind() == reflect.Ptr {
				structFields := parseEmbedStruct(field)

				for _, v := range structFields {
					v.Name = field.Name + "." + v.Name
					v.TableFieldName = tableField + "." + v.TableFieldName
					sfSlice = append(sfSlice, v)
				}
				continue
			}

			sf := &StructField{
				Name:           field.Name,
				TableFieldName: tableField,
				FieldType:      structFieldType,
			}

			sfSlice = append(sfSlice, sf)

		}

	}
	return sfSlice
}
