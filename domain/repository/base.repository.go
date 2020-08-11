package repository

import(
	"context"
	"github.com/xxxmicro/base/domain/model"
)

type BaseRepository interface {
	Create(c context.Context, m model.Model) error

	Update(c context.Context, m model.Model) error

	FindOne(c context.Context, id interface{}, m model.Model) error
	
	// 翻页查询
	// query: 查询条件
	// bean: 数据指针，仅用于帮助推导数据类型
	Page(c context.Context, query *model.PageQuery, m model.Model) (page *model.Page, err error)

	// 根据主键删除数据
	// id: 主键值
	// bean: 数据指针，仅用于帮助推导数据类型
	// has: 数据是否存在
	// err: 异常信息
	Delete(c context.Context, m model.Model) error

	// 游标查询
	// query: 查询条件
	// bean: 数据指针，仅用于帮助推导数据类型
	Cursor(c context.Context, query *model.CursorQuery, m model.Model) (cursor *model.Cursor, err error)
}
