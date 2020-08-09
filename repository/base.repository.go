package repository

import(
	"context"
	_gorm "github.com/jinzhu/gorm"
	"github.com/xxxmicro/base/database/gorm"
	"github.com/xxxmicro/base/database/models"
)

type BaseRepository interface {
	Create(c context.Context, model models.Model) error
	Update(c context.Context, id interface{}, model models.Model) error
	FindOne(c context.Context, id interface{}, model models.Model) error	
}

type BaseRepositoryImpl struct {
	db *_gorm.DB
}

func (r *BaseRepositoryImpl) Create(c context.Context, model models.Model) error {
	db := gorm.SetSpanToGorm(c, r.db)

	model.GenerateID()
	return db.Table(model.GetTable()).Create(model).Error
}

func (r *BaseRepositoryImpl) Update(context context.Context, id interface{}, model models.Model) error {
	db := gorm.SetSpanToGorm(context, r.db)

	return db.Table(model.GetTable()).Where("id=?", id).Save(model).Error
}

func (r *BaseRepositoryImpl) FindOne(context context.Context, id interface{}, model models.Model) error {
	db := gorm.SetSpanToGorm(context, r.db)

	return db.Table(model.GetTable()).Where("id=?", id).Take(model).Error
}


