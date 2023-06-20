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
				s = strings.TrimSpace(s)
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
func In(ls []string, v string) bool {
	return len(utils.Filter(ls, func(t string) bool { return t == v })) != 0
}
