package config

import (
	"flag"
	"github.com/BurntSushi/toml"
	"io/ioutil"
	"log"
	"path/filepath"
)

var (
	configPath = flag.String("config", "./config.toml", "config path")
)

var cfg Config

type Config struct {
	RpcPort int `tomal:"rpc_port"`
}

func Init() {
	path, err := filepath.Abs(*configPath)
	if err != nil {
		log.Panicf("config file path err:%s", err)
		return
	}
	data, err := ioutil.ReadFile(path)
	if err == nil {
		err = toml.Unmarshal(data, &cfg)
	}
	if err != nil {
		log.Panicf("load toml file err:%s", err)
	}
}

func GetConfig() Config {
	return cfg
}
