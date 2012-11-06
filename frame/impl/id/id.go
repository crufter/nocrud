package id

import(
	iface "github.com/opesun/nocrud/frame/interfaces"
	"labix.org/v2/mgo/bson"
)

type Id string

func (i *Id) String() string {
}

func NewId() iface.Id {
	bson.NewObjectId().Hex()
}

func ToId(encodedForm string) iface.Id {
	
}