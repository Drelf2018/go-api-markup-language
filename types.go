package parser

import "strings"

// 类型集合
type Types map[string]struct{}

func (types *Types) Has(key string) bool {
	for k := range *types {
		if key == k {
			return true
		}
	}
	return false
}

func (types *Types) Join() string {
	keys := []string{}
	for k := range *types {
		keys = append(keys, k)
	}
	return strings.Join(keys, "|")
}

func NewTypes(keys ...string) *Types {
	types := make(Types)
	for _, k := range keys {
		types[k] = struct{}{}
	}
	return &types
}

// 支持的变量类型 int str bool float
var VarTypes = NewTypes("int", "str", "bool", "float")

// 支持的请求类型 GET POST
var MethodTypes = NewTypes("GET", "POST")

// 支持的请求字段名
var RequestTypes = NewTypes("data", "params", "cookies", "headers")
