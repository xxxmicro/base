package gorm

import(
	"context"
	_gorm "github.com/jinzhu/gorm"
	"github.com/xxxmicro/base/database/gorm"
	"github.com/xxxmicro/base/domain/repository"
	"github.com/xxxmicro/base/domain/model"
	breflect "github.com/xxxmicro/base/reflect"
)

type baseRepository struct {
	db *_gorm.DB
}

func NewBaseRepository(db *_gorm.DB) repository.BaseRepository {
	return &baseRepository{ db }
}

func (r *baseRepository) Create(c context.Context, m model.Model) error {
	db := gorm.SetSpanToGorm(c, r.db)

	return db.Create(m).Error
}

func (r *baseRepository) Update(c context.Context, m model.Model) error {
	db := gorm.SetSpanToGorm(c, r.db)

	return db.Model(m).Save(m).Error
}

func (r *baseRepository) FindOne(c context.Context, condition interface{}, m model.Model) error {
	db := gorm.SetSpanToGorm(c, r.db)

	return db.Model(m).Where(condition).Take(m).Error
}

func (r *baseRepository) Delete(c context.Context, m model.Model) error {
	// TODO 这里要做主键保护，如果 m 什么都没设置，这里将会删除表的所有记录
	return r.db.Delete(m).Error
}


func (r *baseRepository) Page(c context.Context, query *model.PageQuery, m model.Model) (page *model.Page, err error){
	items := breflect.MakeSlicePtr(m, 0, 0)

	ms := r.db.NewScope(m).GetModelStruct()

	dbHandler := r.db.Model(m)
	dbHandler, err = buildQuery(dbHandler, ms, query.Filters)
	if err != nil {
		return
	}

	dbHandler, err = buildSort(dbHandler, ms, query.Sort)
	if err != nil {
		return
	}

	total, pageCount, err := pageQuery(dbHandler, query.PageNo, query.PageSize, items)
	page = &model.Page{
		Total: total,
		PageNo: query.PageNo,
		PageSize: query.PageSize,
		PageCount: pageCount,
		Content: items,
	}
	
	return
}

func (r *baseRepository) Cursor(c context.Context, query *model.CursorQuery, m model.Model) (cursor *model.Cursor, err error) {
	return
}