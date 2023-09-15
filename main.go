package main

import (
	"fmt"
	"os"
    "os/signal"
    "syscall"
)

func main() {
	getterChannel := make(chan unionTransInfor, 100)
	uploadChannel := make(chan unionTransInfor, 100)
	sysSignal := make(chan os.Signal, 1)
	done := make(chan bool, 1)
	signal.Notify(sysSignal, syscall.SIGINT, syscall.SIGTERM)
	qqTaskMaster, err := qqServerInit(uploadChannel, getterChannel)
	if err != nil {
		panic(err)

	}
	qqTaskEnd, qqTaskStartErr := qqServerStart(qqTaskMaster, []string{"main"})
	if qqTaskStartErr != nil {
		panic(qqTaskStartErr)
	}
	fmt.Println("server start")

	// 任务结束信号
	go func() {
		<-sysSignal
		qqTaskEnd <- true
		done <- true
	}()
	<-done
	fmt.Println("server end")

}
