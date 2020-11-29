package master

import (
	"crontab_go/crontab/common"
	json "github.com/json-iterator/go"
	"net"
	"net/http"
	"strconv"
	"time"
)

type ApiServer struct {
	httpServer *http.Server
}

var (
	G_apiServer *ApiServer
)

func InitApiServer() (err error) {
	var (
		mux        *http.ServeMux
		listener   net.Listener
		httpServer *http.Server
	)
	//配置路由
	mux = http.NewServeMux()
	mux.HandleFunc("/job/save", handleJobSave)     //新增任务
	mux.HandleFunc("/job/delete", handleJobDelete) //删除任务
	mux.HandleFunc("/job/list", handleJobList)     //查询任务列表
	//启动TCP监听
	if listener, err = net.Listen("tcp", ":"+strconv.Itoa(G_config.ApiPort)); err != nil {
		return err
	}
	//创建一个http服务
	httpServer = &http.Server{
		Handler:           mux,
		ReadTimeout:       time.Duration(G_config.ApiReadTimeout) * time.Millisecond,
		ReadHeaderTimeout: time.Duration(G_config.ApiReadTimeout) * time.Millisecond,
	}
	G_apiServer = &ApiServer{httpServer: httpServer}
	//异步拉起服务
	go httpServer.Serve(listener)
	//go httpServer.ListenAndServe()
	return
}

//保存定时任务
func handleJobSave(resp http.ResponseWriter, req *http.Request) {
	var (
		err     error
		postJob string
		job     common.Job
		oldJob  *common.Job
	)
	//任务保存到etcd job={"name":"job1","command":"echo hello","cronExpr":""}
	if err = req.ParseForm(); err != nil {
		goto ERR
	}
	postJob = req.PostForm.Get("job")
	//反序列化job
	if err = json.Unmarshal([]byte(postJob), &job); err != nil {
		goto ERR
	}
	//保存
	if oldJob, err = G_jobMgr.SaveJob(&job); err != nil {
		goto ERR
	}
	if bytes, err := common.BuildResponse(0, "success", oldJob); err == nil {
		resp.Write(bytes)
	}
	return
ERR:
	if bytes, err := common.BuildResponse(-1, err.Error(), nil); err == nil {
		resp.Write(bytes)
	}
}

//删除任务
func handleJobDelete(w http.ResponseWriter, r *http.Request) {
	var (
		err     error
		prevJob *common.Job
		bytes   []byte
		names   string
	)
	if err = r.ParseForm(); err != nil {
		//goto ERR
		goto ERR
	}

	names = r.PostForm.Get("name")
	prevJob, err = G_jobMgr.DelJob(names)
	if err != nil {
		goto ERR
	}
	if bytes, err = common.BuildResponse(200, "success", prevJob); err == nil {
		w.Write(bytes)
	}
	return
ERR:
	if bytes, err := common.BuildResponse(-1, err.Error(), nil); err == nil {
		w.Write(bytes)
	}
}

func handleJobList(w http.ResponseWriter, r *http.Request) {
	var (
		err   error
		bytes []byte
		list  []*common.Job
	)
	list, err = G_jobMgr.JobList()
	if err != nil {
		goto ERR
	}
	if bytes, err = common.BuildResponse(200, "success", list); err == nil {
		w.Write(bytes)
	}
	return
ERR:
	if bytes, err = common.BuildResponse(-1, err.Error(), nil); err == nil {
		w.Write(bytes)
	}
}
