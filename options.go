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

// Qualifier - determines the specific type of dependency injected as an interface.
func Qualifier(typeObject, typeInterface interface{}) Options {
	return Options{
		task: func(objectKey key, cronos *Cronos) {

			qualifierObjectKey := key{reflect.TypeOf(typeObject).Elem()}
			qualificerInterfaceKey := key{reflect.TypeOf(typeInterface)}

			qualifierObject := cronos.Fetch(qualifierObjectKey.typed)
			qualifierObjectValue := reflect.ValueOf(qualifierObject)

			constructor := cronos.cache.constructors[objectKey]
			argsOldConstructor := cronos.getArgs(constructor)

			args := make([]reflect.Value, len(argsOldConstructor))

			for i := 0; i < len(argsOldConstructor); i++ {
				objectInterfaceType := qualificerInterfaceKey.typed.Elem()

				if argsOldConstructor[i].Implements(objectInterfaceType) {
					args[i] = qualifierObjectValue
				} else {
					args[i] = reflect.ValueOf(cronos.Fetch(argsOldConstructor[i]))
				}
			}

			newConstructor := reflect.MakeFunc(reflect.TypeOf(constructor), func(arguments []reflect.Value) []reflect.Value {
				return reflect.ValueOf(constructor).Call(args)
			})

			cronos.cache.constructors[objectKey] = newConstructor.Interface()
		},
		priority: 3,
	}
}
