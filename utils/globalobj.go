package utils

import (
	"encoding/json"
	"io/ioutil"
	"mzinx/consts"
	"mzinx/ziface"
)

type GlobalObj struct {
	TcpServer ziface.IServer
	Host      string
	TcpPort   int32
	Name      string

	Version        string
	MaxConn        int32
	MaxPackageSize uint32
}

var (
	GlobalObject *GlobalObj
)

func (g *GlobalObj) Reload() {
	cfn := "conf/mzinx.json"
	data, err := ioutil.ReadFile(cfn)
	if err != nil {
		panic(err)
	}

	// TODO 精细化
	err = json.Unmarshal(data, &GlobalObject)
	if err != nil {
		panic(err)
	}
}

func init() {
	GlobalObject = &GlobalObj{
		Host:    "0.0.0.0",
		TcpPort: 8999,
		Name:    "MZinxServerApp",

		Version:        "v04",
		MaxConn:        consts.DefaultMaxConn,
		MaxPackageSize: consts.DefaultMaxPackSize,
	}

	GlobalObject.Reload()
}
