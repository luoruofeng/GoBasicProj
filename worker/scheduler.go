package worker

import (
	"log"
	"time"

	"github.com/luoruofeng/gobasicproj/common"
	"github.com/luoruofeng/gobasicproj/worker/task_srv/etcd"
)

type Scheduler struct {
	eventChan chan *etcd.TaskEvent
}

var Sched *Scheduler

func InitScheduler() {
	ec := make(chan *etcd.TaskEvent, 1000)

	Sched = &Scheduler{
		eventChan: ec,
	}

	go Sched.SchedulerLoop()
}

func (s *Scheduler) SchedulerLoop() {
	timer := time.NewTimer(time.Second * 10)

	for {
		select {
		case te := <-s.eventChan:
			if te.Event == common.EventDelete {
				//TODO 做一些事情
				log.Println("delete event:", te)
			}

			if te.Event == common.EventPut {
				//TODO 做一些事情
				log.Println("put event:", te)
			}
		case <-timer.C:
			log.Println("do")
		}
	}
}

func (s *Scheduler) PushEvent(te *etcd.TaskEvent) {
	s.eventChan <- te
}
