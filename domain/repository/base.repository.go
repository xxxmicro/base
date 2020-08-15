package repository

import(
	"context"
	"github.com/xxxmicro/base/domain/model"
)

type BaseRepository interface {
	Create(c context.Context, m model.Model) error

	Update(c context.Context, m model.Model) error

	UpdateSelective(c context.Context, m model.Model, data map[string]interface{}) error

	FindOne(c context.Context, m model.Model) error
	
	// 翻页查询
	// query: 查询条件
	// m: 数据指针，仅用于帮助推导数据类型
	Page(c context.Context, query *model.PageQuery, m model.Model, resultPtr interface{}) (total int, pageCount int, err error)

	// 根据主键删除数据
	// id: 主键值
	// m: 数据指针，仅用于帮助推导数据类型
	Delete(c context.Context, m model.Model) error

	// 游标查询
	// query: 查询条件
	// bean: 数据指针，仅用于帮助推导数据类型
	Cursor(c context.Context, query *model.CursorQuery, m model.Model, resultPtr interface{}) (cursor *model.CursorExtra, err error)
}
