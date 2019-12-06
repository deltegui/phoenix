package injector

import (
	"log"
	"reflect"
)

var builders = make(map[reflect.Type]Builder)

// Builder is a function that expects nothing and retuns
// the type that builds. It's represented as interface
// because golang's type system is shit.
type Builder interface{}

// Add a builder to the dependency injector.
func Add(builder Builder) {
	outputType := reflect.TypeOf(builder).Out(0)
	builders[outputType] = builder
}

// ShowAvailableBuilders prints all registered builders.
func ShowAvailableBuilders() {
	for k := range builders {
		log.Printf("Builder for type: %s\n", k)
	}
}

// Get returns a builded dependency
func Get(name interface{}) interface{} {
	return GetByType(reflect.TypeOf(name))
}

// GetByType returns a builded dependency identified by type
func GetByType(name reflect.Type) interface{} {
	dependencyBuilder := builders[name]
	if dependencyBuilder == nil {
		log.Panicf("Builder not found for type %s\n", name)
	}
	return CallBuilder(dependencyBuilder)
}

func CallBuilder(builder Builder) interface{} {
	var inputs []reflect.Value
	builderType := reflect.TypeOf(builder)
	for i := 0; i < builderType.NumIn(); i++ {
		impl := builders[builderType.In(i)]
		if impl == nil {
			log.Panicf("Builder not found for type %s\n", builderType.In(i))
		}
		result := CallBuilder(impl)
		inputs = append(inputs, reflect.ValueOf(result))
	}
	builderVal := reflect.ValueOf(builder)
	builded := builderVal.Call(inputs)
	return builded[0].Interface()
}

// PopulateStruct fills a struct with the implementations
// that the injector can create. Make sure you pass a reference and
// not a value
func PopulateStruct(userStruct interface{}) {
	ptrStructValue := reflect.ValueOf(userStruct)
	structValue := ptrStructValue.Elem()
	if structValue.Kind() != reflect.Struct {
		log.Panicln("Value passed to PopulateStruct is not a struct")
	}
	for i := 0; i < structValue.NumField(); i++ {
		field := structValue.Field(i)
		if field.IsValid() && field.CanSet() {
			impl := GetByType(field.Type())
			field.Set(reflect.ValueOf(impl))
		}
	}
}
