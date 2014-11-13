package simpleini

import (
	"bufio"
	"fmt"
	"io"
	"sort"
	"strconv"
	"strings"
)

const (
	DICTSIZE = 8
	VERSION  = "0.0.1"
)

type INI struct {
	dict map[string]map[string]string
}

func NewINI() *INI {
	return &INI{dict: make(map[string]map[string]string, DICTSIZE)}
}

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

func (i *INI) Sections() []string {
	var sections []string
	for s := range i.dict {
		sections = append(sections, s)
	}
	return sections
}

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

func (i *INI) SetString(section string, property string, value string) {
	properties, present := i.dict[section]
	if !present {
		properties = make(map[string]string, DICTSIZE)
		i.dict[section] = properties
	}
	properties[property] = value
	return
}

func (i *INI) SetInt(section string, property string, value int) {
	i.SetString(section, property, strconv.Itoa(value))
	return
}

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
