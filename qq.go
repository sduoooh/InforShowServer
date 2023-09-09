package inforShowServer

import (
	"regexp"
	"os"
	"io"
)

// 注意，为方便传递，qq中的id都是int类型，但是本程序内部传递的id为string.
// 尚待与abstract协同完善

type qqContentAttr struct {
	isGroup bool
	isPrivate bool
	isAnonymity bool
	isRichText bool
}

type qqHandlerInfor struct {

	// nil if isGroup is false
	groupId int
	groupName string

	// always not nil
	userId int
	userName string
}

type qqContentSet struct {
	text string

	// nil if isRichText is false
	richText []string // if some rich text exists, such as pics or videos, it will be stored in this slice
}

type qqUserIdInfor struct {
	userId int
	userPassword string
}

func readPort (url string)  (string ,error) {
	re, _ := regexp.Compile("[^#]address: 127.0.0.1:([0-9]+[^#])")
	file, err1 := os.Open(url + "config.yml")
	if err1 != nil {
		return "-1" , err1
	}
	data, err2 := io.ReadAll(file)
	if err2 != nil {
		return "-1" , err2
	}
	return re.FindStringSubmatch(string(data))[1], nil

}

func qqServerStart() {
	qqTaskMaster := taskMasterInfor[qqContentAttr, qqHandlerInfor, qqContentSet, qqUserIdInfor]{}
	creater := func(sourceAddress string) func(string) (*taskInfor[qqContentAttr, qqHandlerInfor, qqContentSet, qqUserIdInfor], error) {
		return func(taskId string) (*taskInfor[qqContentAttr, qqHandlerInfor, qqContentSet, qqUserIdInfor], error) {
			return nil, nil
		}
	}
	qqTaskMaster.init("qq", false, "../qq", ["/123456"], )
}