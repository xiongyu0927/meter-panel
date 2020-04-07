package capacity

import (
	"encoding/json"
	"log"
	"strings"

	"github.com/Shopify/sarama"
	"github.com/spf13/viper"
)

func transferData(data CebCacpityReport) {
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
