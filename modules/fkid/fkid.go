package fkid

import(
	"github.com/opesun/nocrud/frame/context"
	"github.com/opesun/nocrud/frame/misc/scut"
	iface "github.com/opesun/nocrud/frame/interfaces"
	"github.com/opesun/sanitize"
	"strings"
)

type C struct {
}

// func (c *C) Init(uni *context.Uni) {
// 	c.uni = uni
// }

func (c *C) SanitizerMangler(san *sanitize.Extractor) {
	san.AddFuncs(sanitize.FuncMap{
		"fkid": func(dat interface{}, s sanitize.Scheme) (interface{}, error) {
			cs, _ := s.Specific["comma_separated"].(bool)
			if cs {
				ret := []interface{}
				split := strings.Split(dat.(string), ",")
				for _, v := range split {
					idstr := strings.Trim(v, " ")
					id, err := scut.DecodeId(idstr)
					if err != nil {
						return nil, err
					}
					ret = append(ret, id)
				}
				return ret, nil
			}
			str := strings.Trim(dat.(string), " ")
			return scut.DecodeId(str)
		},
	})
}