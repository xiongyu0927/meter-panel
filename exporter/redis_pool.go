package exporter

import (
	"log"
	"meter-panel/configs"
	"meter-panel/store"
	"strconv"
	"strings"
	"time"

	"github.com/go-redis/redis"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/spf13/viper"
	"k8s.io/apimachinery/pkg/labels"
)

// Counter metric for request count
var RedisPoolConnectCount = prometheus.NewGaugeVec(
	prometheus.GaugeOpts{
		Name: "pod_redis_pool_count",
		Help: "pod_redis_pool_count",
	},
	[]string{"pod_name", "pool_type", "data_type"},
)

func init() {
	style := viper.GetString("PROJECT")
	if style != "CEB" {
		return
	}
	cfg = LoadPoolConfig()
	// cfg = localtest()
	prometheus.MustRegister(RedisPoolConnectCount)
	Data()
}

func Data() {
	var opt = &redis.ClusterOptions{
		Addrs:         cfg.address,
		Password:      cfg.password,
		ReadOnly:      true,
		RouteRandomly: true,
		PoolSize:      10,
		MinIdleConns:  1,
		IdleTimeout:   time.Duration(30) * time.Second,
	}
	c := redis.NewClusterClient(opt)
	cmd := redis.NewStringCmd("client", "list")
	go func() {
		for {
			c.Process(cmd)
			a, err := cmd.Result()
			if err != nil {
				log.Println(err)
			}
			data := handleData(a)
			for k, v := range data {
				if v[3] == 0 {
					RedisPoolConnectCount.WithLabelValues(k, "all", "Current").Set(float64(v[0]))
					RedisPoolConnectCount.WithLabelValues(k, "all", "Total").Set(float64(v[2]))
				} else {
					RedisPoolConnectCount.WithLabelValues(k, "Read", "Current").Set(float64(v[0]))
					RedisPoolConnectCount.WithLabelValues(k, "Read", "Total").Set(float64(v[2]))
					RedisPoolConnectCount.WithLabelValues(k, "Write", "Current").Set(float64(v[1]))
					RedisPoolConnectCount.WithLabelValues(k, "Write", "Total").Set(float64(v[3]))
				}
			}
			time.Sleep(cfg.refresh_frequence)
		}
	}()
}

func handleData(a string) map[string][]int {
	var thistime = make(map[string][]int)
	b := strings.SplitN(a, "\n", -1)
	for _, v := range b[:len(b)-1] {
		tmp := strings.SplitN(v, " ", -1)
		ip := awk(tmp[1], "=", ":")
		cmd := awk(tmp[17], "=", "")

		if v, ok := thistime[ip]; ok {
			if IsReadPool(cmd) {
				v[0]++
			} else {
				v[1]++
			}
			continue
		}

		subthistime := []int{1, 1}
		thistime[ip] = subthistime
	}
	return paddingData(thistime)
}

func paddingData(thistime map[string][]int) map[string][]int {
	var tmp = make(map[string][]int)
	var subtmp = make([]int, 4)
	lister := store.AllLister.PodLister[configs.GlobalName]
	for _, v := range cfg.PoolSeting {
		labelset := labels.Set(map[string]string{v[0]: v[1]}).AsSelector()
		pl, err := lister.List(labelset)
		if err != nil {
			log.Println(err)
		}
		for _, v1 := range pl {

			if cn, ok := thistime[v1.Status.PodIP]; ok {
				v2, _ := strconv.Atoi(v[2])
				v3, _ := strconv.Atoi(v[3])
				if v[3] != "" {
					tmp[v1.ObjectMeta.Name] = cn
					tmp[v1.ObjectMeta.Name] = append(tmp[v1.ObjectMeta.Name], v2, v3)
				} else {
					subtmp[0] = cn[0] + cn[1]
					subtmp[2] = v2
					subtmp[3] = v3
					tmp[v1.ObjectMeta.Name] = subtmp
				}
			}
		}
	}
	return tmp
}

func IsReadPool(cmd string) bool {
	cmd = strings.ToLower(cmd)
	if strings.Contains(cmd, "get") || strings.Contains(cmd, "scan") || strings.Contains(cmd, "exits") {
		return true
	}

	if cmd == "smembers" || cmd == "ttl" || cmd == "subscribe" {
		return true
	}

	return false
}

func awk(source string, start string, end string) string {
	s := strings.Index(source, start) + 1
	if end == "" {
		return string([]rune(source)[s:])
	}
	e := strings.Index(source, end)
	return string([]rune(source)[s:e])
}
