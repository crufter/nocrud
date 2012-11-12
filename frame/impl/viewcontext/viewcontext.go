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

func merge(a, b map[string]interface{}) map[string]interface{} {
	for i, v := range b {
		a[i] = v
	}
	return a
}

func (v *ViewContext) Publish(key string, data interface{}) iface.ViewContext {
	val, ok := v.context[key]
	if ok {
		m, ism := val.(map[string]interface{})
		m1, ism1 := data.(map[string]interface{})
		if ism && ism1 {
			v.context[key] = merge(m, m1)
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