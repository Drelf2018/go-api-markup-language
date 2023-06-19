package parser

import (
	"encoding/json"
	"strings"

	"github.com/Drelf2020/utils"
	"gopkg.in/yaml.v2"
)

// 纯净类型
func NameSlice(s string) (name string, args []string) {
	name = strings.Split(s, "<")[0]
	if text := utils.Slice(s, "<", ">", 0); text != "" {
		depth := 0
		utils.ForEach(
			strings.Split(text, ","),
			func(s string) {
				if depth == 0 {
					args = append(args, s)
				} else {
					args[len(args)-1] += "," + s
				}
				depth += strings.Count(s, "<") - strings.Count(s, ">")
			},
		)
	}
	return
}

// json 序列化
func JsonDump(v any, indent string) string {
	b, err := json.MarshalIndent(v, "", indent)
	utils.PanicErr(err)
	return string(b)
}

// yaml 序列化
func YamlDump(v any) string {
	b, err := yaml.Marshal(v)
	utils.PanicErr(err)
	return string(b)
}

// python list in
func In[T string](ls []T, v T) bool {
	return len(utils.Filter[T](ls, func(t T) bool { return t == v })) != 0
}

type Dict map[string]string

func (dict Dict) Get(key, value string) string {
	v, ok := dict[key]
	if ok {
		return v
	}
	return value
}

func (dict Dict) Same(key string) string {
	return dict.Get(key, key)
}

func NewDict(l1, l2 []string, errMsg string) (dict Dict) {
	if len(l1) != len(l2) {
		panic(errMsg)
	}
	dict = make(Dict)
	for i, arg := range l1 {
		dict[arg] = l2[i]
	}
	return
}
