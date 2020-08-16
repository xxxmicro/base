package repository

import(
	"context"
	"github.com/xxxmicro/base/domain/model"
)


type ChangeInfo struct {
	Updated    int
	Removed    int         // Number of documents removed
	Matched    int         // Number of documents matched but not necessarily changed
	UpsertedId interface{} // Upserted _id field, when not explicitly provided
}


type BaseRepository interface {
	Create(c context.Context, m model.Model) error

	Upsert(c context.Context, m model.Model) (*ChangeInfo, error)

	Update(c context.Context, m model.Model, change interface{}) error

	FindOne(c context.Context, m model.Model) error
	
	// 翻页查询
	// query: 查询条件
	// m: 数据指针，仅用于帮助推导数据类型
	Page(c context.Context, m model.Model, query *model.PageQuery, resultPtr interface{}) (total int, pageCount int, err error)

	// 根据主键删除数据
	// m	数据对象
	Delete(c context.Context, m model.Model) error

	// 游标查询
	// @c	上下文
	// @query	查询条件
	// m	数据指针，仅用于帮助推导数据类型
	// resultPtr	返回数据的指针
	Cursor(c context.Context, query *model.CursorQuery, m model.Model, resultPtr interface{}) (cursor *model.CursorExtra, err error)
}
