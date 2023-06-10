package parser

import "strings"

// 类型集合
type TT map[string]struct{}

func (tt *TT) Has(key string) bool {
	for k := range *tt {
		if key == k {
			return true
		}
	}
	return false
}

func (tt *TT) Join() string {
	keys := []string{}
	for k := range *tt {
		keys = append(keys, k)
	}
	return strings.Join(keys, "|")
}

func NewTypes(keys ...string) *TT {
	tt := make(TT)
	for _, k := range keys {
		tt[k] = struct{}{}
	}
	return &tt
}

var TokenTypes = NewTypes("int", "str", "bool", "float")
var RequestTypes = NewTypes("GET", "POST")
