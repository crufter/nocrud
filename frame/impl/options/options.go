package options

import (
	"github.com/opesun/nocrud/frame/impl/nesteddata"
	iface "github.com/opesun/nocrud/frame/interfaces"
)

type Options struct {
	opts iface.NestedData
	m    iface.NestedData
}

func New(opts, m map[string]interface{}) *Options {
	return &Options{
		nesteddata.New(opts),
		nesteddata.New(m),
	}
}

func (o *Options) Document() iface.NestedData {
	return o.opts
}

func (o *Options) Modifiers() iface.NestedData {
	return o.m
}
