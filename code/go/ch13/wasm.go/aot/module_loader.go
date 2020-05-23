package aot

import (
	"errors"
	"plugin"

	"wasm.go/instance"
)

type NewFn = func(instance.Map) (instance.Module, error)

// load compiled module
func Load(filename string, mm instance.Map) (instance.Module, error) {
	p, err := plugin.Open(filename)
	if err != nil {
		return nil, err
	}

	f, err := p.Lookup("Instantiate")
	if err != nil {
		return nil, err
	}

	newFn, ok := f.(NewFn)
	if !ok {
		msg := "'Instantiate' is not 'func(instance.Map) (instance.Module, error)'"
		return nil, errors.New(msg)
	}

	return newFn(mm) // TODO
}
