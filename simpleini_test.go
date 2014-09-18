package simpleini

import (
	"io"
	"strings"
	"testing"
)

func TestParsingCorrectInputs(t *testing.T) {
	var input io.Reader

	inputs := [...]string{
		`[main]
string = this is a test
integer = 123
boolean = yes

[auxillary]
whatever = something
`,
		`; some comment
; another
    [main]
  string   =  this is a test   
integer =   0123

;integer = 9999
 boolean =  true
; boolean = false

;[auxillary]
[auxillary]

       whatever =  something

`,
	}
	for idx, i := range inputs {
		idx++
		input = strings.NewReader(i)

		ini, err := Parse(input)
		if err != nil {
			t.Errorf("(%d) Parsing failed: %s", idx, err)
			return
		}
		if len(ini.Sections()) != 2 {
			t.Errorf("(%d) Wrong number of sections", idx)
		}
		ps, err := ini.Properties("main")
		if err != nil {
			t.Errorf("(%d) Couldn't get properties for main: %s", idx, err)
		} else {
			if len(ps) != 3 {
				t.Errorf("(%d) Wrong number of properties", idx)
			}
		}

		strVal, err := ini.GetString("main", "string")
		if err != nil {
			t.Errorf("(%d) Couldn't get main/string: %s", idx, err)
		} else {
			if strVal != "this is a test" {
				t.Errorf("(%d) Wrong value for main/string", idx)
			}
		}
		strVal, err = ini.GetString("auxillary", "whatever")
		if err != nil {
			t.Errorf("(%d) Couldn't get auxillary/whatever: %s", idx, err)
		} else {
			if strVal != "something" {
				t.Errorf("(%d) Wrong value for auxillar/whatever", idx)
			}
		}

		intVal, err := ini.GetInt("main", "integer")
		if err != nil {
			t.Errorf("(%d) Couldn't get main/integer: %s", idx, err)
		} else {
			if intVal != 123 {
				t.Errorf("(%d) Wrong value for main/integer", idx)
			}
		}

		boolVal, err := ini.GetBool("main", "boolean")
		if err != nil {
			t.Errorf("(%d) Couldn't get main/boolean: %s", idx, err)
		} else {
			if !boolVal {
				t.Errorf("(%d) Wrong value for main/boolean", idx)
			}
		}
	}
}

func TestKeyWithSpaces(t *testing.T) {
	input := strings.NewReader(`
[ugly]
  key with spaces   =   value with spaces too   `)
	ini, err := Parse(input)
	if err != nil {
		t.Errorf("Parsing failed: %s", err)
		return
	}
	val, err := ini.GetString("ugly", "key with spaces")
	if err != nil {
		t.Errorf("Couldn't get ugly/key with spaces: %s", err)
	} else {
		if val != "value with spaces too" {
			t.Errorf("Wrong value for ugly/key with spaces")
		}
	}
}

func TestBoolValues(t *testing.T) {
	trueVals := [...]string{"true", "yes", "on"}
	falseVals := [...]string{"false", "no", "off"}

	for idx, val := range trueVals {
		input := strings.NewReader("[bools]\nval = " + val)
		ini, err := Parse(input)
		if err != nil {
			t.Errorf("(%d) Parsing failed: %s", idx, err)
			return
		}
		boolVal, err := ini.GetBool("bools", "val")
		if err != nil {
			t.Errorf("(%d) Couldn't get bools/val: %s", idx, err)
		} else {
			if !boolVal {
				t.Errorf("(%d) Value of bools/val should have been true", idx)
			}
		}
	}

	for idx, val := range falseVals {
		input := strings.NewReader("[bools]\nval = " + val)
		ini, err := Parse(input)
		if err != nil {
			t.Errorf("(%d) Parsing failed: %s", idx, err)
			return
		}
		boolVal, err := ini.GetBool("bools", "val")
		if err != nil {
			t.Errorf("(%d) Couldn't get bools/val: %s", idx, err)
		} else {
			if boolVal {
				t.Errorf("(%d) Value of bools/val should have been flase", idx)
			}
		}
	}
}

func TestParsingMalformedSections(t *testing.T) {
	inputs := [...]string{"[", "[aaa", "[]"}
	for idx, i := range inputs {
		idx++
		input := strings.NewReader(i)

		_, err := Parse(input)
		if err == nil {
			t.Errorf("(%d) Malformed section parsed", idx)
		}
	}
}

func TestParsingRedfinedSection(t *testing.T) {
	input := strings.NewReader("[a]\nval=1\n\n[a]\nval=1\n")

	_, err := Parse(input)
	if err == nil {
		t.Errorf("Redefined section parsed")
	}
}

func TestParsingMalformedProperties(t *testing.T) {
	inputs := [...]string{"whatever", "this = that", "[hello]\nthis", "something\n[section]"}
	for idx, i := range inputs {
		idx++
		input := strings.NewReader(i)

		_, err := Parse(input)
		if err == nil {
			t.Errorf("(%d) Malformed property parsed", idx)
		}
	}
}

func TestParsingRedfinedProperty(t *testing.T) {
	input := strings.NewReader("[a]\nval=1\nval=2\n")

	_, err := Parse(input)
	if err == nil {
		t.Errorf("Redefined property parsed")
	}
}

func TestBadProperties(t *testing.T) {
	input := strings.NewReader("[a]\nval=1")

	ini, _ := Parse(input)
	_, err := ini.Properties("b")
	if err == nil {
		t.Errorf("Not defined properties found")
	}
}

func TestBadGetString(t *testing.T) {
	input := strings.NewReader("[a]\nval=1")

	ini, _ := Parse(input)
	_, err := ini.GetString("b", "a")
	if err == nil {
		t.Errorf("Not defined section found")
	}
	_, err = ini.GetString("a", "non")
	if err == nil {
		t.Errorf("Not defined property found")
	}
}

func TestBadGetInt(t *testing.T) {
	input := strings.NewReader("[a]\nsval=foo")

	ini, _ := Parse(input)
	_, err := ini.GetInt("b", "a")
	if err == nil {
		t.Errorf("Not defined section found")
	}
	_, err = ini.GetInt("a", "non")
	if err == nil {
		t.Errorf("Not defined property found")
	}
	_, err = ini.GetInt("a", "sval")
	if err == nil {
		t.Errorf("Non-int property returned as int")
	}
}

func TestBadGetBool(t *testing.T) {
	input := strings.NewReader("[a]\nsval=foo")

	ini, _ := Parse(input)
	_, err := ini.GetBool("b", "a")
	if err == nil {
		t.Errorf("Not defined section found")
	}
	_, err = ini.GetBool("a", "non")
	if err == nil {
		t.Errorf("Not defined property found")
	}
	_, err = ini.GetBool("a", "sval")
	if err == nil {
		t.Errorf("Non-int property returned as bool")
	}
}
