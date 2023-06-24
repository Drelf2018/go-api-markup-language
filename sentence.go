package parser

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/Drelf2020/utils"
)

// 单条语句
//
// 其中 Args 表示其 Name 中携带的参数
//
// 当该语句为字典时 Map 不为空 表示其下包含的语句
type Sentence struct {
	Adorn  string               `json:"adorn,omitempty" yaml:"adron,omitempty"`
	Type   string               `json:"type,omitempty" yaml:"type,omitempty"`
	Name   string               `json:"-" yaml:"-"`
	Hint   string               `json:"hint,omitempty" yaml:"hint,omitempty"`
	Value  string               `json:"-" yaml:"-"`
	Output any                  `json:"value,omitempty" yaml:"value,omitempty"`
	Parent *Sentence            `json:"-" yaml:"-"`
	Args   []string             `json:"-" yaml:"-"`
	Params []string             `json:"-" yaml:"-"`
	List   []*Sentence          `json:"-" yaml:"-"`
	Map    map[string]*Sentence `json:"-" yaml:"-"`
}

// 判断该语句是否为 Api 起始语句
func (sentence *Sentence) IsApi() bool {
	return MethodTypes.Has(sentence.Type)
}

// 判断该语句是否为 type 起始语句
func (sentence *Sentence) IsType() bool {
	return sentence.Type == "type"
}

// 判断该语句是否为枚举
func (sentence *Sentence) IsEnum() bool {
	return sentence.Type == "enum"
}

// 判断该语句是否为起始括号
func (sentence *Sentence) IsOpen() bool {
	return sentence.Value == "{"
}

// 判断该语句是否为起始列表
func (sentence *Sentence) IsBracket() bool {
	return sentence.Value == "["
}

// 判断该语句是否为闭合括号
func (sentence *Sentence) IsClose() bool {
	return sentence.Name == "}" || sentence.Name == "]"
}

// 判断该语句是否为字典
func (sentence *Sentence) IsDict() bool {
	return sentence.IsOpen() || len(sentence.Map) != 0 || sentence.Type == "dict"
}

// 判断该语句是否为列表
func (sentence *Sentence) IsList() bool {
	return sentence.IsBracket() || len(sentence.List) != 0 || sentence.Type == "list" || utils.Startswith(sentence.Type, "[")
}

// 判断是否为多行文本
func (sentence *Sentence) IsMultiLine() string {
	if len(sentence.Value) < 1 {
		return ""
	}
	if chn := sentence.Value[:1]; In([]string{"\"", "'", "`"}, chn) {
		return chn
	}
	return ""
}


func (sentence *Sentence) IsDeprecate() bool {
	return sentence.Adorn == "deprecate"
}

// 判断是否有引号
func (sentence *Sentence) HasQuotation(chn string) bool {
	return strings.Count(sentence.Name, chn) == 1
}

// 判断该语句是否为必要变量
func (sentence *Sentence) IsRequired() bool {
	return sentence.Adorn == "required"
}

// 判断该语句是否为常量
func (sentence *Sentence) IsConstant() bool {
	return utils.Endswith(sentence.Value, ",constant")
}

// 判断该语句是否为选填变量
func (sentence *Sentence) IsOptional() bool {
	return (sentence.Value == "none") || !sentence.IsRequired() && !sentence.IsConstant()
}

// 添加子语句
func (sentence *Sentence) Add(t *Sentence, isList ...bool) {
	if sentence == nil {
		return
	}
	t.Parent = sentence
	if sentence.IsList() || (len(isList) != 0 && isList[0]) {
		sentence.List = append(sentence.List, t)
	} else {
		sentence.Map[t.Name] = t
	}
}

// 获取子语句
func (sentence *Sentence) Get(name string) *Sentence {
	return sentence.Map[name]
}

// 移除子语句
func (sentence *Sentence) Pop(name string) *Sentence {
	if v := sentence.Get(name); v != nil {
		delete(sentence.Map, name)
		return v
	}
	return nil
}

