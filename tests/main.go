package main

import (
	parser "github.com/Drelf2018/go-bilibili-api"
	"github.com/Drelf2020/utils"
)

var log = utils.GetLog()

func main() {
	// am := parser.GetApi("./tests/user.aml")
	// utils.ForMap(am.Output, func(s string, a *parser.Api) { log.Info(s, " | ", *a) })
	// translator.ToPython(am, "./tests/", "user")
	ch := make(chan parser.Token, 16)
	go parser.Parse("./tests/user.aml", ch)
	for t := range ch {
		// log.Debug(t, " | ", []byte(t.Value))
		log.Debug(t.Value)
	}
}
