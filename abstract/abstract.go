package abstract

import (
	"time"
)



type transInfor[A any, I any, C any] struct {
	isUpload 		 bool // true is upload , false is download
	accesssTime      time.Time // the time of the infor's access
	contentAttr		 A	// should add the content's attribute , such as the bool of isGroup and isPrivate and so on
	pusherIdInfor	 I	// should add the infor's pusher 's identify infor
	getterIdInfor	 I	// should add the infor's getter 's identify infor
	contentSet		 C	// should add the content of the message , such as the text and the rich text and so on
}

type taskInfor[A any, I any, C any, U any] struct {
	userIdInfor U// should add the identifier of the app' user , such as the userid and userpassword and so on
	processId int // the process id of the task
	uploadChannel  chan transInfor[A, I, C] // the channel to upload the infor
	downloadChannel chan transInfor[A, I, C] // the channel to download the infor
	execution func() error // control the task life
}

type taskMasterInfor[A any, I any, C any, U any] struct {
	name string // the name of the app
	existTask map[string] *taskInfor[A, I, C, U] // key is the task identifier
	sourceAddress string // the address of this task' app 
	creater func() (*taskInfor[A, I, C, U],error) // the function to create the task
}


func (taskMaster *taskMasterInfor[A, I, C, U]) init(name string, sourceAddress string, creater func(string) func() (*taskInfor[A, I, C, U],error) ) {
	taskMaster.name = name
	taskMaster.sourceAddress = sourceAddress
	taskMaster.existTask = make(map[string] *taskInfor[A, I, C, U])
	taskMaster.creater = creater(taskMaster.sourceAddress) // the creater return the taskAdder function, the taskAdder function return the processId
}

func (taskMaster *taskMasterInfor[A, I, C, U]) addTask (taskId string) error {
	taskPointer, createError := taskMaster.creater()
	taskMaster.existTask[taskId] = taskPointer		
	return createError	
}

func (taskMaster *taskMasterInfor[A, I, C, U]) deleteTask (taskId string) error {
	deleteError := taskMaster.existTask[taskId].execution()
	if deleteError != nil {
		return deleteError
	}
	delete(taskMaster.existTask, taskId)
	return nil
}