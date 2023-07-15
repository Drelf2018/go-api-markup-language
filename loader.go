package aml

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
}

var plugins = make(map[string]Plugin)

// 注册插件
func Load(p Plugin) {
	plugins[p.Cmd] = p
}

// 允许的命令
func GetCMD() (r []string) {
	for c := range plugins {
		r = append(r, c)
	}
	return
}

// 导出
func (p *Parser) Export(cmds ...string) (files []File) {
	for _, cmd := range cmds {
		if plugin, ok := plugins[cmd]; ok {
			files = append(files, plugin.Generate(p)...)
		} else {
			log.Errorf("参数 %v 不支持", cmd)
		}
	}
	return
}
