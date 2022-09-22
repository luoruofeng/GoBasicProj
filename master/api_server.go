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
		ReadTimeout:  time.Duration(c.Cnf.ServerReadTimeout * int(time.Millisecond)),
		WriteTimeout: time.Duration(c.Cnf.ServerWriteTimeout * int(time.Millisecond)),
	}

	listenr, err := net.Listen("tcp", ":"+strconv.Itoa(c.Cnf.ServerPort))
	if err != nil {
		return err
	}

	ApiServ = &ApiServer{
		HttpServer: s,
	}

	go s.Serve(listenr)
	return nil

}

// curl --location --request POST 'http://localhost:8333/savetask' \
// --header 'Content-Type: application/json' \
// --data-raw '{"id":110, "name":"abc", "create":"1985-04-12T23:20:50.52Z"}'
func handleSaveTask(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		br := r.Body
		defer br.Close()

		var task task_srv.Task
		d := json.NewDecoder(br)
		err := d.Decode(&task)
		if err != nil {
			if bytes, err := common.BuildResponse(-1, err.Error(), nil); err == nil {
				w.Write(bytes)
				return
			}
		}

		if o, err := e.EtcdTaskSrv.SaveTask(&task); err != nil {
			if bytes, err := common.BuildResponse(-1, err.Error(), nil); err == nil {
				w.Write(bytes)
				return
			}
		} else {
			if bytes, err := common.BuildResponse(0, "success", o); err == nil {
				w.Write(bytes)
				return
			}
		}
	}
}

// curl --location --request GET 'http://localhost:8333/task?id=123'
func handleGetTask(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		err := r.ParseForm()
		if err != nil {
			if bytes, err := common.BuildResponse(-1, err.Error(), nil); err == nil {
				w.Write(bytes)
			}
			return
		}

		strid := r.FormValue("id")
		id, err := strconv.ParseUint(strid, 10, 64)
		if err != nil {
			if bytes, err := common.BuildResponse(-1, err.Error(), nil); err == nil {
				w.Write(bytes)
			}
			return
		}

		t, err := e.EtcdTaskSrv.GetTaskById(id)
		if err != nil {
			if bytes, err := common.BuildResponse(-1, err.Error(), nil); err == nil {
				w.Write(bytes)
			}
			return
		} else {
			if bytes, err := common.BuildResponse(0, "success", t); err == nil {
				w.Write(bytes)
				return
			}
		}
	}
}

// curl --location --request GET 'http://localhost:8333/tasks'
func handleGetAllTask(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		ts, err := e.EtcdTaskSrv.GetAllTask()
		if err != nil {
			goto ERR
		}
		if bytes, err := common.BuildResponse(0, "success", ts); err == nil {
			w.WriteHeader(http.StatusOK)
			w.Write(bytes)
		} else {
			goto ERR
		}
		return
	ERR:
		if bytes, err := common.BuildResponse(-1, err.Error(), nil); err == nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(bytes)
			return
		}
	}
}

// curl --location --request DELETE 'http://localhost:8333/deltask' \
// --header 'Content-Type: application/json' \
// --data-raw '{
//     "id": 111
// }'
func handleDelTask(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodDelete:
		m := make(map[string]uint64)
		err := json.NewDecoder(r.Body).Decode(&m)
		if err != nil {
			log.Println(err)
			return
		}
		if len(m) < 1 || m["id"] <= 0 {
			if bytes, err := common.BuildResponse(-1, "param is wrong!", nil); err == nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write(bytes)
				return
			}
		}

		o, err := e.EtcdTaskSrv.DeleteTask(m["id"])
		if err != nil {
			if bytes, err := common.BuildResponse(-1, err.Error(), nil); err == nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write(bytes)
				return
			}
		} else {
			bytes, _ := common.BuildResponse(0, "success", o)
			w.Write(bytes)
			return
		}
	}
}
