package mongo

import(
	"gopkg.in/mgo.v2"
	"github.com/xxxmicro/base/domain/model"
)

type Indexed interface {
	Indexes() []mgo.Index
	model.Model
}
