package sql

import (
	"regexp"
	"strings"

	"github.com/licat233/genzero/config"
	"github.com/licat233/genzero/tools"
)

func TableFilterCondition() (mustFilter bool, checkTables []string) {
	if len(config.C.DB.Tables) == 0 {
		return false, nil
	}
	tables := []string{}
	for _, t := range config.C.DB.Tables {
		if t == "" {
			continue
		}
		tables = append(tables, t)
	}
	if len(tables) == 1 {
		if tables[0] == "*" {
			return false, nil
		}
	}
	return len(tables) != 0, tables
}

func goType(sqlType string) string {
	if strings.Contains(sqlType, "(") {
		sqlType = strings.Split(sqlType, "(")[0]
	}
	v, ok := TypeForMysqlToGo[sqlType]
	if ok {
		return v
	}
	return "any"
}

func MutipleFindStringSubmatch(line string, patterns ...string) (string, error) {
	if len(patterns) == 0 {
		return "", nil
	}
	for _, pattern := range patterns {
		reg, err := regexp.Compile(pattern)
		if err != nil {
			return "", err
		}
		arr := reg.FindStringSubmatch(line)
		if len(arr) == 2 {
			return strings.TrimSpace(arr[1]), nil
		}
	}
	return "", nil
}

func MutipleStringSubmatch(line string, patterns ...string) (bool, error) {
	line = strings.TrimSpace(line)
	if len(patterns) == 0 {
		return false, nil
	}
	for _, pattern := range patterns {
		if len(pattern) == 0 {
			continue
		}
		reg, err := regexp.Compile(pattern)
		if err != nil {
			return false, err
		}
		ok := reg.MatchString(line)
		if ok {
			return true, nil
		}
	}
	return false, nil
}

func PickTableComment(line string) string {
	res, err := MutipleFindStringSubmatch(line, `\sCOMMENT\s*=\s*"(.*?)"`, `\sCOMMENT\s*=\s*'(.*?)'`, `\scomment\s*=\s*"(.*?)"`, `\scomment\s*=\s*'(.*?)'`)
	if err != nil {
		tools.Error(err.Error())
	}
	return res
}

func PickFieldComment(line string) string {
	res, err := MutipleFindStringSubmatch(line, `\sCOMMENT\s+"(.*?)"`, `\sCOMMENT\s+'(.*?)'`, `\scomment\s+"(.*?)"`, `\scomment\s+'(.*?)'`)
	if err != nil {
		tools.Error(err.Error())
	}
	return res
}

func PickFieldDefaultValue(line string) string {
	res, err := MutipleFindStringSubmatch(line, `\sDEFAULT\s+"(.*?)"`, `\sDEFAULT\s+'(.*?)'`, `\sdefault\s+([^']*?)\s`, `\sdefault\s+'(.*?)'`)
	if err != nil {
		tools.Error(err.Error())
	}
	return res
}

func PickFieldName(line string) string {
	res, err := MutipleFindStringSubmatch(line, `^\s*`+"`"+`(\w+)`+"`"+`\s\w`)
	if err != nil {
		tools.Error(err.Error())
	}
	return res
}

func PickUniqueFieldName(line string) string {
	res, err := MutipleFindStringSubmatch(line, "\\s*UNIQUE\\s*\\(`([^`]+)`\\)", "\\s*unique\\s*\\(`([^`]+)`\\)")
	if err != nil {
		tools.Error(err.Error())
	}
	return res
}

func PickFieldType(line string) string {
	res, err := MutipleFindStringSubmatch(line, `^.+?\s([A-Za-z]+)\W?`)
	if err != nil {
		tools.Error(err.Error())
	}
	return res
}

func IsFieldDefineString(line string) bool {
	reg, err := regexp.Compile(`^\s*` + "`" + `\w+` + "`" + `\s([A-Za-z]+)\W?`)
	if err != nil {
		tools.Error(err.Error())
		return false
	}
	ok := reg.MatchString(line)
	return ok
}

func IsDeleteField(fieldName string) bool {
	return tools.SliceContain(DelFieldNames, fieldName)
}

func IsDelAtField(fieldName string) bool {
	return tools.SliceContain(DelAtFieldNames, fieldName)
}

func IsUuidField(fieldName string) bool {
	return tools.SliceContain(UuidFieldNames, fieldName)
}

func IsIgnoreField(fieldName string) bool {
	return tools.SliceContain(config.C.DB.IgnoreColumns, fieldName)
}

func IsNameField(fieldName string) bool {
	return tools.SliceContain(NameFieldNames, fieldName)
}
