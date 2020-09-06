package mapper

import (
	"fmt"
	"sync"
	"reflect"
	"github.com/petersunbag/coven"
)

var (
	mutex sync.Mutex
	c_map  = make(map[string]*coven.Converter)
	options = make(map[string]*coven.StructOption)
)

func RegisterOption(src, dst interface{}, option *coven.StructOption) {
	key := fmt.Sprintf("%v_%v", reflect.TypeOf(src).String(), reflect.TypeOf(dst).String())

	options[key] = option
}

func Map(src, dst interface{}) (err error) {
	key := fmt.Sprintf("%v_%v", reflect.TypeOf(src).String(), reflect.TypeOf(dst).String())
	if _, ok := c_map[key]; !ok {
		mutex.Lock()
		defer mutex.Unlock()

		if option, ok := options[key]; ok {
			if c_map[key], err = coven.NewConverterOption(dst, src, option); err != nil {
				return
			}
		} else {
			if c_map[key], err = coven.NewConverter(dst, src); err != nil {
				return
			}
		}
	}

	if err = c_map[key].Convert(dst, src); err != nil {
		return
	}
	return
}