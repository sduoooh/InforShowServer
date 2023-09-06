package inforShowServer

import (
)

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

func qqServerStart() {
	//qqTaskMaster := taskMasterInfor[qqContentAttr, qqHandlerInfor, qqContentSet, qqUserIdInfor]{}
	//qqTaskMaster.init("qq", )
}