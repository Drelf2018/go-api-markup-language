package translator

import (
	"fmt"
	"strings"

	parser "github.com/Drelf2018/go-bilibili-api"
	"github.com/Drelf2020/utils"
)

// 转 value 为 python 格式
func ValueToPython(typ, val string) string {
	val = strings.Replace(val, ",constant", "", 1)
	if typ == "str" {
		if val == "none" {
			return "\"\""
		}
		return "\"" + val + "\""
	} else if typ == "bool" {
		if val == "none" {
			return "False"
		}
		return utils.Capitalize(val)
	} else {
		if val == "none" {
			return "0"
		}
		return val
	}
}

// 转为 Python 格式
func TokenToPython(token *parser.Token) (s string) {
	s = token.Name
	typ := strings.Replace(token.Type, "num", "int", 1)
	val := token.Value
	if typ != "" {
		s += ": " + typ
	}

	if val != "" {
		s += " = " + ValueToPython(typ, val)
	}
	return
}

func ToPythonFunc(format string, api *parser.Api) string {
	r, o, all := []string{}, []string{}, []string{}
	hint := api.Hint
	params := []string{}
	for _, token := range api.Params.Tokens {
		params = append(params, fmt.Sprintf("%v (%v): %v", token.Name, token.Type, token.Hint))
		if token.IsConstant() {
			typ := strings.Replace(token.Type, "num", "int", 1)
			all = append(all, token.Name+"="+ValueToPython(typ, token.Value))
			continue
		}
		if token.IsRequired() {
			r = append(r, TokenToPython(token))
		} else if token.IsOptional() {
			o = append(o, TokenToPython(token))
		}
		all = append(all, token.Name+"="+token.Name)
	}
	if len(api.Params.Tokens) != 0 {
		hint += "\n\n    Args:\n        " + strings.Join(params, "\n\n        ")
	}
	args := append(r, o...)

	format = strings.ReplaceAll(format, "demo", api.Function)
	format = strings.ReplaceAll(format, "args", strings.Join(args, ", "))
	format = strings.ReplaceAll(format, "hint", hint)
	format = strings.ReplaceAll(format, "update(", "update("+strings.Join(all, ", "))
	return format
}

func ToPython(am *parser.ApiManager, path, name string) error {
	am.ToJson(path + name + ".json")
	am.ToYaml(path + name + ".yml")
	s := utils.ReadFile("./template/python/func.py")
	include := s[:strings.Index(s, "# loop")]
	include = strings.Replace(include, "path", path+name+".json", 1)
	function := utils.Slice(s, "# loop", "# end", 0)
	for name, api := range am.Output {
		api.Function = name
		include += ToPythonFunc(function, api)
	}
	utils.WriteFile(path+"api.py", utils.ReadFile("./template/python/api.py"))
	return utils.WriteFile(path+name+".py", include)
}
