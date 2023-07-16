package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	_ "github.com/Drelf2018/aml2py"
	aml "github.com/Drelf2018/go-api-markup-language"
	"github.com/Drelf2020/utils"
)

func init() {
	aml.Plugin{
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
	}.Load()

	aml.Plugin{
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
	}.Load()
}

// 正常运作
func run(path *string) {
	if *path == "" {
		fmt.Print("请输入 AML 文件路径：")
		fmt.Scan(path)
	}
	parser := aml.NewParser(*path)
	os.Mkdir("output", os.ModePerm)
	utils.ForEach(parser.Export(), func(f aml.File) { utils.WriteFile("./output/"+f.Name, f.Content) })
}

// 重复输出
func LoopPrint(msg string, times int) {
	var temp byte
	for times != 0 {
		fmt.Println(msg)
		for {
			fmt.Scanf("%c", &temp)
			if temp == '\n' {
				break
			}
		}
		times--
	}
}

// 引导模式
func guide() {
	fmt.Print(`教程模式，启动！

aml 是一个用来导出 api 信息为可执行代码或文档的工具。

使用:

	aml -path=xxx.aml [arguments]

	其中 -path 代表要解析的 aml 文件路径

参数(arguments):

`)

	for _, p := range aml.GetLoadedPlugin() {
		fmt.Printf(`	-%v
		%v
		作者：%v 版本：%v
		链接：%v

`, p.Cmd, p.Description, p.Author, p.Version, p.Link)
	}

	LoopPrint("好了现在你会了，可以关闭本窗口开始使用了。", 1<<2)
	LoopPrint("不是，怎么还不关呢？", 1)
	LoopPrint("好了现在你会了，可以关闭本窗口开始使用了。", 1<<3)
	LoopPrint("在期待什么？", 1)
	LoopPrint("好了现在你会了，可以关闭本窗口开始使用了。", 1<<4)
	LoopPrint("再不关也不会理你了！", 1)
	LoopPrint("好了现在你会了，可以关闭本窗口开始使用了。", 1<<5)
	LoopPrint("你...就这么坚持么...", 1)
	LoopPrint("好了现在你会了，可以关闭本窗口开始使用了。", 1<<6)
	LoopPrint("你真讨厌！", 1)
	LoopPrint("好了现在你会了，可以关闭本窗口开始使用了。", -1)
}

func main() {
	// 注册并获取文件路径
	path := flag.String("path", "", "")
	flag.Parse()

	// 未选择导出参数则启动教程模式
	if aml.NoFlag() {
		guide()
	} else {
		run(path)
	}
}
