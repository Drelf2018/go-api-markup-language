package parser

import (
	"os"
	"regexp"

	"github.com/Drelf2020/utils"
)

// 读取文件
func ReadFile(path string) string {
	data, err := os.ReadFile(path)
	if utils.CheckErr(err) {
		return ""
	}
	return string(data)
}

// 写入文件
func WriteFile(path, s string) error {
	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY, 0666)
	if utils.LogErr(err) {
		return err
	}

	_, err = file.WriteString(s)
	if utils.LogErr(err) {
		return err
	}
	return nil
}

var re = regexp.MustCompile(` *(?:(` + TokenTypes.Join() + "|" + RequestTypes.Join() + `) )? *([^:^=^\r^\n^ ]+)(?:: *([^=^\r^\n]+))? *(?:= *([^\r^\n]+))?`)

// 找出所有语句
func FindTokens(api string, callback func(*Token)) {
	for _, s := range re.FindAllStringSubmatch(api, -1) {
		callback(NewToken(s[1:]))
	}
}

// 更新字典
func Update(dic ...map[string]any) map[string]any {
	d0 := dic[0]
	for _, d1 := range dic[1:] {
		for k, v := range d1 {
			d0[k] = v
		}
	}
	return d0
}
