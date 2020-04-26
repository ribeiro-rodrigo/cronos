package cronos

import "reflect"

type key struct {
	typed reflect.Type
}

type Cronos struct {
	cache
}

func (cronos *Cronos) Fetch(typed reflect.Type) interface{} {
	return nil
}

func (cronos *Cronos) getArgs(constructor constructor) []reflect.Type {
	return nil
}

type constructor interface{}
type component interface{}

type cache struct {
	components    map[key]component
	constructors  map[key]constructor
	notSingletons map[key]bool
	options       OptionsList
}
