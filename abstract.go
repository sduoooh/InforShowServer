package main

import (
	"errors"
)

type portInfor struct {
	ownerTaskId string
	isUsed      bool
}

type transInfor[A any, I any, C any] struct {
	isUpload      bool   // true is upload , false is download
	accesssTime   string // the time of the infor's access
	contentAttr   A      // should add the content's attribute , such as the bool of isGroup and isPrivate and so on
	pusherIdInfor I      // should add the infor's pusher 's identify infor
	getterIdInfor I      // should add the infor's getter 's identify infor
	contentSet    C      // should add the content of the message , such as the text and the rich text and so on
}

type unionTransInfor struct {
	isUpload   bool   // true is upload , false is download
	targetApp  string // the target app's name
	targetId   string // the target app's id, if the target app is only, it is main.
	transInfor string // should convert the transInfor to string
}

type taskInfor[A any, I any, C any, U any] struct {
	occupiedPort    []string                  // the port that the task occupied
	userIdInfor     U                         // should add the identifier of the app' user , such as the userid and userpassword and so on
	processId       string                    // the process id of the task
	uploadChannel   chan transInfor[A, I, C] // the channel to upload the infor
	downloadChannel chan transInfor[A, I, C] // the channel to download the infor
	execution       func() error              // control the task life
}

type taskMasterInfor[A any, I any, C any, U any] struct {
	name                 string                                       // the name of the app
	sourceAddress        string                                       // the rootAddress of this task' app
	entranceList         map[string]string                            // the entrance of the task
	existTask            map[string]*taskInfor[A, I, C, U]            // key is the task identifier
	occupiedPort         map[string]portInfor                         // the port that the task occupied
	InnerUploadChannel   chan transInfor[A, I, C]                    // the channel to upload the infor
	InnerDownloadChannel chan transInfor[A, I, C]                    // the channel to download the infor
	OuterUploadChannel   chan unionTransInfor                   // the channel to upload the infor
	OuterDownloadChannel chan unionTransInfor                    // the channel to download the infor
	creater              func(string) (*taskInfor[A, I, C, U], error) // the function to create the task, if isHorizontalExpansion is true, the parameter of the creater is 'main', else is the task's entrance, aka task's unique environment folder name
}

func (taskMaster *taskMasterInfor[A, I, C, U]) init(name string, sourceAddress string, entranceList map[string]string, creater func(entrance string,entranceList *map[string]portInfor, innerDownloadChannel, innerUploadChannel chan transInfor[A, I, C]) func(string) (*taskInfor[A, I, C, U], error), outterUpper, outterGetter chan unionTransInfor) {
	taskMaster.name = name
	taskMaster.InnerUploadChannel = make(chan transInfor[A, I, C], 5)
	taskMaster.InnerDownloadChannel = make(chan transInfor[A, I, C], 10)
	taskMaster.OuterUploadChannel = outterUpper
	taskMaster.OuterDownloadChannel = outterGetter
	taskMaster.entranceList = entranceList
	taskMaster.sourceAddress = sourceAddress
	taskMaster.existTask = make(map[string]*taskInfor[A, I, C, U])
	taskMaster.occupiedPort = make(map[string]portInfor)
	taskMaster.creater = creater(taskMaster.sourceAddress, &taskMaster.occupiedPort, taskMaster.InnerDownloadChannel, taskMaster.InnerUploadChannel) // the creater return the taskAdder function, the taskAdder function return the processId
}

func (taskMaster *taskMasterInfor[A, I, C, U]) addTask(taskId string) error {
	taskPointer, createError := taskMaster.creater(taskMaster.entranceList[taskId])
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
	taskMaster.existTask[taskId].execution()
	for _, i := range taskMaster.existTask[taskId].occupiedPort {
		delete(taskMaster.occupiedPort, i)
	}
	delete(taskMaster.existTask, taskId)
	return nil
}