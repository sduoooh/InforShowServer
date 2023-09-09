package main

import (
	"regexp"
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
}

func readPort (url string)  (string ,error) {
	re, _ := regexp.Compile(url + "config.yml")
	data, err := fileOperater(url + "config.yml", fileOperaterOptions{operater: "read"})
	if err != nil {
		return "-1" , err
	}
	return re.FindStringSubmatch(data[0])[1], nil
}

func creater (sourceAddress string, occupiedPort *map[string]portInfor) func(string) (*taskInfor[qqContentAttr, qqHandlerInfor, qqContentSet, qqUserIdInfor], error) {
	return func (entrance string) (*taskInfor[qqContentAttr, qqHandlerInfor, qqContentSet, qqUserIdInfor], error) {
		// 尚待施工，需要先确定特定环境文件夹内容构造；
		// 目前认为还需使在其中读取一些信息，如预期使用的端口号及qq号等，后期或许在此预读数据库


		// 先注起来，没施工完，防止恼人的unused警告
		// port, err := readPort(sourceAddress + entrance + "/")
		// if err != nil {
		// 	return nil, err
		// }

		return nil, nil
	}
}

func qqServerInit() {
	entranceList := map[string]string{ // 应当在文件里读，这里先预存一个方便开发
		"123456789": "123456789",
	}
	qqTaskMaster := taskMasterInfor[qqContentAttr, qqHandlerInfor, qqContentSet, qqUserIdInfor]{}
	qqTaskMaster.init("qq", false, "../qq", entranceList, creater)
}