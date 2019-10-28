package core

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"meter-panel/configs"
	"strings"
	"time"

	"github.com/robfig/cron"
	"github.com/spf13/viper"

	elastic "gopkg.in/olivere/elastic.v5"
)

type EsClient struct {
	Client   *elastic.Client
	Index    []string
	Clusters []string
	// k is cluster name
	// {cluster_name {index [all, warning]}}
	Data map[string]map[string][]int64
}

func (its *EsClient) SearchCount(cluster string, index string) error {
	var result *elastic.SearchResult
	var err error
	boolQuery := elastic.NewBoolQuery()
	boolQuery.Must(elastic.NewTermQuery("detail.cluster_name.keyword", cluster))
	query := elastic.NewConstantScoreQuery(boolQuery)
	result, err = its.Client.Search().Index(index).Type("event").Query(query).Size(0).Do(context.Background())
	if err != nil {
		return err
	}
	if _, ok := its.Data[cluster]; !ok {
		its.Data[cluster] = make(map[string][]int64)
	}
	tmp := make([]int64, 2)
	tmp[0] = result.Hits.TotalHits

	boolQuery.Must(elastic.NewTermQuery("log_level", "1"))
	query = elastic.NewConstantScoreQuery(boolQuery)
	result, err = its.Client.Search().Index(index).Type("event").Query(query).Size(0).Do(context.Background())
	if err != nil {
		return err
	}
	tmp[1] = result.Hits.TotalHits
	its.Data[cluster][index] = tmp
	return nil
}

func NewEsClient(clusters configs.AllK8SConfigs) (*EsClient, error) {
	var tmp []string
	for k, _ := range clusters {
		tmp = append(tmp, k)
	}
	client, err := GenerateESclient()
	if err != nil {
		return nil, err
	}
	var x = &EsClient{
		Client:   client,
		Clusters: tmp,
		Data:     make(map[string]map[string][]int64),
	}
	err = x.GenerateAllIndex()
	if err != nil {
		return nil, err
	}
	x.initDate()
	return x, nil
}

func (its *EsClient) Loop() {
	c := cron.New()
	spec := "* 1 0 * * *"
	c.AddFunc(spec, func() {
		its.GenerateAllIndex()
		its.initDate()
		log.Println(its.Index)
	})
	c.Start()
}

func (its *EsClient) initDate() {
	for _, v := range its.Clusters {
		for _, v1 := range its.Index {
			err := its.SearchCount(v, v1)
			if err != nil {
				log.Println(err)
			}
		}
	}
}

func GenerateESclient() (*elastic.Client, error) {
	Host := strings.SplitN(viper.GetString("ES_HOST"), ",", -1)
	un, err := ioutil.ReadFile("/etc/pass_es/username")
	if err != nil {
		return nil, err
	}
	pwd, err := ioutil.ReadFile("/etc/pass_es/password")
	if err != nil {
		return nil, err
	}
	UserName := strings.Replace(string(un), "\n", "", 1)
	PassWord := strings.Replace(string(pwd), "\n", "", 1)
	log.Println(UserName, PassWord)
	client, err := elastic.NewClient(elastic.SetURL(Host...), elastic.SetBasicAuth(UserName, PassWord),
		elastic.SetSniff(false), elastic.SetHealthcheck(false))
	if err != nil {
		return nil, err
	}
	return client, nil
}

func (its *EsClient) GenerateAllIndex() error {
	num := viper.GetInt("TTL")
	var tmp []string
	today := time.Now().AddDate(0, 0, 1)
	for i := 0; i < num; i++ {
		today = today.AddDate(0, 0, -1)
		riqi := today.Format("20060102")
		target := "event-" + riqi
		tmp = append(tmp, target)
	}
	its.Index = tmp
	err := its.IndexExists()
	if err != nil {
		return err
	}
	return nil
}

func (its *EsClient) IndexExists() error {
	var tmp []string
	for _, v := range its.Index {
		exists, err := its.Client.IndexExists(v).Do(context.Background())
		if err != nil {
			return err
		}
		if !exists {
			tmp = append(tmp, v)
		}
	}
	if tmp != nil {
		log.Println(tmp)
		a := fmt.Sprintf("these index %v doesn't exit", tmp)
		return errors.New(a)
	}
	return nil
}
