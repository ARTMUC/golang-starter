package di

import (
	"fmt"
	"log"
	"reflect"
)

type DependencyContainer interface {
	Provide(constructor interface{})
	Get(dependencyType reflect.Type) (interface{}, error)
}

type container struct {
	dependencies map[reflect.Type]reflect.Value
}

func NewContainer() DependencyContainer {
	return &container{
		dependencies: make(map[reflect.Type]reflect.Value),
	}
}

func (c *container) Provide(constructor interface{}) {
	t := reflect.TypeOf(constructor)
	if t.Kind() != reflect.Func {
		log.Fatalf("provided constructor must be a function: %v", t)
	}

	numParams := t.NumIn()
	params := make([]reflect.Value, numParams)
	for i := 0; i < numParams; i++ {
		paramType := t.In(i)
		dependency, found := c.dependencies[paramType]
		if !found {
			panic(fmt.Sprintf("Dependency not found for type %s", paramType))
		}
		params[i] = dependency
	}

	fnValue := reflect.ValueOf(constructor)
	if t.Kind() == reflect.Ptr {
		fnValue = reflect.Indirect(fnValue)
	}
	result := fnValue.Call(params)

	singleton := result[0].Interface()

	c.dependencies[t.Out(0)] = reflect.ValueOf(singleton)
}

func (c *container) Get(dependencyType reflect.Type) (interface{}, error) {
	dependencyValue, ok := c.dependencies[dependencyType]
	if !ok {
		return nil, fmt.Errorf("dependency not found: %v", dependencyType)
	}

	return dependencyValue, nil
}

func MustGet[T any](c DependencyContainer, t *T) *T {
	val, err := c.Get(reflect.TypeOf(t))
	if err != nil {
		log.Fatalf("dependency not found: %v", t)
	}
	value := val.(reflect.Value)

	return value.Interface().(*T)
}

func GetProvider[T any](c DependencyContainer, t *T) {
	val, err := c.Get(reflect.TypeOf(t))
	if err != nil {
		log.Fatalf("dependency not found: %v", t)
	}
	value := val.(reflect.Value)
	t = value.Interface().(*T)
}
