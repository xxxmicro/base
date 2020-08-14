package mongo

import(
	"time"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"github.com/xxxmicro/base/domain/model"
)

type Indexed interface {
	Indexes() []mgo.Index
	model.Model
}
