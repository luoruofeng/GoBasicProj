package etcd

import (
	"testing"
	"time"

	c "github.com/luoruofeng/gobasicproj/master/config"
	"github.com/luoruofeng/gobasicproj/master/task_srv"
)

var (
	config_path = "./test_master_config.json"
	task        = task_srv.Task{
		Id:     uint64(123),
		Name:   "luoruofeng",
		Create: time.Now(),
	}
)

// go test -run Etcd
// go test -run Etcd/InitEtcdTaskSrv
func TestEtcd(t *testing.T) {
	t.Run("InitEtcdTaskSrv", func(t *testing.T) {
		err := c.InitConfig(config_path)
		if err != nil {
			t.Fatalf("Reading config file failed! err:%v", err)
		}

		err = InitEtcdTaskSrv()
		if err != nil {
			t.Fatalf("Calling InitEtcdTaskSrv funcation failed! err:%v", err)
		}
	})
	t.Run("SaveTask", func(t *testing.T) {
		o, err := EtcdTaskSrv.SaveTask(&task)
		if err != nil {
			t.Fatalf("ETCD save task failed! err:%v", err)
		} else {
			t.Log("ETCD save task successfully! oldvalue:", o)
		}
	})
	t.Run("GetAllTask", func(t *testing.T) {
		ts, err := EtcdTaskSrv.GetAllTask()
		if err != nil {
			t.Fatalf("ETCD get all task failed! err:%v", err)
		} else {
			t.Log("ETCD get all task successfully!:")
		}

		t.Log("get all tasks")
		for _, task := range ts {
			t.Log(task)
		}
	})

	// 删除任务
	// t.Run("DeleteTask", func(t *testing.T) {
	// 	o, err := EtcdTaskSrv.DeleteTask(task.Id)
	// 	if err != nil {
	// 		t.Fatalf("ETCD delete task failed! err:%v", err)
	// 	} else {
	// 		t.Log("ETCD delete task successfully! oldvalue:", o)
	// 	}
	// })

}
