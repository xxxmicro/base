package gorm

import(
	"errors"
	_gorm "github.com/jinzhu/gorm"
)

// ErrFilter
var (
	ErrFilter          		= errors.New("过滤参数错误")
	ErrFilterValueType 		= errors.New("过滤值类型错误")
	ErrFilterValueSize 		= errors.New("过滤值大小错误")
	ErrFilterOperate   		= errors.New("过滤操作错误")
)

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
