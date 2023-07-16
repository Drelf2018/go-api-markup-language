package aml

import "flag"

type File struct {
	Name    string
	Content string
}

type Plugin struct {
	Cmd         string
	Author      string
	Version     string
	Description string
	Link        string
	Generate    func(p *Parser) (files []File)

	need *bool
}

var plugins = make([]Plugin, 0)

// 注册插件
func (p Plugin) Load() {
	p.need = flag.Bool(p.Cmd, false, "导出 "+p.Cmd+" 的参数")
	plugins = append(plugins, p)
}

// 获取已注册插件
func GetLoadedPlugin() []Plugin {
	r := make([]Plugin, len(plugins))
	copy(r, plugins)
	return r
}

// 未输入参数
func NoFlag() bool {
	for _, plugin := range plugins {
		if *plugin.need {
			return false
		}
	}
	return true
}

// 导出文件
//
// cmds: 可选的必须导出参数
func (p *Parser) Export(cmds ...string) (files []File) {
	for _, plugin := range plugins {
		if *plugin.need || In(cmds, plugin.Cmd) {
			files = append(files, plugin.Generate(p)...)
		}
	}
	return
}
