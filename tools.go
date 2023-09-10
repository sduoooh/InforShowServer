package main

import (
	"io"
	"os"
	"regexp"
	"time"
	"strconv"
)

type fileOperaterOptions struct {
	operater    string // read or write
	regexp      string // the regexp to match the string
	replacement string // the replacement of the regexp
	createble   bool   // if the file is not exist, create it
}


// 对于返回值[]string，当operrater为read时，data[0]即为string化后的文本内容；当operater为write时，data[0]为修改前内容， data[1]为修改后内容
// 当operater为write时，fileOperater总会同步地更改文件内容
func fileOperater(url string, options fileOperaterOptions) ([]string, error) { 
	file, err1 := os.Open(url + "config.yml")
	defer file.Close()
	if err1 != nil {
		if options.createble {
			_, err2 := os.Create(url + "config.yml")
			if err2 != nil {
				return make([]string, 0), err2
			}
		}else {
			return make([]string, 0), err1
		}
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
		re, err := regexp.Compile(options.regexp)
		if err != nil {
			return make([]string, 0), err
		}
		returnText := make([]string, 2) // 等待后续操作
		data, err2 := io.ReadAll(file)
		if err2 != nil {
			return make([]string, 0), err2
		}
		returnText[0] = string(data)
		returnText[1] = re.ReplaceAllString(string(data), options.replacement)
		_ , err3 := io.WriteString(file, returnText[1])
		if err3 != nil {
			return make([]string, 0), err3
		}
		return returnText, nil
	}
	return make([]string, 0), nil
}

func returnLog(content string) string {
	return time.Now().String() + ": " + content + "\n"
}

type dataStruct struct {
	// .infor 文件格式： \*index\\key\*key\key\\value\*value\value\\*index\ , 其中*标为变量
	data map[string]string
}

func (d *dataStruct) load(file string) error {
	re1 := regexp.MustCompile(`\\[0-9]+\\(.*)\\[0-9]+\\`)
	re2 := regexp.MustCompile(`\\key\\(.*)\\key\\`)
	re3 := regexp.MustCompile(`\\value\\(.*)\\value\\`)
	for index, value := range re1.FindStringSubmatch(file){
		if index == 0 {
			continue
		}
		d.data[re2.FindStringSubmatch(value)[1]] = re3.FindStringSubmatch(value)[1]
	}
	return nil
}

func (d *dataStruct) push() (string, error) {
	result := ""
	index := 0
	for key, value := range d.data {
		index += 1
		result += "\\" + strconv.Itoa(index) + "\\" + "\\key\\" + key + "\\key\\" + "\\value\\" + value + "\\value\\" + "\\" + strconv.Itoa(index) + "\\"
	}
	return result, nil
}