package parser

import (
	"regexp"
	"strings"
)

// 类型集合
type Types map[string]*Token

// 添加类型
func (types Types) Add(token *Token, names ...string) {
	for _, n := range names {
		types[n] = nil
	}
	if token != nil {
		types[token.Name] = token
	}
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
		for k, v := range *types {
			nt.Add(v, k)
		}
	}
	return &nt
}

// 生成正则表达式
func (types *Types) ToRegexp() *regexp.Regexp {
	return regexp.MustCompile(` *(?:((?:` + types.Join() + `)<?(?:` + types.Join() + `)?>?) )? *([^:^=^\r^\n^ ]+)(?:: *([^=^\r^\n]+))? *(?:= *([^\r^\n]+))?`)
}

// 正则查找字符串
func (types *Types) FindStrings(api string) [][]string {
	return types.ToRegexp().FindAllStringSubmatch(api, -1)
}

// 正则查找语句
func (types *Types) FindTokens(api string) (tokens []*Token) {
	re := types.ToRegexp()
	for _, sList := range re.FindAllStringSubmatch(api, -1) {
		tokens = append(tokens, NewToken(sList[1:]...))
	}
	return
}

// 获取 Token
//
// 当 key 为基础类型(str num bool)时返回 nil
func (types *Types) Get(key string) *Token {
	return (*types)[key]
}

// 构造函数
func NewTypes(keys ...string) *Types {
	types := make(Types)
	types.Add(nil, keys...)
	return &types
}

// 支持的变量类型 auto str num bool
var VarTypes = NewTypes("type", "auto", "str", "num", "bool")

// 支持的请求类型 GET POST
var MethodTypes = NewTypes("GET", "POST")
