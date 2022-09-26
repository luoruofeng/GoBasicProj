package task_srv

import "time"

type Task struct {
	Id     uint64    `json:"id"`
	Name   string    `json:"name"`
	Create time.Time `json:"create"`
}
