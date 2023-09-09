package main

import (
	"io"
	"os"
	"regexp"
)

type fileOperaterOptions struct {
	operater    string // read or write
	writeRange  string // the range of the write, "both" or "some"
	regexp      string // the regexp to match the string, if the writeRange is "both", this field will be ignored
	replacement string // the replacement of the regexp

}

func fileOperater(url string, options fileOperaterOptions) ([]string, error) { // 第一个返回值，当operrater为read时，data[0]即为string化后的文本内容，当operater为write时，data[0]为修改前内容， data[1]为修改后内容
	file, err1 := os.Open(url + "config.yml")
	if err1 != nil {
		return make([]string, 0), err1
	}
	switch options.operater {
	case "read":
		data, err2 := io.ReadAll(file)
		if err2 != nil {
			return make([]string, 0), err2
		}
		returnText := make([]string, 1)
		returnText[0] = string(data)
		return returnText, nil
	case "write":
		_, err := regexp.Compile(options.regexp)
		if err != nil {
			return make([]string, 0), err
		}
		// returnText := make([]string, 2) // 等待后续操作
		// returnText[0] = string(file)
		return make([]string, 2), nil
	}
	return make([]string, 0), nil
}
