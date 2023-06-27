package aml

import (
	"bufio"
	"io"
	"os"
	"strings"

	"github.com/Drelf2020/utils"
)

// 词法分析器
type Lexer struct {
	file      io.Reader
	bufReader *bufio.Reader

	// 当前字符和读取时错误
	current rune
	err     error

	// 暂存的字符
	storage []rune

	// Token 池和发送通道
	pool *Pool[Token]
	ch   chan *Token
}

// 获取下位字符
func (l *Lexer) Next() bool {
	if l.err == io.EOF {
		return false
	}
	l.current, _, l.err = l.bufReader.ReadRune()
	if l.err != nil && l.err != io.EOF {
		// 如果读到文件尾 该次仍返回 true
		// 之后运行本函数再返回 false
		// 否则抛异常
		panic(l.err)
	}
	return true
}

// 读取 string
func (l *Lexer) Read() string {
	return string(l.current)
}

// 获取暂存字符长度
func (l *Lexer) Length() int {
	return len(l.storage)
}

// 暂存当前字符
func (l *Lexer) Store() {
	l.storage = append(l.storage, l.current)
}

// 清空暂存并以 string 返回
func (l *Lexer) Restore() string {
	r := string(l.storage)
	l.storage = make([]rune, 0)
	return r
}

// 获取暂存第一个字符
func (l *Lexer) First() string {
	if l.Length() == 0 {
		return ""
	}
	return string(l.storage[0])
}

// 判断暂存是否以引号起始
func (l *Lexer) HasQuotation() bool {
	if l.Length() == 0 {
		return false
	}
	return strings.Contains("\"'`", l.First())
}

// 关闭文件
func (l *Lexer) Close() {
	if file, ok := l.file.(*os.File); ok {
		file.Close()
	}
}

// 读取 Token
func (l *Lexer) Scan(tokens ...*Token) *Token {
	for _, t := range tokens {
		l.Done(t)
	}
	return <-l.ch
}

// 发送
func (l *Lexer) Send(kind int, value string) {
	l.ch <- l.pool.Get().Set(kind, value)
}

// 销毁
func (l *Lexer) Done(t *Token) {
	l.pool.Put(t)
}

// 发送暂存的字符
func (l *Lexer) SendStorage() {
	if l.Length() != 0 {
		result := l.Restore()
		if utils.IsNumber(result) {
			l.Send(NUMBER, result)
		} else {
			l.Send(GetKind(result), result)
		}
	}
}

// 发送当前
func (l *Lexer) SendNow(kind int) {
	l.SendStorage()
	l.Send(kind, l.Read())
}

// 异步启动解析 Token
func (l Lexer) init() *Lexer {
	go func() {
		for l.Next() {
			s := l.Read()

			// 多行文本
			if l.HasQuotation() {
				if l.First() == s {
					// 多行字符串结束了
					l.Send(STRING, l.Restore()[1:])
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
				l.SendStorage()
				continue
			}

			switch s {
			case ",":
				l.SendNow(COMMA)
			case "<":
				l.SendNow(LANGLE)
			case ">":
				l.SendNow(RANGLE)
			case "(":
				l.SendNow(LGROUP)
			case ")":
				l.SendNow(RGROUP)
			case "[":
				l.SendNow(LBRACKET)
			case "]":
				l.SendNow(RBRACKET)
			case "{":
				l.SendNow(LBRACE)
			case "}":
				l.SendNow(RBRACE)
			case "=":
				l.SendNow(ASSIGNMENT)
			case ":":
				l.SendNow(COLON)
			case "#":
				l.SendStorage()
				fallthrough
			default:
				l.Store()
			}
		}

		// 读取至文件尾 关闭文件和通道
		close(l.ch)
		l.Close()
	}()
	return &l
}

// 获取文件流
func FromFile(path string) (l *Lexer) {
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	// bufio.NewReader(file) 等价 bufio.NewReaderSize(file, 4096) 可根据需求修改 size
	return Lexer{
		file,
		bufio.NewReader(file),
		0,
		nil,
		make([]rune, 0),
		NewPool[Token](),
		make(chan *Token),
	}.init()
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
