package main

import (
	"crontab_go/crontab/master"
	"flag"
	"fmt"
	"runtime"
	"time"
)

var(
	confFile string
)
func initArgs()  {
	flag.StringVar(&confFile,"config","./master.json","指定master.json")
	flag.Parse()
}
func initEnv()  {
	runtime.GOMAXPROCS(runtime.NumCPU())
}
func main()  {
	var (
		err error
	)
	//1.初始化线程
	initEnv()
	initArgs()
	//2.加载配置文件
	err = master.InitConfig(confFile)
	if err != nil {
		goto ERR
	}
	//3.加载etcd
	if err = master.InitJobMgr(); err != nil{
		goto ERR
	}
	//3.启动ApiHTTP服务
	if err = master.InitApiServer(); err != nil{
		goto ERR
	}

	for  {
		time.Sleep(10*time.Second)
	}
	return

ERR:
	fmt.Println(err)
}
