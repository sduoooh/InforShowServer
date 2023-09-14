package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"time"

	"github.com/lxzan/gws"
)

// 尚待与abstract协同完善

type qqContentAttr struct {
	isGroup     bool
	isPrivate   bool
	isAnonymity bool
}

type qqHandlerInfor struct {

	// nil if isGroup is false
	groupId   string
	groupName string

	// always not nil
	userId   string
	userName string
}

type qqContentSet struct {
	text string
}

type qqUserIdInfor struct {
	userId string
}

type QQGwsHandler struct {
	gws.BuiltinEventHandler
	logAddr         string
	downloadChannel chan transInfor[qqContentAttr, qqHandlerInfor, qqContentSet]
}

// 在当前需求下，直接拿正则取效果更好，参见：https://github.com/sduoooh/GoJsonDataTest
// "和\前会加\，以此判断
// []会用ASCII码表示，可以直接取

var postTypeGetter = regexp.MustCompile(`"post_type":"(.*?)(?:\\\\)*"`)
var messageGetter = regexp.MustCompile(`"message":"(.*?)(?:\\\\)*"`)
var timeGetter = regexp.MustCompile(`"time":([0-9]{10})`)
var messageTypeGetter = regexp.MustCompile(`"message_type":"(.*?)(?:\\\\)*"`)
var subTypeGetter = regexp.MustCompile(`"sub_type":"(.*?)(?:\\\\)*"`)
var groupIdGetter = regexp.MustCompile(`"group_id":([0-9]+)`)
var userIdGetter = regexp.MustCompile(`"user_id":([0-9]+)`)
var nicknameGetter = regexp.MustCompile(`"nickname":"(.*?)(?:\\\\)*"`)
var targetIdGetter = regexp.MustCompile(`"target_id":([0-9]+)`)

func (h *QQGwsHandler) OnMessage(c *gws.Conn, message *gws.Message) {
	// 获取消息类型和内容
	opcode := message.Opcode
	payload := message.Data.String()

	// 根据不同的消息类型进行处理
	switch opcode {
	case gws.OpcodeText:
		switch postTypeGetter.FindStringSubmatch(payload)[1] {
		case "message":

			// todo： 针对键名拿数据，然后放chan里
			fmt.Println(messageGetter.FindStringSubmatch(payload)[1])
			isGroup := messageTypeGetter.FindStringSubmatch(payload)[1] == "group"

			// 私聊则无群id
			if !isGroup {
				payload += "\"group_id\":" + "0"
			}

			transInfors := transInfor[qqContentAttr, qqHandlerInfor, qqContentSet]{
				isUpload:    false,
				accesssTime: timeGetter.FindStringSubmatch(payload)[1],
				contentAttr: qqContentAttr{
					isGroup:     isGroup,
					isPrivate:   messageTypeGetter.FindStringSubmatch(payload)[1] == "private",
					isAnonymity: subTypeGetter.FindStringSubmatch(payload)[1] == "anonymous",
				},
				pusherIdInfor: qqHandlerInfor{
					userId:    userIdGetter.FindStringSubmatch(payload)[1],
					userName:  nicknameGetter.FindStringSubmatch(payload)[1],
					groupId:   groupIdGetter.FindStringSubmatch(payload)[1], 
					groupName: "无", // 尚待支持，考虑http请求或者直接在websocket里询问
				},
				getterIdInfor: qqHandlerInfor{
					// 其他不需要
					userId: targetIdGetter.FindStringSubmatch(payload)[1],
				},
				contentSet: qqContentSet{
					text: payload,
				},
			}
			h.downloadChannel <- transInfors
			fmt.Println("send to channel", transInfors)

		default:
			fmt.Println("unsupported message type: ", postTypeGetter.FindStringSubmatch(payload)[1])
		}

	case gws.OpcodeBinary:
		log, err := fileOperater(h.logAddr, fileOperaterOptions{operater: "read"})
		if err != nil {
			panic(err)
		}
		fileOperater(h.logAddr, fileOperaterOptions{operater: "write", regexp: ".*", replacement: log[0] + payload})
	default:
		panic(errors.New("unknown opcode"))
	}
}

