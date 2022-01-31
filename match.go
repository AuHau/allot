package allot

import (
	"errors"
	"fmt"
	"strconv"
)

// MatchInterface describes how to access a Match
type MatchInterface interface {
	String(name string) (string, error)
	Integer(name string) (int, error)
	Match(position int) (string, error)

	Parameter(paramName string) (string, error)
}

// Match is the Match definition
type Match struct {
	Command CommandInterface
	Request string
}

// String returns the value as string
func (m Match) String(name string) (string, error) {
	return m.Parameter(name)
}

// Integer returns the value as integer
func (m Match) Integer(name string) (int, error) {
	str, err := m.Parameter(name)
	if err != nil {
		return 0, err
	}

	return strconv.Atoi(str)
}

// Parameter returns the value for a parameter
func (m Match) Parameter(name string) (string, error) {
	pos := m.Command.Position(name)

	if pos == -1 {
		return "", fmt.Errorf("Unknown parameter '%s'", name)
	}

	expr := m.Command.Expression()
	matches := expr.FindAllStringSubmatch(m.Request, -1)[0][1:]

	return matches[pos], nil
}

// Match returns the match at given position
func (m Match) Match(position int) (string, error) {
	expr := m.Command.Expression()
	matches := expr.FindAllStringSubmatch(m.Request, -1)

	if len(matches) != 1 {
		return "", errors.New("Unable to parse request")
	}

	if position >= len(matches[0]) {
		return "", fmt.Errorf("No parameter at position %d", position)
	}

	return matches[0][position+1], nil
}
