package parser

import (
	"strings"

	"github.com/Drelf2020/utils"
)

var apiText = new(ApiText)

type ApiText struct {
	Text  string
	Lines []string
}

func (at *ApiText) Find(s string) int {
	for i, line := range at.Lines {
		if s == line {
			return i
		}
	}
	return -1
}

func (at *ApiText) Accumulate(o, s string) string {
	st := at.Find(o)
	if st == -1 {
		return ""
	}
	if strings.Count(o, s) != 1 {
		return utils.Slice(o, s, s, 0)
	}
	for i := st + 1; i < len(at.Lines); i++ {
		if strings.Count(at.Lines[i], s) == 1 {
			return utils.Slice(strings.Join(at.Lines[st:i+1], ""), s, s, 0)
		}
	}
	return o[1:]
}

func NewText(api string) *ApiText {
	return &ApiText{
		api,
		strings.Split(strings.ReplaceAll(api, "\r\n", "\n"), "\n"),
	}
}
