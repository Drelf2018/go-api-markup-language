package aml

import (
	"bufio"
	"os"
	"strings"

	"github.com/Drelf2020/utils"
)

// 词法分析器
type Lexer struct {
	*Scanner
	token  *Token
	victim *Token
}

func (l *Lexer) Next() bool {
	if l.token.NotNull() {
		return true
	}
	l.Read()
	return l.token.NotNull()
}

func (l *Lexer) Shift() *Lexer {
	l.token.Shift(l.victim)
	return l
}

func (l *Lexer) Done() (int, string) {
	l.Shift().Next()
	return l.token.Kind, l.token.Value
}

// 保存暂存的字符
//
// 返回是否保存成功
func (l *Lexer) SaveStorage() bool {
	if l.Length() != 0 {
		result := l.Restore()
		if utils.IsNumber(result) {
			l.token.New(NUMBER, result)
		} else {
			l.token.New(GetKind(result), result)
		}
		return true
	}
	return false
}

// 保存暂存和当前的字符
func (l *Lexer) SaveNow(kind int) {
	if l.SaveStorage() {
		l.victim.New(kind, l.Scanner.Read())
	} else {
		l.token.New(kind, l.Scanner.Read())
	}
}

func (l *Lexer) Read() {
	for l.Scanner.Next() {
		s := l.Scanner.Read()

		// 多行文本
		if l.HasQuotation() {
			if l.First() == s {
				// 多行字符串结束了
				l.token.New(STRING, l.Restore()[1:])
				return
			} else {
				l.Store()
			}
			continue
		}

		// 如果开头是 # 则忽略直到换行
		if l.First() == "#" {
			if s != "\n" {
				continue
			}
			l.Restore()
		}

		// 跳过空白字符
		if strings.Contains(" \t\n\r", s) {
			if l.SaveStorage() {
				return
			}
			continue
		}

		switch s {
		case ",":
			l.SaveNow(COMMA)
		case "<":
			l.SaveNow(LANGLE)
		case ">":
			l.SaveNow(RANGLE)
		case "(":
			l.SaveNow(LGROUP)
		case ")":
			l.SaveNow(RGROUP)
		case "[":
			l.SaveNow(LBRACKET)
		case "]":
			l.SaveNow(RBRACKET)
		case "{":
			l.SaveNow(LBRACE)
		case "}":
			l.SaveNow(RBRACE)
		case "=":
			l.SaveNow(ASSIGNMENT)
		case ":":
			l.SaveNow(COLON)
		case "#":
			l.SaveStorage()
			fallthrough
		default:
			l.Store()
			continue
		}

		// 井号开头继续循环 其他直接退出
		return
	}
}

// 获取文件流
func FromFile(path string) (l *Lexer) {
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	// bufio.NewReader(file) 等价 bufio.NewReaderSize(file, 4096) 可根据需求修改 size
	return &Lexer{
		&Scanner{
			file,
			bufio.NewReader(file),
			0,
			nil,
			make([]rune, 0),
		},
		new(Token),
		new(Token),
	}
}

// 获取网络流
func FromURL(url string) (l *Lexer) {
	return new(Lexer)
}

// 自动选择
func NewLexer(path string) *Lexer {
	if utils.Startswith(path, "http") {
		return FromURL(path)
	}
	return FromFile(path)
}