// 交换类型和变量名
func (sentence *Sentence) Exchange(vt *Types) *Sentence {
	return sentence.Copy(sentence.Name+"<"+strings.Join(sentence.Args, ",")+">", vt)
}

// 转字典
func (sentence *Sentence) ToDict() map[string]string {
	dic := make(map[string]string)
	utils.ForMap(sentence.Map, func(s string, t *Sentence) { dic[s] = t.Value })
	return dic
}

// 复制
func (sentence *Sentence) Copy(nt string, vt *Types) *Sentence {
	t := NewSentence(nt, sentence.Name, sentence.Hint, sentence.Value)
	t.SetTypes(vt)
	return t
}

// 修改输出
func (sentence *Sentence) SetOutput() {
	if sentence.IsList() {
		sentence.Output = &sentence.List
	} else if sentence.IsDict() {
		sentence.Output = sentence.Map
	}
}

// 返回类型数组的具体内容和长度
func (sentence *Sentence) GetLength(value ...string) (string, int64) {
	var typ string = sentence.Type
	var length int64 = -1

	if len(value) != 0 {
		typ = value[0]
	}

	// 判断是否为数组
	utils.ForEach(
		regexp.MustCompile(`\[(\d*)\]`).FindAllStringSubmatch(typ, -1),
		func(s []string) {
			typ = strings.Replace(typ, s[0], "", 1)
			length, _ = strconv.ParseInt(s[1], 10, 64)
			if length < 1 {
				length = 1
			}
		},
		func(s []string) bool { return len(s) != 0 },
	)

	return typ, length
}

// 修改类型
func (sentence *Sentence) SetTypes(vt *Types) {
	typ, length := sentence.GetLength()
	tk := vt.Get(typ)

	// 基础类型 不做修改 输出原值
	if tk == nil {
		if length != -1 {
			tk = NewSentence(typ, "", "", "")
			for i := 0; i < int(length); i++ {
				sentence.Add(tk, true)
			}
		}
		// 有子语句(dict | list) 替换输出为子语句集合
		sentence.SetOutput()
		return
	}

	// 枚举 替换输出为具体值
	if tk.IsEnum() {
		sentence.Output = tk.Get(sentence.Value).Value
		return
	}

	// 自定义类型
	var nt *Sentence
	argsMap := NewZip(tk.Args, sentence.Params, sentence.Type+" 的参数个数都能数歪来？")

	Copy := func(t *Sentence) *Sentence {
		p := ""
		if len(t.Params) != 0 {
			p = "<" + strings.Join(t.Params, ",") + ">"
		}
		return t.Copy(argsMap.Same(t.Type)+p, vt)
	}

	if length != -1 {
		nt = NewSentence(tk.Name, tk.Name, tk.Hint, tk.Value)
	} else {
		nt = sentence
	}
	if tk.IsDict() {
		utils.ForMap(tk.Map, func(s string, t *Sentence) { nt.Add(Copy(t)) })
	} else if tk.IsList() {
		utils.ForEach(tk.List, func(t *Sentence) { nt.Add(Copy(t), true) })
	}
	if length != -1 {
		nt.SetOutput()
		for i := 0; i < int(length); i++ {
			sentence.Add(nt, true)
		}
	}

	sentence.SetOutput()
}

// 解析类似 res<T1,T2> 中的子类型 T1 T2
//
// data 顺序 Sentence Type Name Hint Value
func NewSentence(data ...string) *Sentence {
	// 自动推断类型
	typ, output := AutoType(data[0], data[3])
	// 解析类型中的参数
	typ, params := NameSlice(typ)
	// 解析变量名中的参数
	name, args := NameSlice(strings.TrimSpace(data[1]))
	// 去除标注前后空白
	hint := strings.TrimSpace(data[2])

	return &Sentence{"optional", typ, name, hint, data[3], output, nil, args, params, make([]*Sentence, 0), make(map[string]*Sentence)}
}
