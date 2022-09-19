package etcd

import (
	"context"
	"encoding/json"
	"log"
	"strconv"
	"time"

	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/mvcc/mvccpb"
	"github.com/luoruofeng/gobasicproj/common"
	c "github.com/luoruofeng/gobasicproj/master/config"
	"github.com/luoruofeng/gobasicproj/master/task_srv"
)

// 任务管理器
type TaskSrv struct {
	client *clientv3.Client
	kv     clientv3.KV
	lease  clientv3.Lease
}

var (
	EtcdTaskSrv *TaskSrv
)

// 初始化管理器
func InitEtcdTaskSrv() error {
	var (
		config clientv3.Config
		client *clientv3.Client
		kv     clientv3.KV
		lease  clientv3.Lease
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
		log.Println("GEt ETCD client!")
	}

	//test connect
	cxt, _ := context.WithTimeout(context.TODO(), 5*time.Second)
	_, err = client.Get(cxt, "test_key")
	if err != nil {
		log.Println("connect to ETCD failed!")
		return err
	} else {
		log.Println("etcd connect successfully!", c.Cnf.EtcdAddrs)
	}

	// 得到KV和Lease的API子集
	kv = clientv3.NewKV(client)
	lease = clientv3.NewLease(client)

	// 赋值单例
	EtcdTaskSrv = &TaskSrv{
		client: client,
		kv:     kv,
		lease:  lease,
	}
	return nil
}

// 关闭连接
func (t *TaskSrv) CloseTaskSrv() error {
	log.Println("close etcd server connection!")
	return t.client.Close()
}

// 保存任务
func (t *TaskSrv) SaveTask(task *task_srv.Task) (oldTask *task_srv.Task, err error) {
	// 把任务保存到/gobasicproj//任务id -> json
	var (
		taskKey    string
		taskValue  []byte
		putResp    *clientv3.PutResponse
		oldTaskObj task_srv.Task
	)

	// etcd的保存key
	taskKey = common.TaskSaveDir + strconv.FormatUint(task.Id, 10)
	// 任务信息json
	if taskValue, err = json.Marshal(task); err != nil {
		return
	}
	// 保存到etcd
	if putResp, err = t.kv.Put(context.TODO(), taskKey, string(taskValue), clientv3.WithPrevKV()); err != nil {
		return
	}
	// 如果是更新, 那么返回旧值
	if putResp.PrevKv != nil {
		// 对旧值做一个反序列化
		if err = json.Unmarshal(putResp.PrevKv.Value, &oldTaskObj); err != nil {
			err = nil
			return
		}
		oldTask = &oldTaskObj
	}
	return
}

// 获取全部任务
func (t *TaskSrv) GetAllTask() (allTask []*task_srv.Task, err error) {
	var (
		getResp *clientv3.GetResponse
		kvPair  *mvccpb.KeyValue
		task    *task_srv.Task
	)
	if getResp, err = t.kv.Get(context.TODO(), common.TaskSaveDir, clientv3.WithPrefix()); err != nil {
		return
	}
	log.Println(getResp)
	allTask = make([]*task_srv.Task, 0)
	for _, kvPair = range getResp.Kvs {
		task = &task_srv.Task{}
		if err = json.Unmarshal(kvPair.Value, task); err != nil {
			err = nil
			continue
		}
		allTask = append(allTask, task)
	}
	log.Println(allTask)
	return
}

// 删除任务
func (t *TaskSrv) DeleteTask(id uint64) (oldTask *task_srv.Task, err error) {
	var (
		taskKey    string
		delResp    *clientv3.DeleteResponse
		oldTaskObj task_srv.Task
	)

	// etcd中保存任务的key
	taskKey = common.TaskSaveDir + strconv.FormatUint(id, 10)

	// 从etcd中删除它
	if delResp, err = t.kv.Delete(context.TODO(), taskKey, clientv3.WithPrevKV()); err != nil {
		return
	}

	// 返回被删除的任务信息
	if len(delResp.PrevKvs) != 0 {
		// 解析一下旧值, 返回它
		if err = json.Unmarshal(delResp.PrevKvs[0].Value, &oldTaskObj); err != nil {
			err = nil
			return
		}
		oldTask = &oldTaskObj
	}
	return
}
