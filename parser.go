package aml

import (
	"fmt"
	"path/filepath"
	"strconv"

	"github.com/Drelf2020/utils"
)

// 语法分析器
type Parser struct {
	// 词法分析器
	*Lexer
	// 文件路径
	path string
	// 数组长度
	length int64
	// 暂存的 Sentence
	sentence *Sentence
	// 变量类型
	Types map[string]*Sentence
	// Api 字典
	Output map[string]*Api
}

// 判断类型
func (p *Parser) IsType() (*Sentence, bool) {
	switch p.token.Kind {
	case NUM, STR, BOOL, AUTO:
		return nil, true
	case IDENTIFIER:
		s, ok := p.Types[p.token.Value]
		if ok {
			return s, true
		}
		return nil, In(p.sentence.Args, p.token.Value)
	}
	return nil, false
}

// 匹配导入语句
func (p *Parser) MatchImport() (*Include, error) {
	_, path := p.Done()
	k, v := p.Done()
	if k != IMPORT {
		return nil, fmt.Errorf("导入格式 %v 错误", v)
	}
	kind := COMMA
	types := make([]string, 0)
	for kind == COMMA {
		_, typ := p.Done()
		kind, _ = p.Done()
		types = append(types, typ)
	}
	return NewInclude(p.GetDir(), path, types), nil
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
		0, nil, make([]*Sentence, 0), make(map[string]*Sentence),
	}
	p.sentence.SetOutput(nil)
	p.Types[name] = p.sentence
	return nil
}

// 匹配列表
func (p *Parser) MatchList() (err error) {
	if len(p.sentence.List) == 0 {
		k, v := p.Done()

		if k == AUTO {
			return fmt.Errorf("%v 不是一个好的类型", v)
		}
		p.MatchVar(v)
		if p.length <= 0 {
			p.sentence.Length = -1
		} else {
			p.sentence.Length = p.length
		}
		p.length = 0
	} else {
		p.Shift()
	}
	return nil
}

// 匹配类型数组长度
func (p *Parser) MatchLength() (err error) {
	k, v := p.Done()
	if k == NUMBER {
		p.length, err = strconv.ParseInt(v, 10, 64)
		k, v = p.Done()
	}
	if p.length <= 0 {
		p.length = -1
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
	p.sentence = s.Add(typ, name, hint, value, []string{}, p.Types, 0)
	return nil
}

// 匹配变量
func (p *Parser) MatchVar(typ string) *Sentence {
	var kind int
	var args []string
	var name, hint, value string

	if tk, ok := p.IsType(); ok {
		// 检查这个类型 typ 是否需要参数
		if tk != nil && len(tk.Args) != 0 {
			args, _ = p.MatchArgs()
		}
	} else {
		typ, name = "auto", typ
	}

	// 父语句是列表 那就只读取 typ 和 hint
	if p.sentence.base == LBRACKET {
		k, _ := p.Done()
		if k == COLON {
			_, hint = p.Done()
			p.Done()
		}
		return p.sentence.Add(typ, "", hint, "", args, p.Types, 0)
	}

	if name == "" {
		kind, name = p.Done()
		if kind != IDENTIFIER {
			panic(fmt.Errorf("%v 不是一个好的名字", name))
		}
	}

	kind, _ = p.Done()
	if kind == COLON {
		_, hint = p.Done()
		kind, _ = p.Done()
	}
	if kind == ASSIGNMENT {
		_, value = p.Done()
		p.Done()
	}
	if p.length != 0 {
		s := p.sentence.Add(typ, name, hint, value, args, p.Types, p.length)
		p.length = 0
		return s
	}
	return p.sentence.Add(typ, name, hint, value, args, p.Types, 0)
}

// 选择匹配
func (p *Parser) Match() {
	switch p.token.Kind {
	case FROM:
		i, err := p.MatchImport()
		utils.PanicErr(err)
		utils.ForMap(
			NewParser(i.path).Types,
			func(s string, t *Sentence) { p.Types[s] = t },
			func(s string, t *Sentence) bool { return i.Need(s) },
		)
		return
	case TYPE:
		err := p.MatchType()
		utils.PanicErr(err)
	case GET, POST:
		err := p.MatchApi(p.token.Value)
		utils.PanicErr(err)
	case NUMBER:
		p.length, _ = strconv.ParseInt(p.token.Value, 10, 64)
	case RBRACKET:
		err := p.MatchList()
		utils.PanicErr(err)
		p.sentence = p.sentence.parent
		return
	case RBRACE, RGROUP:
		if p.sentence.IsApi() {
			p.Output[p.sentence.Name] = NewApi(p.sentence)
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
		return
	}
	p.Shift()
}

func (p *Parser) GetDir() string {
	return filepath.Dir(p.path)
}

func (p *Parser) NewExt(ext string) string {
	fullname := filepath.Base(p.path)
	suffix := filepath.Ext(fullname)
	return fullname[0:len(fullname)-len(suffix)] + ext
}

func NewParser(path string) *Parser {
	p := Parser{
		NewLexer(path),
		path,
		0,
		new(Sentence),
		make(map[string]*Sentence),
		make(map[string]*Api),
	}
	for p.Next() {
		p.Match()
	}
	return &p
}
