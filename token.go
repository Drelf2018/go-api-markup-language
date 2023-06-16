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
	Type   string            `json:"type,omitempty" yaml:"type,omitempty"`
	Name   string            `json:"-" yaml:"-"`
	Hint   string            `json:"hint,omitempty" yaml:"hint,omitempty"`
	Value  string            `json:"-" yaml:"-"`
	Output any               `json:"value,omitempty" yaml:"value,omitempty"`
	Parent *Token            `json:"-" yaml:"-"`
	Args   []string          `json:"-" yaml:"-"`
	Tokens map[string]*Token `json:"-" yaml:"-"`
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
	return VarTypes.Has(token.Type)
}

// 判断该语句是否为起始括号
func (token *Token) IsOpen() bool {
	return token.Value == "{"
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

// 修改类型
func (token *Token) SetValue(args []string) {
	if tk := VarTypes.Get(token.Type); tk != nil {
		if len(args) != len(tk.Args) {
			panic(token.Type + " 的参数个数都数歪来？")
		}
		argsMap := map[string]string{"str": "str", "num": "num", "bool": "bool"}
		for i, arg := range tk.Args {
			argsMap[arg] = args[i]
		}
		ForMap(
			tk.Tokens,
			func(s string, t *Token) { token.Tokens[s] = NewToken(argsMap[t.Type], t.Name, t.Hint, t.Value) },
		)
		token.Output = token.Tokens
	} else if token.IsOpen() {
		token.Output = token.Tokens
	}
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
	typ, trueVal := AutoType(data[0], data[3])
	// 解析类型中的参数
	typ, params := NameSlice(typ)
	// 解析变量名中的参数
	name, args := NameSlice(data[1])
	// 去除标注前后空白
	hint := strings.Trim(data[2], " ")

	token := Token{typ, name, hint, data[3], trueVal, nil, args, make(map[string]*Token)}
	token.SetValue(params)

	return &token
}
