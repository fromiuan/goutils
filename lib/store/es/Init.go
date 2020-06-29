package es

import (
	"context"
	"fmt"
	"net/url"
	"strconv"

	els "github.com/olivere/elastic"
	"github.com/olivere/elastic/config"
)

type ESConfig struct {
	Endpoint     string
	Index        string
	Shards       int
	Replicats    int
	DefaultTable string
	UserName     string
	Passwd       string
}

var EsClient *els.Client
var defaultType string
var defaultIndex string

func Init(conf *ESConfig) error {
	if EsClient != nil {
		return nil
	}
	cfg := config.Config{
		URL:      "http://127.0.0.1:9200",
		Index:    "",
		Shards:   5,
		Replicas: 1,
	}
	if len(conf.Index) == 0 {
		return fmt.Errorf("es: The index is empty")
	}
	cfg.Index = conf.Index

	if len(conf.Endpoint) > 0 {
		uri, err := url.Parse(conf.Endpoint)
		if err != nil {
			return err
		}
		if port, _ := strconv.Atoi(uri.Port()); port == 0 {
			return fmt.Errorf("es: Not found endpoint port")
		}
		cfg.URL = conf.Endpoint
	}

	if conf.Replicats > 0 {
		cfg.Replicas = conf.Replicats
	}

	if conf.Shards > 0 {
		cfg.Shards = conf.Shards
	}

	if len(conf.UserName) > 0 {
		cfg.Username = conf.UserName
		cfg.Password = conf.Passwd
	}
	if len(conf.DefaultTable) == 0 {
		return fmt.Errorf("es: Not setting default document")
	}

	ctx := context.Background()
	client, err := els.NewClientFromConfig(&cfg)
	if err != nil {
		return err
	}

	if !client.IsRunning() {
		return fmt.Errorf("es: The server not running")
	}

	isFound, err := client.IndexExists(cfg.Index).Do(ctx)
	if err != nil {
		return err
	}

	if !isFound {
		return fmt.Errorf("es: Not found index %v", cfg.Index)
	}

	mapping, err := client.GetMapping().Index(cfg.Index).Do(ctx)
	if err != nil {
		return err
	}
	if len(mapping) == 0 {
		return fmt.Errorf("es: The index mapping is empty | %v", cfg.Index)
	}

	if ok, err := client.TypeExists().Index(cfg.Index).Type(conf.DefaultTable).Do(ctx); err != nil || !ok {
		return fmt.Errorf("es: Not found default table ")
	}
	defaultType = conf.DefaultTable
	defaultIndex = conf.Index
	EsClient = client
	return nil
}

func GetDefalutTable() string {
	return defaultType
}

func GetDefaultIndex() string {
	return defaultIndex
}
func POST(id int, data interface{}) (*els.IndexResponse, error) {
	if EsClient == nil {
		return nil, fmt.Errorf("es: The Es client not initializatioçn")
	}
	ctx := context.Background()
	client := EsClient.Index()
	if id > 0 {
		client.Id(strconv.Itoa(id))
	}
	client.Index(defaultIndex)
	switch data.(type) {
	case string:
		client.BodyString(data.(string))
	case []byte:
		client.BodyString(string(data.([]byte)))
	default:
		client.BodyJson(data)
	}

	return client.Type(defaultType).Do(ctx)
}

func PBulk(data ...interface{}) (*els.BulkResponse, error) {
	if EsClient == nil {
		return nil, fmt.Errorf("es: The Es client not initializatioçn")
	}
	ctx := context.Background()
	esBulk := EsClient.Bulk()
	for _, item := range data {
		doc := els.NewBulkIndexRequest().Index(defaultType).Type(defaultType).Doc(item)
		esBulk.Type(defaultType).Index(defaultIndex).Add(doc)
	}
	return esBulk.Do(ctx)
}
