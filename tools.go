package main

import (
	"io"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

// 文件操作

type fileOperaterOptions struct {
	operater    string // read or write
	regexp      string // the regexp to match the string
	replacement string // the replacement of the regexp
	createble   bool   // if the file is not exist, create it
}

// 对于返回值[]string，当operrater为read时，data[0]即为string化后的文本内容；当operater为write时，data[0]为修改前内容， data[1]为修改后内容
// 当operater为write时，fileOperater总会同步地更改文件内容
func fileOperater(url string, options fileOperaterOptions) ([]string, error) {
	file, err1 := os.Open(url)
	if err1 != nil {
		if options.createble {
			_, err2 := os.Create(url)
			if err2 != nil {
				return make([]string, 0), err2
			}
		} else {
			return make([]string, 0), err1
		}
	}
	defer file.Close()
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
		_, err3 := io.WriteString(file, returnText[1])
		if err3 != nil {
			return make([]string, 0), err3
		}
		return returnText, nil
	}
	return make([]string, 0), nil
}

// 数据存储结构

type dataStruct struct {
	// .infor 文件格式： \*index\\key\*key\key\\value\*value\value\\*index\ , 其中*标为变量
	data map[string]string
}

// todo: 应该更改逻辑， /index/标识内外层更好， 而对于map而言显然某键是第几项是不必要的，多层map可以考虑泛型加上结构嵌入实现
func dataLoad(file string) map[string]string {
	d := dataStruct{make(map[string]string)}
	re1 := regexp.MustCompile(`\\[0-9]+\\(.*)\\[0-9]+\\`)
	re2 := regexp.MustCompile(`\\key\\(.*)\\key\\`)
	re3 := regexp.MustCompile(`\\value\\(.*)\\value\\`)
	for index, value := range re1.FindStringSubmatch(file) {
		if index == 0 {
			continue
		}
		d.data[re2.FindStringSubmatch(value)[1]] = re3.FindStringSubmatch(value)[1]
	}
	return d.data
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

// 本地程序操作
// todo ： 当前只做了windows的实现，需要补充linux的实现，并未做判断

func lookPath(file string) string {
	// 查看file的路径

	path, lookPathErr := exec.LookPath(file)
	if lookPathErr != nil {
		panic(lookPathErr)
	}
	return path
}

func start(file string, address string) (string, error) {
	// 启动file，返回进程pid

	isLinux := false

	if isLinux {
		cmd := exec.Command("./" + file, "&")
		cmd.Dir = address
		cmd.Start()
		cmd2 := exec.Command("pgrep", strings.Split(file, ".")[0])
		cmd2.Dir = address
		processId, err := cmd2.Output()
		if err != nil {
			return "-1", err
		}
		return string(processId), nil
	} else {
		cmd := exec.Command(lookPath("cmd"), "/c", file)
		cmd.Dir = address
		cmd.Start()
		cmd2 := exec.Command("pwsh", "-Command", "(Get-Process -name "+strings.Split(file, ".")[0]+").id")
		cmd2.Dir = address
		processId, err := cmd2.Output()
		if err != nil {
			return "-1", err
		}
		return string(processId), nil
	}
}

func kill(processId string) error {
	// 结束进程

	killCmd, killErr := exec.Command("pwsh", "-Command", "Stop-Process "+processId).Output()
	if killErr != nil {
		return killErr
	}
	println(string(killCmd))
	return nil
}
