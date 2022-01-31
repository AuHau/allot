package allot

import (
	"testing"
)

var resultCommand bool

func checkErr(err error, t *testing.T) {
	if err != nil {
		t.Error(err)
	}
}

func BenchmarkMatches(b *testing.B) {
	var r bool

	for n := 0; n < b.N; n++ {
		cmd, err := New("command <lorem:integer> <ipsum:string>", nil)
		if err != nil {
			b.Error(err)
		}

		r = cmd.Matches("command 12345 abcdef")
	}

	resultCommand = r
}

func TestMatches(t *testing.T) {
	var data = []struct {
		command string
		request string
		matches bool
	}{
		{"command", "example", false},
		{"command", "command", true},
		{"command", "command example", false},
		{"command <lorem>", "command", false},
		{"command <lorem>", "command example", true},
		{"command <lorem>", "command 1234567", true},
		{"command <lorem>", "command command", true},
		{"command <lorem>", "example command", false},
		{"command <lorem:integer>", "command example", false},
		{"command <lorem:integer>", "command 1234567", true},
		{"command <lorem>", "command command command", false},
	}

	for _, set := range data {
		cmd, err := New(set.command, nil)
		checkErr(err, t)

		matches := cmd.Matches(set.request)

		if matches != set.matches {
			t.Errorf("Matches() returns unexpected values. Got \"%v\", expected \"%v\"\nExpression: \"%s\" not matching \"%s\"", matches, set.matches, cmd.Expression().String(), set.request)
		}
	}
}

func TestNotCompiling(t *testing.T) {
	_, err := New("command with []() invalid regex syntax", nil)

	if err == nil {
		t.Error("Compilation of regex should have had failed")
	}
}

func TestDuplicatedNames(t *testing.T) {
	_, err := New("command with <name:string> duplicated <name:string>", nil)

	if err == nil {
		t.Error("Compilation of regex should have had failed")
	}
}

func TestEscapeMatches(t *testing.T) {
	var data = []struct {
		command string
		request string
		matches bool
	}{
		{"[command]", "example", false},
		{"[command]", "[command]", true},
		{"command", "command example", false},
		{"[command] (<lorem>)", "command", false},
		{"[command] (<lorem>)", "[command] (example)", true},
		{"[command] (<lorem>)", "[command] (1234)", true},
		{"[command] (<lorem:integer>)", "[command] (1234)", true},
	}

	for _, set := range data {
		cmd, err := NewWithEscaping(set.command, nil)
		checkErr(err, t)

		matches := cmd.Matches(set.request)

		if matches != set.matches {
			t.Errorf("Matches() returns unexpected values. Got \"%v\", expected \"%v\"\nExpression: \"%s\" not matching \"%s\"", matches, set.matches, cmd.Expression().String(), set.request)
		}
	}
}

func TestPosition(t *testing.T) {
	var data = []struct {
		command  string
		param    string
		position int
	}{
		{"command <lorem>", "lorem", 0},
		{"command <lorem> <ipsum> <dolor> <sit> <amet>", "dolor", 2},
		{"command <lorem> <ipsum> <dolor:integer> <sit> <amet>", "dolor", 2},
		{"command <lorem:integer> <ipsum> <dolor:integer> <sit> <amet>", "lorem", 0},
		{"command <lorem:string> <ipsum> <sit> <amet>", "lorem", 0},
	}

	var cmd Command
	var err error
	for _, set := range data {
		cmd, err = New(set.command, nil)
		checkErr(err, t)
		
		pos := cmd.Position(set.param)

		if pos != set.position {
			t.Errorf("Position() should be \"%d\", but is \"%d\"", set.position, pos)
		}
	}
}

func TestHas(t *testing.T) {
	var data = []struct {
		command   string
		parameter string
		has       bool
	}{
		{"command <lorem>", "lorem", true},
		{"command <lorem>", "dorem", false},
	}

	var cmd CommandInterface
	var err error
	for _, set := range data {
		cmd, err = New(set.command, nil)
		checkErr(err, t)

		has := cmd.Has(set.parameter)

		if has != set.has {
			t.Errorf("HasParameter is \"%v\", expected \"%v\"", has, set.has)
		}
	}
}

func TestParameters(t *testing.T) {
	var data = []struct {
		command    string
		parameters []string
	}{
		{"command <lorem>", []string{"lorem"}},
		{"cmd <lorem:string>", []string{"lorem"}},
		{"command <lorem:integer>", []string{"lorem"}},
		{"example <lorem> <ipsum> <dolor>", []string{"lorem", "ipsum", "dolor"}},
		{"command <lorem> <ipsum> <dolor:string>", []string{"lorem", "ipsum", "dolor"}},
		{"command <lorem> <ipsum:string> <dolor>", []string{"lorem", "ipsum", "dolor"}},
		{"command <lorem:string> <ipsum> <dolor>", []string{"lorem", "ipsum", "dolor"}},
		{"command <lorem:string> <ipsum> <dolor:string>", []string{"lorem", "ipsum", "dolor"}},
		{"command <lorem:string> <ipsum> <dolor:integer>", []string{"lorem", "ipsum", "dolor"}},
		{"command <lorem:integer> <ipsum:string> <dolor:integer>", []string{"lorem", "ipsum", "dolor"}},
	}

	var cmd Command
	var err error
	for _, set := range data {
		cmd, err = New(set.command, nil)
		checkErr(err, t)

		if cmd.Text() != set.command {
			t.Errorf("cmd.Text() must be \"%s\", but is \"%s\"", set.command, cmd.Text())
		}

		for _, param := range set.parameters {
			if !cmd.Has(param) {
				t.Errorf("\"%s\" missing parameter.Item \"%+v\"", cmd.Text(), param)
			}
		}
	}
}

func TestCustomTypes(t *testing.T) {
	var data = []struct {
		command string
		request string
		matches bool
		types Types
	}{
		{"command", "example", false, map[string]string{"string":  "[^\\s]+", "integer": "[0-9]+"}},
		{"command", "command", true, map[string]string{"string":  "[^\\s]+", "integer": "[0-9]+"}},
		{"command <custom:type>", "command 123", true, map[string]string{"string":  "[^\\s]+", "type": "[0-9]+"}},
		{"command <rest:rest>", "command something else", true, map[string]string{"string":  "[^\\s]+", "rest": ".*"}},
	}

	for _, set := range data {
		cmd, err := New(set.command, set.types)
		checkErr(err, t)

		matches := cmd.Matches(set.request)

		if matches != set.matches {
			t.Errorf("Matches() returns unexpected values. Got \"%v\", expected \"%v\"\nExpression: \"%s\" not matching \"%s\"", matches, set.matches, cmd.Expression().String(), set.request)
		}
	}
}