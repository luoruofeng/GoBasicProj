package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	ServerPort         int      `json:"server_port"`
	ServerReadTimeout  int      `json:"server_read_timeout"`
	ServerWriteTimeout int      `json:"server_write_timeout"`
	ServerWebRoot      string   `json:"server_web_root"`
	EtcdAddrs          []string `json:"etcd_addrs"`
	EtcdDialTimeout    int      `json:"etcd_dial_timeout"`
	MongoUri           string   `json:"mongo_uri"`
	MongoDialTimeout   int      `json:"mongo_dial_timeout"`
}

var (
	Cnf Config
)

func InitConfig(configPath string) error {
	var c Config
	b, err := os.ReadFile(configPath)
	if err != nil {
		return err
	}
	err = json.Unmarshal(b, &c)
	if err != nil {
		return err
	}

	Cnf = c
	return nil
}
