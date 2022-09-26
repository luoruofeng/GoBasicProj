package etcd

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/coreos/etcd/clientv3"
	"github.com/luoruofeng/gobasicproj/common"
	c "github.com/luoruofeng/gobasicproj/worker/config"

	
)

// 任务管理器
type TaskSrv struct {
	client  *clientv3.Client
	kv      clientv3.KV
	lease   clientv3.Lease
	watcher clientv3.Watcher
}

var (
	EtcdTaskSrv *TaskSrv
)

type Task struct {
	Id     uint64    `json:"id"`
	Name   string    `json:"name"`
	Create time.Time `json:"create"`
}

type TaskEvent struct {
	Task  *Task
	Event common.TaskEventType
}

func BuildTaskEvent(eventType common.TaskEventType, task *Task) (taskEvent *TaskEvent) {
	return &TaskEvent{
		Event: eventType,
		Task:  task,
	}
}

// 初始化管理器
func InitEtcdTaskSrv() error {
	var (
		config  clientv3.Config
		client  *clientv3.Client
		kv      clientv3.KV
		lease   clientv3.Lease
		watcher clientv3.Watcher
	)

	// 初始化配置
	config = clientv3.Config{
		Endpoints:   c.Cnf.EtcdAddrs,                                         // 集群地址
		DialTimeout: time.Duration(c.Cnf.EtcdDialTimeout) * time.Millisecond, // 连接超时
	}

	// 建立连接
	client, err := clientv3.New(config)
	if err != nil {
		return err
	} else {
		log.Println("GEt ETCD client!", c.Cnf.EtcdAddrs)
	}

	//test connect
	cxt, _ := context.WithTimeout(context.TODO(), 5*time.Second)
	_, err = client.Get(cxt, "test_key")
	if err != nil {
		log.Println("connect to ETCD failed!")
		return err
	} else {
		log.Println("etcd connect successfully!")
	}

	// 得到KV和Lease的API子集
	kv = clientv3.NewKV(client)
	lease = clientv3.NewLease(client)
	watcher = clientv3.NewWatcher(client)

	// 赋值单例
	EtcdTaskSrv = &TaskSrv{
		client:  client,
		kv:      kv,
		lease:   lease,
		watcher: watcher,
	}
	return nil
}

// 关闭连接
func (t *TaskSrv) CloseTaskSrv() error {
	log.Println("close etcd server connection!")
	return t.client.Close()
}

func (t *TaskSrv) watchTasks() (err error) {
	// 获取以后的所有任务
	getResp, err := t.kv.Get(context.TODO(), common.TaskSaveDir, clientv3.WithPrefix())
	if err != nil {
		return
	}

	var task Task
	for _, kvpair := range getResp.Kvs {
		if e := json.Unmarshal(kvpair.Value, &task); e == nil {
			te := BuildTaskEvent(common.EventPut, &task)
			
		}
	}

}
