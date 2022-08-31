package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/luoruofeng/gobasicproj/master"
	"github.com/luoruofeng/gobasicproj/master/task_srv/etcd"
)

func main() {
	master.InitFlag()

	err := master.InitConfig(master.ConfigPath)
	if err != nil {
		panicAndExit(err)
	}

	err = etcd.InitEtcdTaskSrv()
	if err != nil {
		panicAndExit(err)
	}

	err = master.InitApiSever()
	if err != nil {
		panicAndExit(err)
	} else {
		log.Println("API server is running...")
	}

	// Setup our Ctrl+C handler
	SetupCloseHandler()

}

func SetupCloseHandler() {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c
	log.Println("\r- Ctrl+C pressed in Terminal")
	DeleteFiles() //process exit to do something
	os.Exit(0)
}

func DeleteFiles() {
	log.Println("- Run Clean Up - Delete some File")
	// _ = os.Remove(something)
	log.Println("- Good bye!")
}

func panicAndExit(err error) {
	DeleteFiles()
	log.Fatal(err)
	os.Exit(1)
}
