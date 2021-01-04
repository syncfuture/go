package u

import (
	"reflect"
	"regexp"
	"strings"
)

const (
	_regStr = `^([A-Za-z0-9+/]{4})*([A-Za-z0-9+/]{3}=|[A-Za-z0-9+/]{2}==)?$`
)

var (
	_regex = regexp.MustCompile(_regStr)
)

func IsMissing(x interface{}) bool {
	if x == nil {
		return true
	}

	v := reflect.ValueOf(x)
	k := v.Kind()

	for {
		if k == reflect.Ptr || k == reflect.Interface {
			v = v.Elem()
			k = v.Kind()
		} else {
			break
		}
	}

	if k == reflect.Invalid { // nil
		return true
	} else if k == reflect.String {
		return strings.TrimSpace(x.(string)) == ""
	} else if k == reflect.Slice || k == reflect.Array || k == reflect.Chan || k == reflect.Map {
		return v.Len() == 0
	}

	return false
}

func IsBase64String(str string) bool {
	return _regex.MatchString(str)
}
