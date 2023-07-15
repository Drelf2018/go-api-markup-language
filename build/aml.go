package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	// _ "github.com/Drelf2018/aml2py"
	aml "github.com/Drelf2018/go-api-markup-language"
	"github.com/Drelf2020/utils"
)

var log = utils.GetLog()

func init() {
	// 导出 json 插件
	aml.Load(aml.Plugin{
		Cmd:         "json",
		Author:      "Drelf2018",
		Version:     "0.0.1",
		Description: "导出 json 文件",
		Link:        "https://github.com/Drelf2018/go-api-markup-language/build",
		Generate: func(p *aml.Parser) []aml.File {
			output := aml.JsonDump(p.Output, "    ")
			utils.ForMap(
				p.Output,
				func(s string, a *aml.Api) {
					info := aml.JsonDump(a.Info.ToDict(), "        ")
					output = strings.Replace(output, "\"function\": \""+s+"\"", utils.Slice(info, "\"", "\"", 3), 1)
				},
				func(s string, a *aml.Api) bool { return a.Function != "" },
			)
			return []aml.File{{
				Name:    p.NewExt(".json"),
				Content: output,
			}}
		},
	})

	// 导出 yaml 插件
	aml.Load(aml.Plugin{
		Cmd:         "yaml",
		Author:      "Drelf2018",
		Version:     "0.0.1",
		Description: "导出 yml 文件",
		Link:        "https://github.com/Drelf2018/go-api-markup-language/build",
		Generate: func(p *aml.Parser) []aml.File {
			output := aml.YamlDump(p.Output)
			utils.ForMap(
				p.Output,
				func(s string, a *aml.Api) {
					info := aml.YamlDump(map[string]map[string]string{"info": a.Info.ToDict()})
					output = strings.Replace(output, "  function: "+s+"\n", strings.Replace(info, "info:\n", "", 1), 1)
				},
				func(s string, a *aml.Api) bool { return a.Function != "" },
			)
			return []aml.File{{
				Name:    p.NewExt(".yml"),
				Content: output,
			}}
		},
	})
}

func main() {
	cmds := make(map[string]*bool)
	utils.ForEach[string](
		aml.GetCMD(),
		func(s string) { cmds[s] = flag.Bool(s, false, "") },
	)
	path := flag.String("path", "", "")
	flag.Parse()
	keys := []string{}
	for k, ok := range cmds {
		if *ok {
			keys = append(keys, k)
		}
	}
	if *path == "" {
		log.Info("请输入 AML 文件路径：")
		fmt.Scan(path)
	}
	parser := aml.NewParser(*path)
	os.Mkdir("output", os.ModePerm)
	utils.ForEach(
		parser.Export(keys...),
		func(f aml.File) {
			utils.WriteFile("./output/"+f.Name, f.Content)
		},
	)
}
