// See LICENSE.txt for licensing information.

// Package simpleini implements yet another interface to a particular subset of INI files.
// In particular it doesn't allow 'out-of-section' keys, has no required sections, is
// line-ending agnostic and provides dedicated getters and setter for the basic data types.
package simpleini

import (
	"bufio"
	"fmt"
	"io"
	"sort"
	"strconv"
	"strings"
)

// Internal constants.
const (
	DICTSIZE = 8       // default size of the dictionary
	VERSION  = "0.0.1" // current version of the library
)

// INI represents a INI file data.
// Currently it consist only of the raw dictionary. It has been defined as a struct to
// make future additions backward compatible (e.g. original file path).
type INI struct {
	dict map[string]map[string]string
}

// NewINI returns an INI with default settings.
func NewINI() *INI {
	return &INI{dict: make(map[string]map[string]string, DICTSIZE)}
}

// Parse tries to parse an input into an INI.
func Parse(input io.Reader) (*INI, error) {
	scn := bufio.NewScanner(input)
	ini := NewINI()

	var section string
	lineNum := 1
	for scn.Scan() {
		line := strings.Trim(scn.Text(), " ")
		if len(line) < 1 {
			continue
		}
		switch line[0] {
		case ';':
		case '[':
			if len(line) < 3 {
				return nil, fmt.Errorf("Line %d: Malformed section", lineNum)
			}
			if line[len(line)-1] != ']' {
				return nil, fmt.Errorf("Line %d: Malformed section", lineNum)
			}
			section = line[1 : len(line)-1]
			if _, present := ini.dict[section]; present {
				return nil, fmt.Errorf("Line %d: Section '%s' has been defined previosuly", lineNum, section)
			}
			ini.dict[section] = make(map[string]string, DICTSIZE)
		default:
			if section == "" {
				return nil, fmt.Errorf("Line %d: Property defined outside of a section", lineNum)
			}
			parts := strings.Split(line, "=")
			if len(parts) != 2 {
				return nil, fmt.Errorf("Line %d: Malformed property", lineNum)
			}
			property := strings.Trim(parts[0], " ")
			if _, present := ini.dict[section][property]; present {
				return nil, fmt.Errorf("Line %d: Property '%s' has been defined previously", lineNum, property)
			}
			ini.dict[section][property] = strings.Trim(parts[1], " ")
		}
		lineNum++
	}
	return ini, nil
}

// Sections return a slice of sections from an INI.
func (i *INI) Sections() []string {
	var sections []string
	for s := range i.dict {
		sections = append(sections, s)
	}
	return sections
}

// Properties returns a slice of properties (keys) form a given section of an INI.
func (i *INI) Properties(section string) ([]string, error) {
	properties, present := i.dict[section]
	if !present {
		return nil, fmt.Errorf("Section '%s' not found", section)
	}
	var ps []string
	for p := range properties {
		ps = append(ps, p)
	}
	return ps, nil
}

// GetString tries to return a string representation from a section - property pair of an INI.
func (i *INI) GetString(section string, property string) (string, error) {
	properties, present := i.dict[section]
	if !present {
		return "", fmt.Errorf("Section '%s' not found", section)
	}
	value, present := properties[property]
	if !present {
		return "", fmt.Errorf("Property '%s' not found in section '%s'", property, section)
	}
	return value, nil
}

// GetInt tries to return an integer representation from a section - property pair of an INI.
func (i *INI) GetInt(section string, property string) (int, error) {
	strVal, err := i.GetString(section, property)
	if err != nil {
		return 0, err
	}
	intVal, err := strconv.Atoi(strVal)
	if err != nil {
		return 0, fmt.Errorf("Property '%s/%s' is not an int: %s", section, property, err)
	}
	return intVal, err
}

// GetBool tries to return a bool representation from a section - property pair of an INI.
// Values that map to true: 'true', 'yes', 'on'; and respective values map to false.
func (i *INI) GetBool(section string, property string) (bool, error) {
	strVal, err := i.GetString(section, property)
	if err != nil {
		return false, err
	}
	switch strVal {
	case "true", "yes", "on":
		return true, nil
	case "false", "no", "off":
		return false, nil
	default:
		return false, fmt.Errorf("Property '%s/%s' is not a boolean", section, property)
	}
}

// SetString sets a section - property pair to the given value, creating it if it wasn't
// already present.
func (i *INI) SetString(section string, property string, value string) {
	properties, present := i.dict[section]
	if !present {
		properties = make(map[string]string, DICTSIZE)
		i.dict[section] = properties
	}
	properties[property] = value
	return
}

// SetInt sets a section - property pair to the given value, creating it if it wasn't
// already present.
func (i *INI) SetInt(section string, property string, value int) {
	i.SetString(section, property, strconv.Itoa(value))
	return
}

// SetBool sets a section - property pair to the given value, creating it if it wasn't
// already present.
func (i *INI) SetBool(section string, property string, value bool) {
	var strVal string
	if value {
		strVal = "yes"
	} else {
		strVal = "no"
	}
	i.SetString(section, property, strVal)
	return
}

// Write tries to output the INI onto an output.
// The prettify option sorts the sections and properties within for better user
// experience.
func (i *INI) Write(output io.Writer, prettify bool) error {
	buf := bufio.NewWriter(output)
	sections := i.Sections()
	if prettify {
		sort.Strings(sections)
	}
	for _, section := range sections {
		buf.WriteString("[" + section + "]" + CRLF)
		properties, _ := i.Properties(section)
		if prettify {
			sort.Strings(properties)
		}
		for _, property := range properties {
			buf.WriteString(property + " = " + i.dict[section][property] + CRLF)
		}
		if prettify {
			buf.WriteString(CRLF)
		}
	}
	return buf.Flush()
}

func DeleteSection(ini *INI, section string) {
	delete(ini.dict, section)
}

func DeleteProperty(ini *INI, section string, property string) {
	properties, present := ini.dict[section]
	if present {
		delete(properties, property)
	}
}
