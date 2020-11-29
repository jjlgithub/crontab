package common

import json "github.com/json-iterator/go"

//定时任务
type Job struct {
	Name string `json:"name"`
	Command string `json:"command"` //shell命令
	CronExpr string `json:"cronExpr"` //cron表达式
}


type Response struct {
	Errno int `json:"errno"`
	Msg string `json:"msg"`
	Data interface{} `json:"data"`
}

func BuildResponse(errno int,msg string,data interface{}) (resp []byte,err error) {
	var(
		response Response
	)
	response.Errno = errno
	response.Msg = msg
	response.Data = data

	resp , err = json.Marshal(response)
	return
}