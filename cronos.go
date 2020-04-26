package cronos

import "reflect"

type key struct {
	typed reflect.Type
}

type Cronos struct {
}

type constructor interface{}
type component interface{}

type cache struct {
	components    map[key]component
	constructors  map[key]constructor
	notSingletons map[key]bool
	options       OptionsList
}
