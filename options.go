package cronos

import "reflect"

//OptionsList - list of options
type OptionsList []Options

// Options - dependency injection options
type Options struct {
	key      key
	task     func(objectKey key, cronos *Cronos)
	priority int
}

func (ol OptionsList) Len() int {
	return len(ol)
}

func (ol OptionsList) Less(i, j int) bool {
	return ol[i].priority < ol[j].priority
}

func (ol OptionsList) Swap(i, j int) {
	ol[i], ol[j] = ol[j], ol[i]
}

// Singleton - determines whether the dependency is a singleton, the default is true.
func Singleton(isSingleton bool) Options {
	return Options{
		task: func(objectKey key, cronos *Cronos) {
			if !isSingleton {
				cronos.cache.notSingletons[objectKey] = !isSingleton
			}
		},
		priority: 1,
	}
}

// As - specify whether constructor returns an interface
func As(typei interface{}) Options {
	return Options{
		task: func(objectKey key, cronos *Cronos) {
			object, ok := cronos.cache.components[objectKey]

			if !ok {
				object = cronos.Fetch(objectKey.typed)
			}

			ikey := key{reflect.TypeOf(typei).Elem()}

			if _, ok := cronos.cache.components[ikey]; !ok {
				cronos.cache.components[ikey] = object
			}
		},
		priority: 2,
	}
}
