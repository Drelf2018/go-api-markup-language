package parser

import (
	"strings"

	"github.com/Drelf2020/utils"
)

// 单条语句
type Token struct {
	Type  string `json:"type,omitempty"`
	Name  string `json:"-"`
	Hint  string `json:"hint,omitempty"`
	Value string `json:"value,omitempty"`
}

// 去除参数的纯类型
func (token *Token) PureType() string {
	if t := Slice(token.Type, "", "<", 0); t != "" {
		return t
	}
	return token.Type
}

// 去除参数的纯名字
func (token *Token) PureName() string {
	if t := Slice(token.Name, "", "<", 0); t != "" {
		return t
	}
	return token.Name
}

// 判断该语句是否为 Api 起始语句
func (token *Token) IsApi() bool {
	return MethodTypes.Has(token.Type)
}

// 判断该语句是否为 type 起始语句
func (token *Token) IsType() bool {
	return token.Type == "type"
}

// 判断该语句是否为变量类型
func (token *Token) InType() bool {
	return VarTypes.Has(token.PureType())
}

// 判断该语句是否为闭合括号
func (token *Token) IsClosed() bool {
	return token.Name == "}"
}

// 判断该语句是否为必要变量
func (token *Token) IsRequired() bool {
	return token.Value == ""
}

// 判断该语句是否为常量
func (token *Token) IsConstant() bool {
	return utils.Endswith(token.Value, ",constant")
}

// 判断该语句是否为选填变量
func (token *Token) IsOptional() bool {
	return !token.IsRequired() && !token.IsConstant()
}

// 转为 Python 格式
func (token *Token) ToPython() (s string) {
	s = token.Name
	if token.Type != "" {
		s += ": " + token.Type
	}
	if token.Value != "" {
		if token.Type == "str" || token.Type == "" {
			s += " = \"" + token.Value + "\""
		} else if token.Type == "bool" {
			s += " = " + Capitalize(token.Value)
		} else {
			s += " = " + token.Value
		}
	}
	return
}

func NewToken(data []string) *Token {
	data[2] = strings.Trim(data[2], " ")
	return &Token{data[0], data[1], data[2], data[3]}
}

type Tokens struct {
	// 类型名
	Name string
	// 参数名
	Args []string
	// 子语句
	Tokens map[string]*Token
}

// 添加数据
func (ts Tokens) Add(token *Token) {
	ts.Tokens[token.Name] = token
}

// 转字典
func (ts Tokens) ToDict() map[string]string {
	dic := make(map[string]string)
	for k, v := range ts.Tokens {
		dic[k] = v.Value
	}
	return dic
}

func NewTokens(name string, args ...string) *Tokens {
	return &Tokens{
		name,
		args,
		make(map[string]*Token),
	}
}

func MakeTokens() Tokens {
	return Tokens{
		"",
		[]string{},
		make(map[string]*Token),
	}
}
