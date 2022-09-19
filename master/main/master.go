package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/luoruofeng/gobasicproj/master"
	c "github.com/luoruofeng/gobasicproj/master/config"
	"github.com/luoruofeng/gobasicproj/master/task_srv/etcd"
)

func main() {
	// 初始化flag
	master.InitFlag()

	err := c.InitConfig(master.ConfigPath)
	if err != nil {
		panicAndExit(err)
	}

	//初始化 etcd
	err = etcd.InitEtcdTaskSrv()
	if err != nil {
		panicAndExit(err)
	}

	// 初始化Server
	err = master.InitApiSever()
	if err != nil {
		panicAndExit(err)
	} else {
		log.Println("API server is running...", c.Cnf.ServerPort)
	}

	// Setup our Ctrl+C handler
	SetupCloseHandler()

}

//优雅退出
func SetupCloseHandler() {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c
	log.Println("\r- Ctrl+C pressed in Terminal")
	CloseAppHandle() //process exit to do something
	os.Exit(0)
}

func CloseAppHandle() {
	log.Println("- Run Clean Up - Delete some File")
	// _ = os.Remove(something)
	etcd.EtcdTaskSrv.CloseTaskSrv()
	log.Println("- Good bye!")
}

func panicAndExit(err error) {
	CloseAppHandle()
	log.Fatal(err)
	os.Exit(1)
}
