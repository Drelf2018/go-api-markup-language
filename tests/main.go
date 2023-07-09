package main

import (
	aml "github.com/Drelf2018/go-api-markup-language"
	"github.com/Drelf2018/go-api-markup-language/translator"
	"github.com/Drelf2020/utils"
)

var log = utils.GetLog()

func main() {
	parser := aml.NewParser("./tests/user.aml")
	am := parser.Parse()
	utils.ForMap(am.Output, func(s string, a *aml.Api) { log.Info(s, " | ", *a) })
	translator.ToPython(am, "./tests/", "user")
}
