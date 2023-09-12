package main

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"

	"github.com/lxzan/gws"
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

type QQGwsHandler struct {
	gws.BuiltinEventHandler
	logAddr         string
	downloadChannel *chan transInfor[qqContentAttr, qqHandlerInfor, qqContentSet]
}

func (h *QQGwsHandler) OnMessage(c *gws.Conn, message *gws.Message) {
	// 获取消息类型和内容
	opcode := message.Opcode
	payload := message.Data

	fmt.Println("receive message")
	fmt.Println(payload.String())
	// 根据不同的消息类型进行处理
	switch opcode {
	case gws.OpcodeText:
		fmt.Println(payload.String())
	case gws.OpcodeBinary:
		log, err := fileOperater(h.logAddr, fileOperaterOptions{operater: "read"})
		if err != nil {
			panic(err)
		}
		fileOperater(h.logAddr, fileOperaterOptions{operater: "write", regexp: ".*", replacement: log[0] + payload.String()})
	default:
		panic(errors.New("unknown opcode"))
	}
}

func (h *QQGwsHandler) OnClose(c *gws.Conn, err error) {
	fmt.Println("close")
}

func (h *QQGwsHandler) OnOpen(c *gws.Conn) {
	fmt.Println("open")
}

func creater(sourceAddress string, occupiedPort *map[string]portInfor) func(string) (*taskInfor[qqContentAttr, qqHandlerInfor, qqContentSet, qqUserIdInfor], error) {
	return func(entrance string) (*taskInfor[qqContentAttr, qqHandlerInfor, qqContentSet, qqUserIdInfor], error) {
		// 目前认为还需使在其中读取一些信息，如预期使用的端口号及qq号等，后期或许在此预读数据库
		

		// 创建task
		task := taskInfor[qqContentAttr, qqHandlerInfor, qqContentSet, qqUserIdInfor]{}
		uploadChannel := make(chan transInfor[qqContentAttr, qqHandlerInfor, qqContentSet], 1)
		downloadChannel := make(chan transInfor[qqContentAttr, qqHandlerInfor, qqContentSet], 1)
		task.occupiedPort = make([]string, 1)

		task.uploadChannel = &uploadChannel
		task.downloadChannel = &downloadChannel

		portInfor, err := fileOperater(sourceAddress+entrance+"/qqPortList.txt", fileOperaterOptions{operater: "read", createble: false})
		if err != nil {
			return nil, errors.New("can't read qqPortList")
		}
		task.occupiedPort[0] = portInfor[0]

		userId, getUserIdErr := strconv.Atoi(entrance)
		if getUserIdErr != nil {
			return nil, errors.New(getUserIdErr.Error() + "getUserId error")
		}
		task.userIdInfor = qqUserIdInfor{userId}

		pid, startErr := start("go-cqhttp_windows_amd64.exe", sourceAddress + entrance)
		if startErr != nil {
			return nil, errors.New(startErr.Error() + "start error")
		}
		task.processId = pid

		task.execution = func() error {
			killErr := kill(task.processId)
			if killErr != nil {
				return errors.New(killErr.Error() + "kill error")
			}
			return nil
		}

		// 连接进程
		// 当前进度： 连接成功，但是OnMessage函数无法被调用,不知何原因；
		// 同时OnOpen、OnClose函数也无法被调用
		QQGwsHandler := QQGwsHandler{logAddr: sourceAddress + entrance + "/log.txt", downloadChannel: task.downloadChannel}
		app, _, gwsErr := gws.NewClient(&QQGwsHandler, &gws.ClientOption{Addr: "ws://127.0.0.1:" + portInfor[0]})
		if gwsErr != nil {
			return nil, errors.New(gwsErr.Error()+ "gws error")
		}
		go app.ReadLoop()

		return &task, nil
	}
}

func qqServerInit() (*taskMasterInfor[qqContentAttr, qqHandlerInfor, qqContentSet, qqUserIdInfor], error) {
	qqMap, err := fileOperater("../qq/qqMap.txt", fileOperaterOptions{operater: "read", createble: false})
	if err != nil {
		return nil, errors.New("can't read qqMap.txt")
	}
	entranceList := dataLoad(qqMap[0])
	qqTaskMaster := taskMasterInfor[qqContentAttr, qqHandlerInfor, qqContentSet, qqUserIdInfor]{}
	qqTaskMaster.init("qq", "../qq/", entranceList, creater)
	return &qqTaskMaster, nil
}
