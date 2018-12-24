/*
Package factory is a fixtures replacement for Go inspired by factory_bot in rails and factory_girl in javascript.
It is built to be simple to set up and use.
*/
package factory

import (
	"fmt"
	"reflect"
)

// Builder is the interface of the factory builder.
type Builder interface {
	// Defines a factory which later can be used to create instances of a struct
	Define(model interface{}, definitions Definition) error
	// Build creates an instance of a previously defined struct, model should be a pointer to the struct that will be populated
	Build(model interface{}, overrides Definition) error
}

type factory struct {
	model       reflect.Type
	definitions Definition
}

type defaultFactoryBuilder struct {
	factories map[string]factory
}

func (f *defaultFactoryBuilder) Define(model interface{}, definitions Definition) error {
	val := typeOf(model)
	fac := factory{model: val, definitions: definitions}
	name := fac.model.Name()
	f.factories[name] = fac
	return nil
}

func typeOf(m interface{}) reflect.Type {
	val := reflect.TypeOf(m)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	return val
}

func nameOf(model interface{}) string {
	val := reflect.TypeOf(model)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	return val.Name()
}

func (f *defaultFactoryBuilder) Build(model interface{}, overrides Definition) error {
	modelType := reflect.TypeOf(model)
	name := nameOf(model)

	fac, ok := f.factories[name]
	if !ok {
		return fmt.Errorf("unregistered model: %s", model)
	}
	if modelType.Kind() != reflect.Ptr {
		return fmt.Errorf("model needs to be a pointer")
	}

	obj := reflect.ValueOf(model).Elem()
	o := reflect.TypeOf(model).Elem()

	names, indexes := getFields([]int{}, o)

	for i, n := range names {
		v, ok := overrides[n]
		if !ok {
			v, ok = fac.definitions[n]
		}
		if ok {
			field := obj.FieldByIndex(indexes[i])
			val := reflect.ValueOf(v)
			if val.Kind() == reflect.Func {
				res := val.Call([]reflect.Value{})
				val = res[0]
			}
			field.Set(val)
		}
	}
	return nil
}

func getFields(parent []int, obj reflect.Type) ([]string, [][]int) {
	fc := obj.NumField()
	names := make([]string, 0, 10)
	indexes := make([][]int, 0, 10)
	for i := 0; i < fc; i += 1 {
		field := obj.Field(i)
		if field.Anonymous {
			an, ai := getFields(append(parent, field.Index...), field.Type)
			names = append(names, an...)
			indexes = append(indexes, ai...)
		} else {
			names = append(names, field.Name)
			idx := append(parent, field.Index...)
			indexes = append(indexes, idx)
		}
	}
	return names, indexes
}

/*
Definition is the template on how to populate the struct
It is a map where the key is the name of the field to populate. The value can either be the value to populate the field,
or a function that returns the value.
*/
type Definition map[string]interface{}

var builder Builder

func init() {
	builder = New()
}

// New creates a new builder with a separate definition registry
func New() Builder {
	return &defaultFactoryBuilder{
		factories: make(map[string]factory),
	}
}

// Define defines the model in the global registry
func Define(model interface{}, definitions Definition) error {
	return builder.Define(model, definitions)
}

// Build creates an instance of a model previously registered in the global registry
func Build(model interface{}, overrides Definition) error {
	return builder.Build(model, overrides)
}
