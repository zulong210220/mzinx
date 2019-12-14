package config

import (
	"encoding/json"
	"io/ioutil"
	"mzinx/consts"
	"mzinx/ziface"
)

type Config struct {
	TcpServer ziface.IServer
	Host      string
	TcpPort   int32
	Name      string

	Version        string
	MaxConn        int
	MaxPackageSize uint32
	MsgChanSize    int

	WorkerPoolSize uint32
	TaskQueueSize  int
	IsWorker       bool
}

var (
	config *Config
)

func Reload() {
	cfn := "conf/mzinx.json"
	data, err := ioutil.ReadFile(cfn)
	if err != nil {
		panic(err)
	}

	// TODO 精细化
	err = json.Unmarshal(data, &config)
	if err != nil {
		panic(err)
	}
}

func init() {
	config = &Config{
		Host:    "0.0.0.0",
		TcpPort: consts.DefaultPort,
		Name:    consts.DefaultServerName,

		Version:        consts.DefaultVersion,
		MaxConn:        consts.DefaultMaxConn,
		MaxPackageSize: consts.DefaultMaxPackSize,
		MsgChanSize:    consts.DefaultMsgChanSize,

		WorkerPoolSize: consts.DefaultWorkerPoolSize,
		TaskQueueSize:  consts.DefaultTaskQueueSize,
		IsWorker:       true,
	}

	Reload()
}

func GetConfig() *Config {
	return config
}
