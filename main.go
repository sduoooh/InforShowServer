package main

func main() {
	qqTaskMaster, err := qqServerInit()
	if err != nil {
		panic(err)
	}
	addTaskError := qqTaskMaster.addTask("main")
	if addTaskError != nil {
		panic(addTaskError)
	}
}