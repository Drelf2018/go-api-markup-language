package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/Drelf2018/aml2py"
	aml "github.com/Drelf2018/go-api-markup-language"
	"github.com/Drelf2020/utils"
)

var log = utils.GetLog()

func GetPrefix(p string) string {
	fullname := filepath.Base(p)
	suffix := filepath.Ext(fullname)
	return fullname[0 : len(fullname)-len(suffix)]
}

func main() {
	var path string
	if len(os.Args) < 2 {
		log.Info("请输入 AML 文件路径：")
		fmt.Scan(&path)
	} else {
		path = os.Args[1]
	}
	name := GetPrefix(path)
	am := aml.NewParser(path).Parse()
	os.Mkdir("output", os.ModePerm)
	am.ToJson("./output/" + name + ".json")
	am.ToYaml("./output/" + name + ".yml")
	api, file := aml2py.ToPython(am, name+".json")
	utils.WriteFile("./output/api.py", api)
	utils.WriteFile("./output/"+name+".py", file)
}
