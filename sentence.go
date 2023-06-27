package aml

import (
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
	for i, s := range sentence.Args {
		if s == arg {
			return i
		}
	}
	return -1
}

// 使用类型
func (sentence *Sentence) Use(vt Types) *Sentence {
	log.Debug("type | ", sentence.Type)
	tk := vt[sentence.Type]
	if tk != nil {
		sentence.base = tk.base
		if tk.IsList() {
			//
		} else {
			for s1, s2 := range tk.Map {
				if s2.index >= 0 {
					typ := sentence.Args[s2.index]
					sentence.Map[s1] = s2.Copy(typ).Use(vt)
				} else {
					sentence.Map[s1] = s2
				}
			}
		}
		sentence.SetOutput(nil)
	}
	return sentence
}

// 添加子语句
func (sentence *Sentence) Add(typ, name, hint, value string, args []string, vt Types) *Sentence {
	typ, val := AutoType(typ, value)

	base := -1
	if tk, ok := vt[typ]; GetKind(typ) == IDENTIFIER && ok {
		base = tk.base
		if tk.IsEnum() {
			base = -1
			val = tk.Map[value].Value
		}
	}

	idx := -1
	if sentence != nil {
		idx = sentence.Find(typ)
	}
	s := &Sentence{
		typ, name, hint, value, args,
		base, idx, sentence,
		nil, make([]*Sentence, 0), make(map[string]*Sentence),
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

// 复制
func (s *Sentence) Copy(typ string) *Sentence {
	typ, args := NameSlice(typ)
	sentence := &Sentence{
		typ, s.Name, s.Hint, s.Value, args,
		s.base, s.index, nil,
		nil, make([]*Sentence, 0), make(map[string]*Sentence),
	}
	return sentence
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
