package reflect

import (
	"fmt"
	"bytes"
	"encoding/json"
	"reflect"
	"errors"
)

// 将string转化为结构
func CastStr2Struct(str string, result interface{}) (err error) {
	decoder := json.NewDecoder(bytes.NewReader([]byte(str)))
	decoder.UseNumber()
	err = decoder.Decode(result)
	return
}

// 将map转化为结构
func CastStruct(bean interface{}, result interface{}) (err error) {
	b, err := json.Marshal(bean)
	if err != nil {
		return
	}
	decoder := json.NewDecoder(bytes.NewReader(b))
	decoder.UseNumber()
	err = decoder.Decode(result)
	return
}

// 根据bean的实际类型构造一个slice，并返回slice的指针
func MakeSlicePtr(bean interface{}, len int, capacity int) interface{} {
	t := reflect.TypeOf(bean)
	slice := reflect.MakeSlice(reflect.SliceOf(t), len, capacity)
	p := reflect.New(slice.Type())
	p.Elem().Set(slice)
	return p.Interface()
}

// 根据slicePtr构造一个相同类型的slice，并返回slice的指针
func DuplicateSlicePtr(slicePtr interface{}, len int, capacity int) interface{} {
	t := reflect.Indirect(reflect.ValueOf(slicePtr)).Type().Elem()
	slice := reflect.MakeSlice(reflect.SliceOf(t), len, capacity)
	p := reflect.New(slice.Type())
	p.Elem().Set(slice)
	return p.Interface()
}

// 获取slicePtr指向的slice的长度
func SlicePtrLen(slicePtr interface{}) int {
	return reflect.Indirect(reflect.ValueOf(slicePtr)).Len()
}

// 获取slicePtr指向的slice的元素
func SlicePtrIndexOf(slicePtr interface{}, index int) interface{} {
	return reflect.Indirect(reflect.ValueOf(slicePtr)).Index(index).Interface()
}

// 对slicePtr指向的slice进行re-slice操作，并使toPtr指向这个新的slice
func SlicePtrSlice3To(slicePtr interface{}, i int, j int, k int, toPtr interface{}) {
	v := reflect.Indirect(reflect.ValueOf(slicePtr)).Slice3(i, j, k)
	to := reflect.ValueOf(toPtr)
	to.Elem().Set(v)
}

// 将toPtr指向slicePtr的内容
func SlicePtrCloneTo(slicePtr interface{}, toPtr interface{}) {
	to := reflect.ValueOf(toPtr)
	v := reflect.Indirect(reflect.ValueOf(slicePtr))
	to.Elem().Set(v)
}

// 将slicePtr指向的slice反向，
func SlicePtrReverse(slicePtr interface{}) {
	v := reflect.Indirect(reflect.ValueOf(slicePtr))
	n := v.Len()
	swap := reflect.Swapper(v.Interface())
	for i, j := 0, n-1; i < j; i, j = i+1, j-1 {
		swap(i, j)
	}
}

// 根据ptr指向的类型构造一个新指针
func NewPtr(ptr interface{}) interface{} {
	return reflect.New(reflect.TypeOf(ptr).Elem()).Interface()
}

// 将interface转化为slice形式
func ToSlice(bean interface{}) []interface{} {
	v := reflect.ValueOf(bean)
	size := v.Len()
	results := make([]interface{}, size)
	for i := 0; i < size; i++ {
		results[i] = v.Index(i).Interface()
	}
	return results
}

// 解引用指针
func DereferencePtr(ptr interface{}) interface{} {
	return reflect.ValueOf(ptr).Elem().Interface()
}

// 解引用指向slice的指针
func DereferencePtrToSlice(ptr interface{}) []interface{} {
	return ToSlice(reflect.Indirect(reflect.ValueOf(ptr)).Interface())
}

func indirectValue(reflectValue reflect.Value) reflect.Value {
	for reflectValue.Kind() == reflect.Ptr {
		reflectValue = reflectValue.Elem()
	}
	return reflectValue
}

func GetStructField(v interface{}, fieldName string) (reflect.Value, error) {
	field := indirectValue(reflect.ValueOf(v)).FieldByName(fieldName)

	if !field.IsValid() || !field.CanSet() {
		return reflect.Value{}, errors.New(fmt.Sprintf("ErrInvalid %s", v))
	}

	return field, nil
}
