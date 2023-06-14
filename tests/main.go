package main

import (
	parser "github.com/Drelf2018/go-bilibili-api"
	"github.com/Drelf2018/go-bilibili-api/translator"
	"github.com/Drelf2020/utils"
)

var log = utils.GetLog()

func main() {
	am := parser.GetApi("./tests/user.aml")
	for _, api := range am.Apis {
		log.Debug(api.Info.Tokens["response"])
	}
	translator.ToPython(am, "./tests/", "user")
}
