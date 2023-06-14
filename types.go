package parser

import (
	"regexp"
	"strings"
)

// 类型集合
type Types map[string]*Tokens

// 添加类型
func (types Types) Add(key string, args ...string) {
	types[key] = NewTokens(key, args...)
}

// 判断字段
func (types *Types) Has(key string) bool {
	for k := range *types {
		if key == k {
			return true
		}
	}
	return false
}

// 以字符串形式连接多个 type
func (types *Types) Join() string {
	keys := []string{}
	for k := range *types {
		keys = append(keys, k)
	}
	return strings.Join(keys, "|")
}

// 合并多个类型组
func (ts *Types) Union(typess ...*Types) *Types {
	nt := *ts
	for _, types := range typess {
		for k := range *types {
			nt.Add(k)
		}
	}
	return &nt
}

// 正则查找
func (types *Types) FindTokens(api string) (tokens []*Token) {
	re := regexp.MustCompile(` *(?:((?:` + types.Join() + `)<?(?:` + types.Join() + `)?>?) )? *([^:^=^\r^\n^ ]+)(?:: *([^=^\r^\n]+))? *(?:= *([^\r^\n]+))?`)
	for _, sList := range re.FindAllStringSubmatch(api, -1) {
		tokens = append(tokens, NewToken(sList[1:]))
	}
	return
}

// 获取 Tokens
func (types *Types) Get(key string) *Tokens {
	return (*types)[key]
}

// 构造函数
func NewTypes(keys ...string) *Types {
	types := make(Types)
	for _, k := range keys {
		types.Add(k)
	}
	return &types
}

// 支持的变量类型 auto str num bool
var VarTypes = NewTypes("type", "auto", "str", "num", "bool")

// 支持的请求类型 GET POST
var MethodTypes = NewTypes("GET", "POST")

// 支持的请求字段名 data params headers cookies
var RequestTypes = NewTypes("data", "params", "headers", "cookies")
