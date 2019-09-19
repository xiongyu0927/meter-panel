package exporter

import (
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"

	"github.com/spf13/viper"
)

var cfg config

type config struct {
	PoolSeting        [][]string
	address           []string
	password          string
	refresh_frequence time.Duration
}

func localtest() config {
	var config config
	viper.SetConfigName("redis_pool")
	viper.AddConfigPath("./")
	err := viper.ReadInConfig()
	if err != nil {
		log.Println(err) // 读取配置文件失败致命错误
	}
	err = viper.Unmarshal(&config)
	if err != nil {
		log.Println(err)
	}

	config.address = []string{"10.0.129.2:26379", "10.0.128.28:26379", "10.0.128.97:26379"}
	config.password = "alauda_redis_passwd"
	frequence := viper.GetInt("MONNITOR_FREQUENCE")
	config.refresh_frequence = time.Duration(frequence) * time.Second
	return config
}

func LoadPoolConfig() config {
	var config config
	viper.SetConfigName("redis_pool")
	viper.AddConfigPath("/ace")
	err := viper.ReadInConfig()
	if err != nil {
		log.Println(err) // 读取配置文件失败致命错误
	}
	err = viper.Unmarshal(&config)
	if err != nil {
		log.Println(err)
	}

	frequence := viper.GetInt("MONNITOR_FREQUENCE")
	config.refresh_frequence = time.Duration(frequence) * time.Second
	config.address = strings.SplitN(os.Getenv("REDIS_HOST"), ",", -1)
	password, err := ioutil.ReadFile("/etc/paas_redis/password")
	if err != nil {
		log.Println(err)
	}
	config.password = strings.Replace(string(password), "\n", "", 1)
	return config
}
