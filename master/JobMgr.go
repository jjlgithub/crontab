package master

import (
	"context"
	"crontab_go/crontab/common"
	json "github.com/json-iterator/go"
	"go.etcd.io/etcd/clientv3"
	"time"
)

type JobMgr struct {
	client *clientv3.Client
	kv     clientv3.KV
	lease  clientv3.Lease
}

var (
	G_jobMgr *JobMgr
)

func InitJobMgr() (err error) {
	var (
		config clientv3.Config
		client *clientv3.Client
		kv     clientv3.KV
		lease  clientv3.Lease
	)
	config = clientv3.Config{
		Endpoints:   G_config.EtcdEndpoints,
		DialTimeout: time.Duration(G_config.DialTimeout) * time.Microsecond,
	}
	//简历etcd连接
	if client, err = clientv3.New(config); err != nil {
		return
	}
	//得到kv和lease字集
	kv = clientv3.NewKV(client)
	lease = clientv3.NewLease(client)
	G_jobMgr = &JobMgr{
		client: client,
		kv:     kv,
		lease:  lease,
	}
	return
}

//add task
func (JobMgr *JobMgr) SaveJob(job *common.Job) (oldJob *common.Job, err error) {
	var (
		jobKey   string
		jobValue []byte
		putResp  *clientv3.PutResponse
	)

	jobKey = common.SAVE_DIR_NAME + job.Name
	//任务的值
	if jobValue, err = json.Marshal(job); err != nil {
		return nil, err
	}
	if putResp, err = JobMgr.kv.Put(context.TODO(), jobKey, string(jobValue), clientv3.WithPrevKV()); err != nil {
		return
	}
	//如果是更新返回旧址
	if putResp.PrevKv != nil {
		if err = json.Unmarshal(putResp.PrevKv.Value, &oldJob); err != nil {
			err = nil
			return nil, err
		}
	}
	return
}

//删除任务 del task
func (jobMgr *JobMgr) DelJob(name string) (oldJob *common.Job, err error) {
	var (
		del *clientv3.DeleteResponse
	)
	jobKey := common.SAVE_DIR_NAME + name
	if del, err = jobMgr.kv.Delete(context.TODO(), jobKey, clientv3.WithPrevKV()); err != nil {
		return
	}
	if len(del.PrevKvs) > 0 {
		err := json.Unmarshal(del.PrevKvs[0].Value, &oldJob)
		if err != nil {
			err = nil
			return nil, err
		}
	}
	return
}

//select task
func (j JobMgr) JobList() (jobList []*common.Job,err error) {
	var (
		lists *clientv3.GetResponse
		job *common.Job
	)
	key := common.SAVE_DIR_NAME
	if lists, err = j.kv.Get(context.TODO(),key,clientv3.WithPrefix()); err != nil{
		return
	}
	jobList = make([]*common.Job,0)
	for _, item := range lists.Kvs {
		job = &common.Job{}
		if err = json.Unmarshal(item.Value,job);err != nil{
			err = nil
			continue
		}
		jobList = append(jobList,job)
	}
	return
}