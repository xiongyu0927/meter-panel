package capacity

import (
	"encoding/json"
	"log"
	"strings"
	"time"

	"github.com/Shopify/sarama"
	"github.com/robfig/cron"
	"github.com/spf13/viper"
)

const (
	A = "上地机房"
	B = "酒仙桥机房"
	S = "单中心"
)

var (
	siteInfo string
)

// CebCapacityReport 容量报文结构体
type CebCapacityReport struct {
	// key is site
	Kind        string      `json:"kind"`
	Environment string      `json:"environment"`
	DataCenter  datacenter  `json:"dataCenter"`
	Data        interface{} `json:"data"`
	Date        string      `json:"date"`
}

type datacenter struct {
	Description string `json:"description"`
	Name        string `json:"name"`
}

// NewCebCapacityReport 创建通用报文实体
func NewCebCapacityReport(kind string) *CebCapacityReport {
	var description, environment string
	env := viper.GetString("ENVIRONMENT_INFO")
	switch env {
	case "as":
		description = A
		environment = "stage"
	case "ap":
		description = A
		environment = "production"
	case "bs":
		description = B
		environment = "stage"
	case "bp":
		description = B
		environment = "production"
	case "s":
		description = S
		environment = "development"
	}

	return &CebCapacityReport{
		Kind:        kind,
		Environment: environment,
		Date:        getDateYMD(),
		DataCenter: datacenter{
			Description: description,
			Name:        env,
		},
	}
}

func (r *CebCapacityReport) isYourSiteByClusterName(cluster string) bool {
	if r.DataCenter.Name == "s" {
		return true
	}

	slice := strings.SplitN(cluster, "-", -1)
	for _, v := range slice {
		if v == r.DataCenter.Name {
			return true
		}
	}
	return false
}

//Start entrance
func Start() {
	c := cron.New()
	spec := "0 0 18 * * *"
	c.AddFunc(spec, func() {
		clusterScalReportStart()
	})

	c.Start()
}

func transferData(data CebCapacityReport) {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Partitioner = sarama.NewRandomPartitioner
	config.Producer.Return.Successes = true
	config.Producer.Return.Errors = true

	host := strings.SplitN(viper.GetString("KAFKA_ENDPOINT"), ",", -1)
	topic := viper.GetString("KAFKA_TOPIC")
	if topic == "" {
		log.Println("didn't set kafka topic")
	}

	producer, e := sarama.NewAsyncProducer(host, config)
	if e != nil {
		log.Println(e)
		return
	}

	defer producer.AsyncClose()

	go func(p sarama.AsyncProducer) {
		for {
			select {
			case <-p.Successes():
				//fmt.Println("offset: ", suc.Offset, "timestamp: ", suc.Timestamp.String(), "partitions: ", suc.Partition)
			case fail := <-p.Errors():
				log.Println(fail.Err)
			}
		}
	}(producer)

	data2, _ := json.Marshal(data)
	reportData := sarama.ByteEncoder(data2)
	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: reportData,
	}
	producer.Input() <- msg
}

func getDateYMD() string {
	return time.Now().Format("2006-01-02")
}
