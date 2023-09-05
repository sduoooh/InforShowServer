package infor

import (
	"time"
)



type transInfor[A any, I any, C any] struct {
	accesssTime      time.Time // the time of the infor's access
	contentAttr		 A	// should add the content's attribute , such as the bool of isGroup and isPrivate and so on
	inforIdInfor	 I	// should add the identify infor of this infor's target or source
	contentSet		 C	// should add the content of the message , such as the text and the rich text and so on
}

type taskInfor[A any, I any, C any, U any] struct {
	userIdInfor U// should add the identifier of the app' user , such as the userid and userpassword and so on
	processId int // the process id of the task
	uploadChannel  chan transInfor[A, I, C] // the channel to upload the infor
	downloadChannel chan transInfor[A, I, C] // the channel to download the infor
	endChannel    chan bool // control the task life
}

type taskMasterInfor[A any, I any, C any, U any] struct {
	name string // the name of the app
	existTask map[string] *taskInfor[A, I, C, U] // key is the task identifier
	sourceAddress string // the address of this task' app 
	creater func() *taskInfor[A, I, C, U] // the function to create the task
}


func (taskMaster *taskMasterInfor[A, I, C, U]) init(name string, sourceAddress string, creater func(string) func() *taskInfor[A, I, C, U] ) {
	taskMaster.name = name
	taskMaster.sourceAddress = sourceAddress
	taskMaster.existTask = make(map[string] *taskInfor[A, I, C, U])
	taskMaster.creater = creater(taskMaster.sourceAddress) // the creater return the taskAdder function, the taskAdder function return the processId
}

func (taskMaster *taskMasterInfor[A, I, C, U]) addTask (taskId string) {
	taskPointer := taskMaster.creater()
	taskMaster.existTask[taskId] = taskPointer				
}

func (taskMaster *taskMasterInfor[A, I, C, U]) deleteTask (taskId string) {
	taskMaster.existTask[taskId].endChannel <- true
	delete(taskMaster.existTask, taskId)
}