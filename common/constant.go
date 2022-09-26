package common

const (
	TaskSaveDir = "/gobasicproj/task/"
)

type TaskEventType int

const (
	EventPut    TaskEventType = 0
	EventDelete TaskEventType = 1
)
