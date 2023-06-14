package parser

import (
	"strings"

	"github.com/Drelf2020/utils"
)

// 单条语句
//
// 其中 Args 表示其 Name 中携带的参数
//
// 当该语句为字典时 Tokens 不为空 表示其下包含的语句
type Token struct {
	Type   string            `json:"type,omitempty"`
	Name   string            `json:"-"`
	Hint   string            `json:"hint,omitempty"`
	Value  string            `json:"value,omitempty"`
	Parent *Token            `json:"-"`
	Args   []string          `json:"-"`
	Tokens map[string]*Token `json:"values,omitempty"`
	open   bool
}

// 去除参数的纯类型
func (token *Token) PureType() string {
	if t := Slice(token.Type, "", "<", 0); t != "" {
		return t
	}
	return token.Type
}

// 判断该语句是否为 Api 起始语句
func (token *Token) IsApi() bool {
	return MethodTypes.Has(token.Type)
}

// 判断该语句是否为 type 起始语句
func (token *Token) IsType() bool {
	return token.Type == "type"
}

// 判断该语句是否为在变量类型内
func (token *Token) InTypes() bool {
	return VarTypes.Has(token.PureType())
}

// 判断该语句是否为起始括号
func (token *Token) IsOpen() bool {
	return token.open
}

// 判断该语句是否为闭合括号
func (token *Token) IsClose() bool {
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

// 添加子语句
func (token *Token) Add(t *Token) {
	t.Parent = token
	token.Tokens[t.Name] = t
}

// 移除子语句
func (token *Token) Pop(name string) *Token {
	if v, ok := token.Tokens[name]; ok {
		delete(token.Tokens, name)
		return v
	}
	return nil
}

// 转字典
func (token *Token) ToDict() map[string]string {
	dic := make(map[string]string)
	for k, v := range token.Tokens {
		dic[k] = v.Value
	}
	return dic
}

// 解析类似 res<T1,T2> 中的子类型 T1 T2
//
// data 顺序 Type Name Hint Value
func NewToken(data ...string) *Token {
	// 自动推断类型
	typ := data[0]
	val := data[3]
	tokens := make(map[string]*Token)
	if typ == "" || typ == "auto" {
		typ = AutoType(val)
	} else {
		var as []string
		typ, as = NameSlice(typ)
		if token := VarTypes.Get(typ); token != nil {
			tokens = token.Tokens
			for i, a := range as {
				ForMap(
					tokens,
					func(s string, t *Token) {
						t.Type = a
						// 这里还要替换值
					},
					func(s string, t *Token) bool { return t.Type == token.Args[i] },
				)
			}
		}
	}
	if val == "{" {
		val = ""
	}

	// 处理变量名及参数
	name, args := NameSlice(data[1])

	// 去除标注前后空白
	hint := strings.Trim(data[2], " ")

	return &Token{typ, name, hint, val, nil, args, tokens, data[3] == "{"}
}
