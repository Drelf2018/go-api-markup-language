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
	// 文件路径
	dir string
	// 词法分析器
	lexer *Lexer
	// 暂存的 Sentence
	sentence *Sentence
	// api 管理器
	*ApiManager
}

// 匹配导入语句
func (p *Parser) MatchImport(old ...*Token) (*Include, error) {
	token := p.lexer.Scan(old...)
	path := token.Value
	token = p.lexer.Scan(token)
	if token.Kind == IMPORT {
		types := p.lexer.Scan(token)
		defer p.lexer.Done(types)
		return NewInclude(path, types.Value), nil
	}
	return nil, fmt.Errorf("导入格式 %v 错误", token.Value)
}

// 匹配参数
func (p *Parser) MatchArgs(old ...*Token) (args []string, token *Token) {
	token = p.lexer.Scan(old...)
	if token.Kind != LANGLE {
		return
	}
	depth := 1
	args = append(args, "")
	for depth > 0 {
		token = p.lexer.Scan(token)
		if depth == 1 && token.Kind == COMMA {
			args = append(args, "")
		} else {
			if token.Kind == LANGLE {
				depth++
			} else if token.Kind == RANGLE {
				depth--
			}
			if depth > 0 {
				args[len(args)-1] += token.Value
			}
		}
	}
	return
}

// 匹配定义语句
func (p *Parser) MatchType(old ...*Token) error {
	token := p.lexer.Scan(old...)
	if token.Kind != IDENTIFIER {
		return fmt.Errorf("%v 不是一个好的变量名", token.Value)
	}
	name := token.Value
	var base int
	var hint string
	args, token := p.MatchArgs(token)
	if len(args) != 0 {
		token = p.lexer.Scan(token)
	}
	if token.Kind == COLON {
		token = p.lexer.Scan(token)
		hint = token.Value
		token = p.lexer.Scan(token)
	}
	if token.Kind == ASSIGNMENT {
		token = p.lexer.Scan(token)
		base = token.Kind
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
func (p *Parser) MatchLength(old ...*Token) (length int64, err error) {
	arg := p.lexer.Scan()
	if arg.Kind == NUMBER {
		length, err = strconv.ParseInt(arg.Value, 10, 64)
		arg = p.lexer.Scan()
	}
	if arg.Kind != RBRACKET {
		return -1, fmt.Errorf("%v 不是合法的中括号", arg.Value)
	}
	return
}

// 匹配 Api
func (p *Parser) MatchApi(old ...*Token) (*Sentence, error) {
	typ := old[0].Value
	token := p.lexer.Scan(old...)
	name := token.Value
	var hint, value string
	token = p.lexer.Scan(token)
	if token.Kind != COLON && token.Kind != ASSIGNMENT {
		return nil, fmt.Errorf("%v 不是一个好的 Api 格式", token.Value)
	}
	if token.Kind == COLON {
		token = p.lexer.Scan(token)
		hint = token.Value
		token = p.lexer.Scan(token)
	}
	if token.Kind == ASSIGNMENT {
		token = p.lexer.Scan(token)
		value = token.Value
	}
	p.lexer.Done(token)
	var s *Sentence
	return s.Add(typ, name, hint, value, []string{}, p.Types), nil
}

// 选择匹配
func (p *Parser) Match(t *Token) *Token {
	// var length int64 = -1
	switch t.Kind {
	case FROM:
		i, err := p.MatchImport(t)
		utils.PanicErr(err)
		utils.ForMap(
			NewParser(i.ToApi(p.dir)).Parse().Types,
			func(s string, t *Sentence) { p.Types[s] = t },
			func(s string, t *Sentence) bool { return t != nil && i.Need(s) },
		)
	case TYPE:
		err := p.MatchType()
		utils.PanicErr(err)
	case GET, POST:
		s, err := p.MatchApi(t)
		utils.PanicErr(err)
		p.sentence = s
	case RBRACE, RBRACKET, RGROUP:
		if p.sentence.IsApi() {
			p.Add(p.sentence)
		}
		p.sentence = p.sentence.parent
	case LBRACKET:
		var err error
		// length, err = p.MatchLength()
		_, err = p.MatchLength()
		utils.PanicErr(err)
		t = p.lexer.Scan()
		fallthrough
	case NUM, STR, BOOL, AUTO, IDENTIFIER:
		typ := t.Value
		var args []string
		var token *Token = t
		var name, hint, value string
		if t.Kind == IDENTIFIER {
			tk := p.Types[typ]
			if tk != nil && len(tk.Args) != 0 {
				args, token = p.MatchArgs(t)
			}
		}
		if p.sentence.base == LBRACKET {
			return nil
		}
		token = p.lexer.Scan(token)
		if token.Kind == COLON || token.Kind == ASSIGNMENT {
			name = typ
			typ = "auto"
		} else {
			name = token.Value
			token = p.lexer.Scan(token)
		}
		if token.Kind == COLON {
			token = p.lexer.Scan(token)
			hint = token.Value
			token = p.lexer.Scan(token)
		}
		if token.Kind == ASSIGNMENT {
			token = p.lexer.Scan(token)
			value = token.Value
			token = p.lexer.Scan(token)
		}
		s := p.sentence.Add(typ, name, hint, value, args, p.Types)
		if s.IsOpen() {
			p.sentence = s
		}
		return token
	}
	return nil
}

// 从文件解析出 Api
func (p *Parser) Parse() *ApiManager {
	for t := range p.lexer.ch {
		for t != nil {
			t = p.Match(t)
		}
		// p.lexer.Done(t)
	}
	return p.ApiManager
}

func NewParser(path string) *Parser {
	return &Parser{filepath.Dir(path), NewLexer(path), new(Sentence), NewApiManager()}
}
