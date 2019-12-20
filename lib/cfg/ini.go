package cfg

import (
	"bytes"
	"errors"
	"io/ioutil"
	"strconv"
)

const (
	DefaultLinkeSplit = "\n"
	DefaultComments   = "#"
	DefaultCommentss  = ";"
	DefaultSep        = "="
)

type iniSection struct {
	name    []string
	section map[string]string
}

type IniConfig struct {
	FileName string
	iniSection
	sectionMap map[string]*iniSection
}

func NewIniConfig(filename string) (*IniConfig, error) {
	i := new(IniConfig)
	i.FileName = filename
	i.sectionMap = make(map[string]*iniSection)
	i.section = make(map[string]string)
	err := i.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return i, nil
}

func (i *IniConfig) ReadFile(filename string) error {
	contents, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	return i.parseIni(contents)
}

func (i *IniConfig) parseIni(contents []byte) error {
	lines := bytes.Split(contents, []byte(DefaultLinkeSplit))

	var sectionName string
	var section *iniSection

	for _, line := range lines {
		line = bytes.TrimSpace(line)
		if len(line) == 0 {
			continue
		}

		if bytes.HasPrefix(line, []byte(DefaultComments)) || bytes.HasPrefix(line, []byte(DefaultCommentss)) {
			continue
		}

		if bytes.HasPrefix(line, []byte("[")) {
			if !bytes.HasSuffix(line, []byte("]")) {
				return errors.New("section names must be surrounded by [ and ], as in [section]")
			}
			sectionName = string(line[1 : len(line)-1])

			if secet, ok := i.sectionMap[sectionName]; !ok {
				// reset data map
				section = new(iniSection)
				section.name = []string{sectionName}
				section.section = make(map[string]string)
				i.sectionMap[sectionName] = section
			} else {
				i.sectionMap[sectionName] = secet
			}
			continue
		}

		index := bytes.Index(line, []byte(DefaultSep))
		if index <= 0 {
			err := errors.New("Came accross an error : " + string(line) + " is NOT a valid key/value pair")
			return err
		}

		key := bytes.TrimSpace(line[0:index])
		value := bytes.TrimSpace(line[index+len(DefaultSep):])

		section.name = append(section.name, string(key))
		section.section[string(key)] = string(value)

		i.section[string(key)] = string(value)
	}
	return nil
}

func get(values map[string]string, key string) string {
	if len(key) == 0 || values == nil {
		return ""
	}
	if data, ok := values[key]; ok {
		return data
	}
	return ""
}

func GetString(values map[string]string, key string) string {
	val := get(values, key)
	if len(val) > 0 {
		return val
	}
	return ""
}

func GetInt64(values map[string]string, key string) int64 {
	val := get(values, key)
	if len(val) > 0 {
		integer, err := strconv.ParseInt(val, 10, 64)
		if err != nil {
			return 0
		}
		return integer
	}
	return 0
}

func GetInt(values map[string]string, key string) int {
	val := get(values, key)
	if len(val) > 0 {
		inte, err := strconv.Atoi(val)
		if err != nil {
			return 0
		}
		return inte
	}
	return 0
}

func GetFloat(values map[string]string, key string) float64 {
	val := get(values, key)
	if len(val) > 0 {
		floa, err := strconv.ParseFloat(val, 64)
		if err != nil {
			return float64(0)
		}
		return floa
	}
	return float64(0)
}

func (i *IniConfig) IsExist(key string) bool {
	val := get(i.section, key)
	if len(val) > 0 {
		return true
	}
	return false
}

func (i *IniConfig) String(key string) string {
	return GetString(i.section, key)
}

func (i *IniConfig) Int64(key string) int64 {
	return GetInt64(i.section, key)
}

func (i *IniConfig) Int(key string) int {
	return GetInt(i.section, key)
}

func (i *IniConfig) Float(key string) float64 {
	return GetFloat(i.section, key)
}

func (i *IniConfig) GetSectionMap() map[string]*iniSection {
	return i.sectionMap
}

func (i *IniConfig) GetSectionName(sectionkey string) []string {
	if len(sectionkey) == 0 {
		return nil
	}
	inisection, ok := i.sectionMap[sectionkey]
	if !ok {
		return nil
	}
	return inisection.name
}

func (i *IniConfig) GetGlobalSection() map[string]string {
	return i.section
}

func (i *IniConfig) GetGlobalName() []string {
	return i.name
}