func creater(sourceAddress string, occupiedPort *map[string]portInfor, innerGetter, innerUpper chan transInfor[qqContentAttr, qqHandlerInfor, qqContentSet]) func(string) (*taskInfor[qqContentAttr, qqHandlerInfor, qqContentSet, qqUserIdInfor], error) {
	return func(entrance string) (*taskInfor[qqContentAttr, qqHandlerInfor, qqContentSet, qqUserIdInfor], error) {
		// 目前认为还需使在其中读取一些信息，如预期使用的端口号及qq号等，后期或许在此预读数据库

		// 创建task
		// 下载信息最大缓存10条
		task := taskInfor[qqContentAttr, qqHandlerInfor, qqContentSet, qqUserIdInfor]{}
		task.occupiedPort = make([]string, 1)

		task.uploadChannel = innerUpper
		task.downloadChannel = innerGetter

		portInfor, err := fileOperater(sourceAddress+entrance+"/qqPortList.txt", fileOperaterOptions{operater: "read", createble: false})
		if err != nil {
			return nil, errors.New("can't read qqPortList")
		}
		task.occupiedPort[0] = portInfor[0]

		userId := entrance
		task.userIdInfor = qqUserIdInfor{userId}

		pid, startErr := start("go-cqhttp_windows_amd64.exe", sourceAddress+entrance)
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
		time.Sleep(8 * time.Second) // 必要的时延，包括启动的5秒在内

		// 连接进程
		QQGwsHandler := QQGwsHandler{logAddr: sourceAddress + entrance + "/log.txt", downloadChannel: task.downloadChannel}
		app, _, gwsErr := gws.NewClient(&QQGwsHandler, &gws.ClientOption{Addr: "ws://127.0.0.1:" + portInfor[0]})
		if gwsErr != nil {
			return nil, errors.New(gwsErr.Error() + "gws error")
		}
		go app.ReadLoop()

		return &task, nil
	}
}

func qqServerInit(outterUpper, outterGetter chan unionTransInfor) (*taskMasterInfor[qqContentAttr, qqHandlerInfor, qqContentSet, qqUserIdInfor], error) {
	qqMap, err := fileOperater("../qq/qqMap.txt", fileOperaterOptions{operater: "read", createble: false})
	if err != nil {
		return nil, errors.New("can't read qqMap.txt")
	}
	entranceList := dataLoad(qqMap[0])
	qqTaskMaster := taskMasterInfor[qqContentAttr, qqHandlerInfor, qqContentSet, qqUserIdInfor]{}
	qqTaskMaster.init("qq", "../qq/", entranceList, creater, outterUpper, outterGetter)

	return &qqTaskMaster, nil
}

func qqServerStart(qqTaskMaster *taskMasterInfor[qqContentAttr, qqHandlerInfor, qqContentSet, qqUserIdInfor], nameList []string) (error, chan bool) {
	end := make(chan bool, 1)
	for _, i := range nameList {
		addTaskError := qqTaskMaster.addTask(i)
		if addTaskError != nil {
			return addTaskError, nil
		}
	}
	go func() {
		for {
			select {
			// 有消息来，就转化为unionTransInfor，放入outterGetter
			case qqTransInfor := <-qqTaskMaster.InnerDownloadChannel:
				content, _ := json.Marshal(qqTransInfor) // 显然得自己写
				unionTransInfor := unionTransInfor{
					isUpload:   false,
					targetApp:  "qq",
					targetId:   qqTransInfor.getterIdInfor.userId,
					transInfor: string(content),
				}
				qqTaskMaster.OuterDownloadChannel <- unionTransInfor
				fmt.Println("qqTransInfor", unionTransInfor)
				
			case qqTransInfor2 := <-qqTaskMaster.OuterUploadChannel:
				fmt.Println("qqTransInfor2", qqTransInfor2)
				// todo
				continue
			case <-end:
				for _, i := range qqTaskMaster.existTask {
					i.execution()
				}
				fmt.Println("qqServer end")
				break
			}
		}
	}()
	return nil, end
}
