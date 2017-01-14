package utils

import (
	"regexp"
	"strings"
)

type Attribute map[string]interface{}

func (attrs Attribute) Item(name string, value ...interface{}) (attr interface{}) {

	attr = attrs[name]

	if value != nil && len(value) > 0 {
		attr = value[0]
		if attr == nil {
			delete(attrs, name)
		} else {
			attrs[name] = attr
		}
	}

	return
}

func (attrs Attribute) Set(key string, value interface{}) interface{} {
	attrs[key] = value
	return value
}

func (attrs Attribute) Get(key string) interface{} {
	if value, ok := attrs[key]; ok {
		return value
	}
	return nil
}

func (attrs Attribute) GetString(key string) string {
	return ToString(attrs.Get(key))
}

func (props Attribute) SetProperty(key string, value string, keepToken ...bool) string {
	if len(keepToken) <= 0 || !keepToken[0] {
		reg := regexp.MustCompile(`\{\s*([\w\._\-]+)\s*\}`)
		value = reg.ReplaceAllStringFunc(value, func(s string) string {

			skey := strings.Trim(s, "{")
			skey = strings.Trim(skey, "}")
			skey = strings.Trim(skey, " ")

			if v := props.Get(skey); v != nil {
				return ToString(v)
			}

			return s
		})
	}
	// props.Set(key, value)
	props[key] = value

	return value
}
