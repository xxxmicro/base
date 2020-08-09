package repository

import(
	"context"
	"github.com/xxxmicro/base/repository/models"
)

type BaseRepository interface {
	Create(c context.Context, model models.Model) error
	Update(c context.Context, id interface{}, model models.Model) error
	FindOne(c context.Context, id interface{}, model models.Model) error	
}
