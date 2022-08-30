package master

import (
	"net"
	"net/http"
	"strconv"
	"time"
)

type ApiServer struct {
	HttpServer http.Server
}

var ApiServ *ApiServer

func InitApiSever() error {
	mux := http.NewServeMux()
	mux.HandleFunc("/task", handleCreateTask)

	s := http.Server{
		Handler:      mux,
		ReadTimeout:  time.Duration(Cnf.ServerReadTimeout),
		WriteTimeout: time.Duration(Cnf.ServerWriteTimeout),
	}

	l, err := net.Listen("tcp", ":"+strconv.Itoa(Cnf.ServerPort))
	if err != nil {
		return err
	}

	ApiServ = &ApiServer{
		HttpServer: s,
	}

	go s.Serve(l)
	return nil

}

func handleCreateTask(w http.ResponseWriter, r *http.Request) {

}
