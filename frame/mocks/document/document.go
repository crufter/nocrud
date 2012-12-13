// This mock is mainly designed to be used with a "user" implementation, where a user is not logged in,
// but a nil Document as an embedded field will cause runtime panics.
package document

import(
	iface "github.com/opesun/nocrud/frame/interfaces"
	"fmt"
)

type Document struct {
	data	map[string]interface{}
}

func New(a map[string]interface{}) iface.Document {
	return &Document{
		a,
	}
}

func NewEmpty() iface.Document {
	return &Document{
		map[string]interface{}{},
	}
}

func (d *Document) Id() iface.Id {
	return nil
}

func (d *Document) Data() map[string]interface{} {
	return d.data
}

var noDb = fmt.Errorf("This is a mock document, no connection to database.")

func (d *Document) Update(q map[string]interface{}) error {
	return noDb
}

func (d *Document) Remove() error {
	return noDb
}

