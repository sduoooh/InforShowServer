package main

import (
	"bytes"
	"errors"
	"os/exec"
	"regexp"
	"strconv"
)

// 注意，为方便传递，qq中的id都是int类型，但是本程序内部传递的id为string.
// 尚待与abstract协同完善

type qqContentAttr struct {
	isGroup     bool
	isPrivate   bool
	isAnonymity bool
	isRichText  bool
}

type qqHandlerInfor struct {

	// nil if isGroup is false
	groupId   int
	groupName string

	// always not nil
	userId   int
	userName string
}

type qqContentSet struct {
	text string

	// nil if isRichText is false
	richText []string // if some rich text exists, such as pics or videos, it will be stored in this slice
}

type qqUserIdInfor struct {
	userId int
}

func readPort(url string) (string, error) {
	re, _ := regexp.Compile(url + "config.yml")
	data, err := fileOperater(url+"config.yml", fileOperaterOptions{operater: "read"})
	if err != nil {
		return "-1", err
	}
	return re.FindStringSubmatch(data[0])[1], nil
}

func creater(sourceAddress string, occupiedPort *map[string]portInfor) func(string) (*taskInfor[qqContentAttr, qqHandlerInfor, qqContentSet, qqUserIdInfor], error) {
	return func(entrance string) (*taskInfor[qqContentAttr, qqHandlerInfor, qqContentSet, qqUserIdInfor], error) {
		// 尚待施工，需要先确定特定环境文件夹内容构造；
		// 目前认为还需使在其中读取一些信息，如预期使用的端口号及qq号等，后期或许在此预读数据库

		task := taskInfor[qqContentAttr, qqHandlerInfor, qqContentSet, qqUserIdInfor]{}
		portInfor, err := fileOperater(sourceAddress+entrance+"/qqPortList", fileOperaterOptions{operater: "read", createble: false})
		if err != nil {
			return nil, errors.New("can't read qqPortList")
		}
		task.occupiedPort = make([]string, 1)
		task.occupiedPort[0] = portInfor[0]
		userId, _ := strconv.Atoi(entrance)
		task.userIdInfor = qqUserIdInfor{userId}
		// task.processId = entrance
		task.execution = func () error {
		exec.Command()
		task.uploadChannel = make(chan transInfor[qqContentAttr, qqHandlerInfor, qqContentSet], 1)
		task.downloadChannel = make(chan transInfor[qqContentAttr, qqHandlerInfor, qqContentSet], 1)

		cmd := exec.Command(sourceAddress + entrance + "/go-cqhttp_windows_amd64.exe")
		var stdout, stderr bytes.Buffer
		cmd.Stdout = &stdout 
		cmd.Stderr = &stderr 
		err1 := cmd.Run()
		outStr, errStr := string(stdout.Bytes()), string(stderr.Bytes())
		if err1 != nil {
			return nil, err1
		}
		fileOperater("../qq/qqLog.txt", fileOperaterOptions{operater: "write", regexp: ".*", replacement: returnLog(outStr), createble: true})
		fileOperater("../qq/qqLog.txt", fileOperaterOptions{operater: "write", regexp: ".*", replacement: returnLog(errStr), createble: true})
		return nil, nil
	}
}

func qqServerInit() error {
	qqMap, err := fileOperater("../qq/qqMap.txt", fileOperaterOptions{operater: "read", createble: false})
	if err != nil {
		return errors.New("can't read qqMap.txt")
	}
	temp := dataStruct{make(map[string]string)}
	temp.load(qqMap[0])
	entranceList := temp.data
	qqTaskMaster := taskMasterInfor[qqContentAttr, qqHandlerInfor, qqContentSet, qqUserIdInfor]{}
	qqTaskMaster.init("qq", "../qq", entranceList, creater)
	return nil
}
