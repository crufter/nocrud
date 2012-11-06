package viewcontext

import(
	iface "github.com/opesun/nocrud/frame/interfaces"
)

type ViewContext struct {
	context		map[string]interface{}
}

func New() *ViewContext {
	return &ViewContext{
		map[string]interface{}{},
	}
}

func (v *ViewContext) Publish(key string, data interface{}) iface.ViewContext {
	val, ok := v.context[key]
	if ok {
		m, ism := val.(map[string]interface{})
		m1, ism1 := data.(map[string]interface{})
		if ism && ism1 {
			for i, v := range m1 {
				m[i] = v
			}
			v.context[key] = m
			return v
		} else {
			panic("You can't overwrite data.")
		}
	}
	v.context[key] = data
	return v
}

func (v *ViewContext) Get() map[string]interface{} {
	return v.context
}