package translator

import (
	"fmt"
	"strings"

	parser "github.com/Drelf2018/go-bilibili-api"
)

func ToPythonFunc(format string, api *parser.Api) string {
	r, o, all := []string{}, []string{}, []string{}
	hint := api.Hint
	params := []string{}
	for _, token := range api.Params {
		params = append(params, fmt.Sprintf("%v (%v): %v", token.Name, token.Type, token.Hint))
		if token.IsConstant() {
			continue
		}
		if token.IsRequired() {
			r = append(r, token.ToPython())
		} else if token.IsOptional() {
			o = append(o, token.ToPython())
		}
		all = append(all, token.Name+"="+token.Name)
	}
	if len(api.Params) != 0 {
		hint += "\n\n    Args:\n        " + strings.Join(params, "\n\n        ")
	}
	args := append(r, o...)

	format = strings.ReplaceAll(format, "demo", api.Function)
	format = strings.ReplaceAll(format, "args", strings.Join(args, ", "))
	format = strings.ReplaceAll(format, "{hint}", hint)
	format = strings.ReplaceAll(format, "update(", "update("+strings.Join(all, ", "))
	return format
}

func ToPython(am *parser.ApiManager, path, name string) error {
	am.ToJson(path + name + ".json")
	s := parser.ReadFile("./template/python/func.py")
	include := s[:strings.Index(s, "# loop")]
	include = strings.Replace(include, "{path}", path+name+".json", 1)
	function := parser.Slice(s, "# loop", "# end", 0)
	for _, api := range am.Apis {
		include += ToPythonFunc(function, api)
	}
	parser.WriteFile(path+"api.py", parser.ReadFile("./template/python/api.py"))
	return parser.WriteFile(path+name+".py", include)
}
