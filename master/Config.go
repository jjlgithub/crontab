package master

import (
	"encoding/json"
	"os"
)

type Config struct {
	ApiPort int `json:"apiPort"`
	ApiReadTimeout int `json:"apiReadTimeout"`
	ApiWriteTimeout int `json:"apiWriteTimeout"`
	DialTimeout int `json:"dialTimeout"`
	EtcdEndpoints []string `json:"etcdEndpoints"`
}

var (
	G_config *Config
)

func InitConfig(filename string) (err error) {
	//打开文件
	var(
		file *os.File
	)
	file, err = os.Open(filename)
	if err != nil {
		return err
	}
	var config Config
	err = json.NewDecoder(file).Decode(&config)
	if err != nil {
		return err
	}

	G_config = &config
	return
}

