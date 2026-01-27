package tools

import (
	"fmt"
	"my_agent/tools/jsonschema"
	"reflect"
)

// Tool represents a registerable function.
type Tool struct {
	Name        string
	Description string

	// the reflection value
	// this returns the value the two important things we need
	Func reflect.Value
	// this returns the datatype
	// these reflect values tell the underlying value and type of function lets say WeatherTool()
	ArgsType reflect.Type

	// maps is a kv data store , something like dict , here the string is the key type and value can be any
	Schema map[string]any
}

type Registry struct {
	tools map[string]Tool
}

func NewRegistry() *Registry {
	return &Registry{
		tools: make(map[string]Tool),
	}
}

// check if function -- get back its args -- generate json -- save it
func (r *Registry) Register(name string, description string, function any) error {

	fnType := reflect.TypeOf(function)

	if fnType.Kind() != reflect.Func {
		return fmt.Errorf("this is not a valid function please try again")
	}

	if fnType.NumIn() != 1 {
		return fmt.Errorf("function must have exactly 1 argument")
	}

	argType := fnType.In(0)

	// Generate schema using our helper
	schema := jsonschema.GenerateSchema(argType)

	// Store the tool
	r.tools[name] = Tool{
		Name:        name,
		Description: description,
		Func:        reflect.ValueOf(function),
		ArgsType:    argType,
		Schema:      schema,
	}

	return nil
}
