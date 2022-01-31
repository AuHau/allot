package allot

import (
	"fmt"
	"regexp"
	"strings"
)

func MakeBasicTypes() Types {
	return map[string]string{
		"string":  "[^\\s]+",
		"integer": "[0-9]+",
	}
}

const defaultType = "string"

// Expression returns the regexp for a data type
func Expression(data string, types Types) *regexp.Regexp {
	if exp, ok := types[data]; ok {
		return regexp.MustCompile(exp)
	}

	return nil
}

// ParameterInterface describes how to access a Parameter
type ParameterInterface interface {
	Equals(param ParameterInterface) bool
	Expression() *regexp.Regexp
	Name() string
	Data() string
}

// Parameter is the Parameter definition
type Parameter struct {
	name string
	data string
	expr *regexp.Regexp
}

// Expression returns the regexp behind the type
func (p Parameter) Expression() *regexp.Regexp {
	return p.expr
}

// Name returns the Parameter name
func (p Parameter) Name() string {
	return p.name
}

// Data returns the Parameter name
func (p Parameter) Data() string {
	return p.data
}

// Equals checks if two parameter are equal
func (p Parameter) Equals(param ParameterInterface) bool {
	return p.Name() == param.Name() && p.Expression().String() == param.Expression().String()
}

// NewParameterWithType returns
func NewParameterWithType(name string, data string, regex* regexp.Regexp) Parameter {
	return Parameter{name, data, regex}
}

// Parse parses parameter info
func Parse(text string, types Types) (Parameter, error) {
	var splits []string
	var name, data string
	var param Parameter

	name = strings.Replace(text, "<", "", -1)
	name = strings.Replace(name, ">", "", -1)
	data = defaultType

	if strings.Contains(name, ":") {
		splits = strings.Split(name, ":")

		name = splits[0]
		data = splits[1]

		_, exists := types[data]
		if !exists {
			return param, fmt.Errorf("data types '%s' is not defined", data)
		}
	}

	regex := Expression(data, types)
	param = NewParameterWithType(name, data, regex)

	return param, nil
}
