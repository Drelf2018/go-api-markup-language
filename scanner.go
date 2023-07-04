package aml

import (
	"bufio"
	"io"
	"os"
	"strings"
)

type Scanner struct {
	file      io.Reader
	bufReader *bufio.Reader

	// 当前字符和读取时错误
	current rune
	err     error

	// 暂存的字符
	storage []rune
}

// 获取下位字符
func (s *Scanner) Next() bool {
	if s.err == io.EOF {
		return false
	}
	s.current, _, s.err = s.bufReader.ReadRune()
	if s.err != nil && s.err != io.EOF {
		// 如果读到文件尾 该次仍返回 true
		// 之后运行本函数再返回 false
		// 否则抛异常
		panic(s.err)
	}
	return true
}

// 读取 string
func (s *Scanner) Read() string {
	return string(s.current)
}

// 获取暂存字符长度
func (s *Scanner) Length() int {
	return len(s.storage)
}

// 暂存当前字符
func (s *Scanner) Store() {
	s.storage = append(s.storage, s.current)
}

// 清空暂存并以 string 返回
func (s *Scanner) Restore() string {
	r := string(s.storage)
	s.storage = make([]rune, 0)
	return r
}

// 获取暂存第一个字符
func (s *Scanner) First() string {
	if s.Length() == 0 {
		return ""
	}
	return string(s.storage[0])
}

// 判断暂存是否以引号起始
func (s *Scanner) HasQuotation() bool {
	if s.Length() == 0 {
		return false
	}
	return strings.Contains("\"'`", s.First())
}

// 关闭文件
func (s *Scanner) Close() {
	if file, ok := s.file.(*os.File); ok {
		file.Close()
	}
}
