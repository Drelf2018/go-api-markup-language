package translator

import (
	"fmt"
	"strings"

	parser "github.com/Drelf2018/go-bilibili-api"
)

const Python = `import json
from dataclasses import dataclass, field
from typing import Dict

import httpx


@dataclass
class Api:
    url: str
    method: str
    comment: str = ""
    data: Dict[str, dict] = field(default_factory=dict)
    params: Dict[str, dict] = field(default_factory=dict)

    def __post_init__(self):
        self.method = self.method.upper()
        self.original_data = self.data.copy()
        self.original_params = self.params.copy()
        self.data = {k: v.get("value", "").replace(",constant", "") for k, v in self.data.items()}
        self.params = {k: v.get("value", "").replace(",constant", "") for k, v in self.params.items()}
        self.__result = None

    def request(self):
        return httpx.request(self.method, self.url, data=self.data, params=self.params).text

    def update_data(self, **kwargs):
        self.data.update(kwargs)
        self.__result = None
        return self

    def update_params(self, **kwargs):
        self.params.update(kwargs)
        self.__result = None
        return self

    def update(self, **kwargs):
        if self.method == "GET":
            return self.update_params(**kwargs)
        else:
            return self.update_data(**kwargs)


def get_api(path: str):
    with open(path, "r", encoding="utf-8") as fp:
        return json.load(fp)


API = get_api("%v")
`

func PythonFunc(api *parser.Api) string {
	r, o, all := []string{}, []string{}, []string{}
	hint := api.Hint
	if len(api.Params) != 0 {
		hint += "\n\n    Args:"
	}
	for _, token := range api.Params {
		hint += fmt.Sprintf("\n        %v (%v): %v\n", token.Name, token.Type, token.Hint)
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
	args := append(r, o...)
	return fmt.Sprintf(`def %v(%v):
    """
    %v
    """
    api = Api(**API["%v"])
    api.update(%v)
    return api.request()
`, api.Function, strings.Join(args, ", "), hint, api.Function, strings.Join(all, ", "))
}

func ToPython(am *parser.ApiManager, api string) error {
	s := fmt.Sprintf(Python, api)
	for _, api := range am.Apis {
		s += "\n\n" + PythonFunc(api)
	}
	return parser.WriteFile(strings.ReplaceAll(api, ".json", ".py"), s)
}
