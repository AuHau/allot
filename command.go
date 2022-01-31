package allot

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

// CommandInterface describes how to access a Command
type CommandInterface interface {
	Expression() *regexp.Regexp
	Has(name string) bool
	Match(req string) (MatchInterface, error)
	Matches(req string) bool
	Parameters() []Parameter
	Position(name string) int
	Text() string
}

type Types = map[string]string

// Command is a Command definition
type Command struct {
	text   string
	escape bool
	params []Parameter
	expression *regexp.Regexp
}

// Text returns the command text
func (c Command) Text() string {
	return c.text
}

// Expression returns the regular expression matching the command text
func (c Command) Expression() *regexp.Regexp {
	return c.expression
}

// Parameters returns the list of defined parameters
func (c Command) Parameters() []Parameter {
	return c.params
}

// Has checks if the parameter is found in the command
func (c Command) Has(name string) bool {
	pos := c.Position(name)

	return pos != -1
}

// Position returns the position of a parameter
func (c Command) Position(name string) int {
	params := c.Parameters()

	for index, item := range params {
		if item.name == name {
			return index
		}
	}

	return -1
}

// Match returns the parameter matching the expression at the defined position
func (c Command) Match(req string) (MatchInterface, error) {
	result := c.Matches(req)

	if result {
		return Match{c, req}, nil
	}

	return nil, errors.New("Request does not match Command.")
}

// Matches checks if a command definition matches a request
func (c Command) Matches(req string) bool {
	expr:= c.Expression()

	return expr.MatchString(req)
}

func makeExpression(commandString string, params []Parameter, escape bool) (*regexp.Regexp, error) {
	var expr string

	if escape {
		expr = regexp.QuoteMeta(commandString)
	} else {
		expr = commandString
	}

	for _, param := range params {
		expr = strings.Replace(
			expr,
			fmt.Sprintf("<%s:%s>", param.Name(), param.Data()),
			fmt.Sprintf("(%s)", param.Expression().String()),
			-1,
		)

		expr = strings.Replace(
			expr,
			fmt.Sprintf("<%s>", param.Name()),
			fmt.Sprintf("(%s)", param.Expression().String()),
			-1,
		)
	}

	regex, err := regexp.Compile(fmt.Sprintf("^%s$", expr))
	if err != nil {
		return nil, err
	}

	return regex, nil
}

func parseParameters(command string, types Types) ([]Parameter, error) {
	var list []Parameter
	namesSet := make(map[string]struct{})

	re, err := regexp.Compile("<(.*?)>")
	if err != nil {
		return nil, err
	}

	result := re.FindAllStringSubmatch(command, -1)

	for _, p := range result {
		if len(p) != 2 {
			continue
		}

		param, err := Parse(p[1], types)
		if err != nil {
			return nil, err
		}

		if _, exist := namesSet[param.name]; exist {
			return nil, fmt.Errorf("the parameter with name '%s' is present multiple times", param.name)
		}

		namesSet[param.name] = struct{}{}
		list = append(list, param)
	}

	return list, nil
}

func makeCommand(commandString string, escaping bool, types Types) (Command, error) {
	if types == nil {
		types = MakeBasicTypes()
	}

	var command Command
	params, err := parseParameters(commandString, types)
	if err != nil {
		return command, err
	}

	expr, err := makeExpression(commandString, params, escaping)
	if err != nil {
		return command, err
	}

	return Command{commandString, escaping, params,expr }, nil
}

// New returns a new command
func New(command string, types Types) (Command, error) {
	return makeCommand(command, false, types)
}

// NewWithEscaping returns a new command that escapes regex characters
func NewWithEscaping(command string, types Types) (Command, error) {
	return makeCommand(command, true, types)
}
