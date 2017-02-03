package notifier

import (
	log "github.com/Sirupsen/logrus"
	elastic "gopkg.in/olivere/elastic.v3"
	"time"
)

type ElasticSearchNotifier struct {
	Host      string
	IndexName string
	NotifName string
}

type ElasticSearchMessage struct {
	Node    string `json:"node"`
	Service string `json:"service"`
	Checks  string `json:"checks"`
	Notes   string `json:"notes"`
	Output  string `json:"output"`
	Status  string `json:"status"`
	Time    int64  `json:"time"`
}

func (es *ElasticSearchNotifier) NotifierName() string {
	return es.NotifName
}

func (es *ElasticSearchNotifier) Notify(messages Messages) bool {

	client, err := elastic.NewClient(elastic.SetURL(es.Host), elastic.SetSniff(false))
	if err != nil {
		log.Println("unable to access elasticsearch. can't send notification. ", err)
		return false
	}

	err = es.setupIndex(client)
	if err != nil {
		log.Println("unable to setup elasticsearch index. ", err)
		return false
	}

	bulkIdxRequest, err := es.PrepareEsIndexRequest(client, messages)
	if err != nil {
		log.Println("unable to create bulk index request. ", err)
	}

	err = es.Send(client, bulkIdxRequest)
	if err != nil {
		log.Println("unable to send notifications. ", err)
	}

	log.Println("elasticsearch notification sent.")
	return true
}

func (es *ElasticSearchNotifier) setupIndex(client *elastic.Client) error {

	exists, err := client.IndexExists(es.IndexName).Do()
	if err != nil {
		return err
	}

	if !exists {
		createIndex, err := client.CreateIndex(es.IndexName).Do()
		if err != nil {
			return err
		}
		if !createIndex.Acknowledged {
			return err
		}
	}
	return nil
}

func (es *ElasticSearchNotifier) PrepareEsIndexRequest(client *elastic.Client, messages Messages) ([]elastic.BulkableRequest, error) {

	bulkIdxRequest := make([]elastic.BulkableRequest, 0)

	for _, message := range messages {
		esMessage := &ElasticSearchMessage{}

		esMessage.Node = message.Node
		esMessage.Service = message.Service
		esMessage.Checks = message.Check
		esMessage.Notes = message.Notes
		esMessage.Output = message.Output
		esMessage.Status = message.Status
		esMessage.Time = time.Now().Unix()

		idxRequest := elastic.
			NewBulkIndexRequest().
			Index(es.IndexName).
			Type("consul").
			Doc(esMessage)

		bulkIdxRequest = append(bulkIdxRequest, idxRequest)
	}
	return bulkIdxRequest, nil
}

func (es *ElasticSearchNotifier) Send(client *elastic.Client, bulkIdxRequest []elastic.BulkableRequest) error {
	bulkRequest := client.Bulk()

	for _, bulkIndex := range bulkIdxRequest {
		bulkRequest = bulkRequest.Add(bulkIndex)
	}

	_, err := bulkRequest.Do()
	if err != nil {
		return err
	}
	return nil
}
