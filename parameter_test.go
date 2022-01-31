package allot

import (
	"regexp"
	"testing"
)

var basicTypes = MakeBasicTypes()

func TestExpression(t *testing.T) {
	var data = []struct {
		data       string
		expression string
	}{
		{"string", "[^\\s]+"},
		{"integer", "[0-9]+"},
		{"unknown", ""},
	}

	for _, set := range data {
		exp, err := regexp.Compile(set.expression)

		if err != nil {
			t.Errorf("TextExpression regexp does not compile: %s", set.expression)
		}

		test := Expression(set.data, basicTypes)

		if test == nil && set.expression != "" {
			t.Errorf("Expression for data \"%s\" is not valid", set.data)
		}

		if test != nil && test.String() != exp.String() {
			t.Errorf("Expression() not matching test data! got \"%s\", expected \"%s\"", test.String(), exp.String())
		}
	}
}

func TestParse(t *testing.T) {
	var data = []struct {
		text string
		name string
		data string
	}{
		{"<lorem>", "lorem", "string"},
		{"<ipsum:integer>", "ipsum", "integer"},
	}

	var param Parameter
	var err error
	for _, set := range data {
		param, err = Parse(set.text, basicTypes)
		if err != nil {
			t.Errorf("Parse returned err: %s", err)
		}

		if param.Name() != set.name {
			t.Errorf("param.Name() should be \"%s\", but is \"%s\"", set.name, param.Name())
		}
	}

	_, err = Parse("<some:nonexistingtype>", basicTypes)
	if err == nil {
		t.Error("Parse should return error for non-existing types")
	}
}

