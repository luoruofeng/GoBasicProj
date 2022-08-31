package master

import (
	"encoding/json"
	"log"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/luoruofeng/gobasicproj/common"
	c "github.com/luoruofeng/gobasicproj/master/config"
	"github.com/luoruofeng/gobasicproj/master/task_srv"
	e "github.com/luoruofeng/gobasicproj/master/task_srv/etcd"
)

type ApiServer struct {
	HttpServer http.Server
}

var ApiServ *ApiServer

func InitApiSever() error {
	mux := http.NewServeMux()
	mux.HandleFunc("/task", handleGetTask)
	mux.HandleFunc("/tasks", handleGetAllTask)
	mux.HandleFunc("/savetask", handleSaveTask)
	mux.HandleFunc("/deltask", handleDelTask)

	s := http.Server{
		Handler:      mux,
		ReadTimeout:  time.Duration(c.Cnf.ServerReadTimeout),
		WriteTimeout: time.Duration(c.Cnf.ServerWriteTimeout),
	}

	l, err := net.Listen("tcp", ":"+strconv.Itoa(c.Cnf.ServerPort))
	if err != nil {
		return err
	}

	ApiServ = &ApiServer{
		HttpServer: s,
	}

	go s.Serve(l)
	return nil

}

func handleSaveTask(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		if err := r.ParseForm(); err != nil {
			if bytes, err := common.BuildResponse(-1, err.Error(), nil); err == nil {
				w.Write(bytes)
				return
			}
		}
		strtask := r.PostForm.Get("task")
		var task *task_srv.Task
		if err := json.Unmarshal([]byte(strtask), task); err != nil {
			if bytes, err := common.BuildResponse(-1, err.Error(), nil); err == nil {
				w.Write(bytes)
				return
			}
		}
		if o, err := e.EtcdTaskSrv.SaveTask(task); err != nil {
			if bytes, err := common.BuildResponse(-1, err.Error(), nil); err == nil {
				w.Write(bytes)
				return
			}
		} else {
			if bytes, err := common.BuildResponse(0, "success", o); err == nil {
				w.Write(bytes)
			}
		}
	}
}

func handleGetTask(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
	}
}

func handleGetAllTask(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		ts, err := e.EtcdTaskSrv.GetAllTask()
		if err != nil {
			log.Print("Getting all task failed!")
			if bytes, err := common.BuildResponse(-1, err.Error(), nil); err == nil {
				w.Write(bytes)
				return
			}
		}
		if bytes, err := common.BuildResponse(0, "success", ts); err == nil {
			w.Write(bytes)
		}
	}
}

func handleDelTask(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodDelete:
	}
}
