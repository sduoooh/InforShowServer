package main

func main() {
	getterChannel := make(chan unionTransInfor, 100)
	uploadChannel := make(chan unionTransInfor, 100)
	qqTaskMaster, err := qqServerInit(uploadChannel, getterChannel)
	if err != nil {
		panic(err)

	}
	qqTaskStartErr, _ := qqServerStart(qqTaskMaster, []string{"main"})
	if qqTaskStartErr != nil {
		panic(qqTaskStartErr)
	}
	// todo
	// 接受signal，结束qqServer
	for {
		
	}

}
