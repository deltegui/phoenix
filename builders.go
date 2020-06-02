package phoenix

import (
	"log"
	"reflect"
)

// Builder is a function that expects nothing and retuns
// the type that builds. It's represented as interface.
type Builder interface{}

type Injector struct {
	builders map[reflect.Type]Builder
}

func NewInjector() *Injector {
	return &Injector{
		builders: make(map[reflect.Type]Builder),
	}
}

// Add a builder to the dependency injector.
func (injector Injector) Add(builder Builder) {
	outputType := reflect.TypeOf(builder).Out(0)
	injector.builders[outputType] = builder
}

// ShowAvailableBuilders prints all registered builders.
func (injector Injector) ShowAvailableBuilders() {
	for k := range injector.builders {
		log.Printf("Builder for type: %s\n", k)
	}
}

// Get returns a builded dependency
func (injector Injector) Get(name interface{}) interface{} {
	return injector.GetByType(reflect.TypeOf(name))
}

// GetByType returns a builded dependency identified by type
func (injector Injector) GetByType(name reflect.Type) interface{} {
	dependencyBuilder := injector.builders[name]
	if dependencyBuilder == nil {
		log.Panicf("Builder not found for type %s\n", name)
	}
	return injector.CallBuilder(dependencyBuilder)
}

func (injector Injector) CallBuilder(builder Builder) interface{} {
	var inputs []reflect.Value
	builderType := reflect.TypeOf(builder)
	for i := 0; i < builderType.NumIn(); i++ {
		impl := injector.builders[builderType.In(i)]
		if impl == nil {
			log.Panicf("Builder not found for type %s\n", builderType.In(i))
		}
		result := injector.CallBuilder(impl)
		inputs = append(inputs, reflect.ValueOf(result))
	}
	builderVal := reflect.ValueOf(builder)
	builded := builderVal.Call(inputs)
	return builded[0].Interface()
}

// PopulateStruct fills a struct with the implementations
// that the injector can create. Make sure you pass a reference and
// not a value
func (injector Injector) PopulateStruct(userStruct interface{}) {
	ptrStructValue := reflect.ValueOf(userStruct)
	structValue := ptrStructValue.Elem()
	if structValue.Kind() != reflect.Struct {
		log.Panicln("Value passed to PopulateStruct is not a struct")
	}
	for i := 0; i < structValue.NumField(); i++ {
		field := structValue.Field(i)
		if field.IsValid() && field.CanSet() {
			impl := injector.GetByType(field.Type())
			field.Set(reflect.ValueOf(impl))
		}
	}
}
