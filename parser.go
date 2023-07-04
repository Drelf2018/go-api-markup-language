package aml

import (
	"fmt"
	"path/filepath"
	"strconv"

	"github.com/Drelf2020/utils"
)

var log = utils.GetLog()

// 语法分析器
type Parser struct {
	// 词法分析器
	*Lexer
	// 文件路径
	dir string
	// 数组长度
	length int64
	// 暂存的 Sentence
	sentence *Sentence
	// api 管理器
	*ApiManager
}

// 匹配导入语句
func (p *Parser) MatchImport() (*Include, error) {
	_, path := p.Done()
	k, v := p.Done()
	if k != IMPORT {
		return nil, fmt.Errorf("导入格式 %v 错误", v)
	}
	_, types := p.Done()
	return NewInclude(p.dir, path, types), nil
}

// 匹配参数
func (p *Parser) MatchArgs() (args []string, kind int) {
	kind, _ = p.Done()
	if kind != LANGLE {
		return
	}
	depth := 1
	args = append(args, "")
	for depth > 0 {
		k, v := p.Done()
		if depth == 1 && k == COMMA {
			args = append(args, "")
		} else {
			if k == LANGLE {
				depth++
			} else if k == RANGLE {
				depth--
			}
			if depth > 0 {
				args[len(args)-1] += v
			}
		}
	}
	return
}

// 匹配定义语句
func (p *Parser) MatchType() error {
	k, name := p.Done()
	if k != IDENTIFIER {
		return fmt.Errorf("%v 不是一个好的变量名", name)
	}
	var base int
	var hint string
	args, kind := p.MatchArgs()
	if len(args) != 0 {
		kind, _ = p.Done()
	}
	if kind == COLON {
		_, hint = p.Done()
		kind, _ = p.Done()
	}
	if kind == ASSIGNMENT {
		base, _ = p.Done()
	}
	p.sentence = &Sentence{
		"type", name, hint, "", args,
		base, -1, nil,
		nil, make([]*Sentence, 0), make(map[string]*Sentence),
	}
	p.sentence.SetOutput(nil)
	p.Types[name] = p.sentence
	return nil
}

// 匹配类型数组
func (p *Parser) MatchLength() (err error) {
	k, v := p.Done()
	if k == NUMBER {
		p.length, err = strconv.ParseInt(v, 10, 64)
		k, v = p.Done()
	}
	if k != RBRACKET {
		return fmt.Errorf("%v 不是合法的中括号", v)
	}
	return
}

// 匹配 Api
func (p *Parser) MatchApi(typ string) error {
	_, name := p.Done()
	var hint, value string
	k, v := p.Done()
	if k != COLON && k != ASSIGNMENT {
		return fmt.Errorf("%v 不是一个好的 Api 格式", v)
	}
	if k == COLON {
		_, hint = p.Done()
		k, _ = p.Done()
	}
	if k == ASSIGNMENT {
		_, value = p.Done()
	}
	var s *Sentence
	p.sentence = s.Add(typ, name, hint, value, []string{}, p.Types)
	return nil
}

// 匹配变量
func (p *Parser) MatchVar(typ string) *Sentence {
	if p.length != -1 {
		log.Debug(p.length)
	}
	var kind int = LANGLE
	var args []string
	var name, hint, value string
	if GetKind(typ) == IDENTIFIER {
		tk := p.Types[typ]
		if tk != nil && len(tk.Args) != 0 {
			args, kind = p.MatchArgs()
		}
	}
	if p.sentence.base == LBRACKET {
		return p.sentence.Add(typ, name, hint, value, args, p.Types)
	}
	if kind == LANGLE {
		kind, name = p.Done()
	}
	if kind == COLON || kind == ASSIGNMENT {
		name = typ
		typ = "auto"
	} else {
		kind, _ = p.Done()
	}
	if kind == COLON {
		_, hint = p.Done()
		kind, _ = p.Done()
	}
	if kind == ASSIGNMENT {
		_, value = p.Done()
		p.Done()
	}
	return p.sentence.Add(typ, name, hint, value, args, p.Types)
}

// 选择匹配
func (p *Parser) Match() {
	switch p.token.Kind {
	case FROM:
		i, err := p.MatchImport()
		utils.PanicErr(err)
		utils.ForMap(
			NewParser(i.path).Parse().Types,
			func(s string, t *Sentence) { p.Types[s] = t },
			func(s string, t *Sentence) bool { return i.Need(s) },
		)
	case TYPE:
		err := p.MatchType()
		utils.PanicErr(err)
	case GET, POST:
		err := p.MatchApi(p.token.Value)
		utils.PanicErr(err)
	case RBRACE, RBRACKET, RGROUP:
		if p.sentence.IsApi() {
			p.Add(p.sentence)
		}
		p.sentence = p.sentence.parent
	case LBRACKET:
		err := p.MatchLength()
		utils.PanicErr(err)
		p.Done()
		fallthrough
	case NUM, STR, BOOL, AUTO, IDENTIFIER:
		s := p.MatchVar(p.token.Value)
		if s.IsOpen() {
			p.sentence = s
		}
		if s.parent.base != LBRACKET {
			return
		}
	}
	p.Shift()
}

// 从文件解析出 Api
func (p *Parser) Parse() *ApiManager {
	for p.Next() {
		p.Match()
	}
	return p.ApiManager
}

func NewParser(path string) *Parser {
	return &Parser{NewLexer(path), filepath.Dir(path), -1, new(Sentence), NewApiManager()}
}
