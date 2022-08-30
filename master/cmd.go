package master

import "flag"

var ConfigPath string

func InitFlag() {
	// master -config ./master_config.json
	// master -h
	flag.StringVar(&ConfigPath, "config", "./master_config.json", "master config path")
	flag.Parse()
}
