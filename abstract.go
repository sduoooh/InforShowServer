package inforShowServer

import (
	"errors"
	"time"
)

type portInfor struct {
	ownerTaskId string
	isUsed      bool
}

type transInfor[A any, I any, C any] struct {
	isUpload      bool      // true is upload , false is download
	accesssTime   time.Time // the time of the infor's access
	contentAttr   A         // should add the content's attribute , such as the bool of isGroup and isPrivate and so on
	pusherIdInfor I         // should add the infor's pusher 's identify infor
	getterIdInfor I         // should add the infor's getter 's identify infor
	contentSet    C         // should add the content of the message , such as the text and the rich text and so on
}

type taskInfor[A any, I any, C any, U any] struct {
	occupiedPort    []string                 // the port that the task occupied
	userIdInfor     U                        // should add the identifier of the app' user , such as the userid and userpassword and so on
	processId       string                   // the process id of the task
	uploadChannel   chan transInfor[A, I, C] // the channel to upload the infor
	downloadChannel chan transInfor[A, I, C] // the channel to download the infor
	execution       func() error             // control the task life
}

type taskMasterInfor[A any, I any, C any, U any] struct {
	name                  string                                       // the name of the app
	isHorizontalExpansion bool                                         // if the task is horizontal expansion
	entranceList          []string                                     // the entrance of the task
	existTask             map[string]*taskInfor[A, I, C, U]            // key is the task identifier
	occupiedPort          map[string]portInfor                         // the port that the task occupied
	sourceAddress         string                                       // the address of this task' app
	creater               func(string) (*taskInfor[A, I, C, U], error) // the function to create the task， if isHorizontalExpansion is true, the parameter of the creater is ''
}

func (taskMaster *taskMasterInfor[A, I, C, U]) init(name string, isHorizontalExpansion bool, sourceAddress string, entranceList []string, creater func(string, *map[string]portInfor) func(string) (*taskInfor[A, I, C, U], error)) {
	taskMaster.name = name
	taskMaster.sourceAddress = sourceAddress
	taskMaster.existTask = make(map[string]*taskInfor[A, I, C, U])
	taskMaster.occupiedPort = make(map[string]portInfor)
	taskMaster.creater = creater(taskMaster.sourceAddress, &taskMaster.occupiedPort) // the creater return the taskAdder function, the taskAdder function return the processId
}

func (taskMaster *taskMasterInfor[A, I, C, U]) addTask(taskId string) error {
	taskPointer, createError := taskMaster.creater(taskMaster.sourceAddress)
	if createError != nil {
		return createError
	}
	taskMaster.existTask[taskId] = taskPointer
	for index, i := range taskPointer.occupiedPort {
		if port, ok := taskMaster.occupiedPort[i]; ok {
			if port.isUsed {
				for j, k := range taskPointer.occupiedPort {
					if j < index {
						taskMaster.occupiedPort[k] = portInfor{ownerTaskId: taskId, isUsed: false}
					} else {
						break
					}
				}
				taskMaster.deleteTask(taskId)
				return errors.New("the port " + i + " is occupied by " + taskMaster.occupiedPort[i].ownerTaskId + " and can't be used by " + taskId)
			}
		}
		taskMaster.occupiedPort[i] = portInfor{ownerTaskId: taskId, isUsed: true}
	}
	return createError
}

func (taskMaster *taskMasterInfor[A, I, C, U]) deleteTask(taskId string) error {
	deleteError := taskMaster.existTask[taskId].execution()
	if deleteError != nil {
		return deleteError
	}
	for _, i := range taskMaster.existTask[taskId].occupiedPort {
		delete(taskMaster.occupiedPort, i)
	}
	delete(taskMaster.existTask, taskId)
	return deleteError
}
