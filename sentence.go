package aml

import (
	"fmt"

	"github.com/Drelf2020/utils"
)

// 单条语句
//
// 其中 Args 表示其 Name 中携带的参数
//
// 当该语句为字典时 Map 不为空 表示其下包含的语句
type Sentence struct {
	Type   string   `json:"type,omitempty" yaml:"type,omitempty"`
	Name   string   `json:"-" yaml:"-"`
	Hint   string   `json:"hint,omitempty" yaml:"hint,omitempty"`
	Value  string   `json:"-" yaml:"-"`
	Args   []string `json:"-" yaml:"-"`
	base   int
	index  int
	parent *Sentence
	Length int64                `json:"length,omitempty" yaml:"length,omitempty"`
	Output any                  `json:"value,omitempty" yaml:"value,omitempty"`
	List   []*Sentence          `json:"-" yaml:"-"`
	Map    map[string]*Sentence `json:"-" yaml:"-"`
}

// 判断该语句是否为 Api 起始语句
func (sentence *Sentence) IsApi() bool {
	switch GetKind(sentence.Type) {
	case GET, POST:
		return true
	}
	return false
}

// 判断该语句是否为起始括号
func (sentence *Sentence) IsBrace() bool {
	return sentence.Value == "{"
}

// 判断该语句是否为起始列表
func (sentence *Sentence) IsBracket() bool {
	return sentence.Value == "["
}

// 判断该语句是否为起始枚举
func (sentence *Sentence) IsGroup() bool {
	return sentence.Value == "("
}

// 判断该语句是否未闭合
func (sentence *Sentence) IsOpen() bool {
	return sentence.IsBrace() || sentence.IsBracket() || sentence.IsGroup()
}

// 判断该语句是否为字典
func (sentence *Sentence) IsDict() bool {
	return sentence.IsBrace() || sentence.base == LBRACE || sentence.Type == "dict"
}

// 判断该语句是否为列表
func (sentence *Sentence) IsList() bool {
	return sentence.IsBracket() || len(sentence.List) != 0 || sentence.base == LBRACKET || sentence.Type == "list" || utils.Startswith(sentence.Type, "[")
}

// 判断该语句是否为枚举
func (sentence *Sentence) IsEnum() bool {
	return sentence.IsGroup() || sentence.base == LGROUP || sentence.Type == "enum"
}

// 判断该语句是否为必要变量
func (sentence *Sentence) IsRequired() bool {
	return sentence.Value == ""
}

// 判断该语句是否为常量
func (sentence *Sentence) IsConstant() bool {
	return utils.Endswith(sentence.Value, ",constant")
}

// 判断该语句是否为选填变量
func (sentence *Sentence) IsOptional() bool {
	return (sentence.Value == "none") || !sentence.IsRequired() && !sentence.IsConstant()
}

// 查找类型参数序号
func (sentence *Sentence) Find(arg string) int {
	if sentence == nil {
		return -1
	}
	for i, s := range sentence.Args {
		if s == arg {
			return i
		}
	}
	return -1
}

// 添加子语句
func (sentence *Sentence) Add(typ, name, hint, value string, args []string, vt Types, length int64) *Sentence {
	typ, val := AutoType(typ, value)
	s := &Sentence{
		typ, name, hint, value, args,
		-1, sentence.Find(typ), sentence,
		0, nil, make([]*Sentence, 0), make(map[string]*Sentence),
	}

	if length != 0 {
		s.base = LBRACKET
		if length < 0 {
			s.Length = -1
		} else {
			s.Length = length
		}
		s.Add(typ, "", "", "", []string{}, vt, 0)
		s.Type = fmt.Sprintf("[%v]%v", length, s.Type)
	}

	if tk, ok := vt[typ]; ok {
		if tk.IsEnum() {
			val = tk.Map[value].Value
		} else {
			s.base = tk.base
			if tk.IsList() {
				s.Length = tk.Length
				for _, s2 := range tk.List {
					if s2.index >= 0 {
						ty, args := NameSlice(s.Args[s2.index])
						s.Add(ty, "", s2.Hint, s2.Value, args, vt, 0)
					} else {
						s.List = append(s.List, s2)
					}
				}
			} else {
				for s1, s2 := range tk.Map {
					if s2.index >= 0 {
						ty, args := NameSlice(s.Args[s2.index])
						s.Add(ty, s1, s2.Hint, s2.Value, args, vt, 0)
					} else {
						s.Map[s1] = s2
					}
				}
			}
		}
	}
	s.SetOutput(val)

	if sentence != nil {
		if sentence.base == LBRACKET {
			sentence.List = append(sentence.List, s)
		} else {
			sentence.Map[s.Name] = s
		}
	}
	return s
}

// 移除子语句
func (sentence *Sentence) Pop(name string) *Sentence {
	if v, ok := sentence.Map[name]; ok {
		delete(sentence.Map, name)
		return v
	}
	return nil
}

// 转字典
func (sentence *Sentence) ToDict() map[string]string {
	dic := make(map[string]string)
	utils.ForMap(sentence.Map, func(s string, t *Sentence) { dic[s] = t.Value })
	return dic
}

// 修改输出
func (sentence *Sentence) SetOutput(val any) {
	if sentence.IsList() {
		sentence.Output = &sentence.List
	} else if sentence.IsDict() || sentence.IsEnum() {
		sentence.Output = sentence.Map
	} else {
		sentence.Output = val
	}
}
