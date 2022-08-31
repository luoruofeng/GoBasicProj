package task_srv

import "time"

type Task struct {
	Id     uint64
	Name   string
	Create time.Time
}
