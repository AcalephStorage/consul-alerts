package notifier

import (
	log "github.com/uchiru/consul-alerts/Godeps/_workspace/src/github.com/Sirupsen/logrus"
	"github.com/uchiru/consul-alerts/Godeps/_workspace/src/github.com/influxdb/influxdb/client"
)

type InfluxdbNotifier struct {
	Enabled    bool
	Host       string `json:"host"`
	Username   string `json:"username"`
	Password   string `json:"password"`
	Database   string `json:"database"`
	SeriesName string `json:"series-name"`
}

// NotifierName provides name for notifier selection
func (influxdb *InfluxdbNotifier) NotifierName() string {
	return "influxdb"
}

func (influxdb *InfluxdbNotifier) Copy() Notifier {
	notifier := *influxdb
	return &notifier
}

//Notify sends messages to the endpoint notifier
func (influxdb *InfluxdbNotifier) Notify(messages Messages) bool {
	config := &client.ClientConfig{
		Host:     influxdb.Host,
		Username: influxdb.Username,
		Password: influxdb.Password,
		Database: influxdb.Database,
	}

	influxdbClient, err := client.New(config)
	if err != nil {
		log.Println("unable to access influxdb. can't send notification. ", err)
		return false
	}

	seriesList := influxdb.toSeries(messages)
	err = influxdbClient.WriteSeries(seriesList)

	if err != nil {
		log.Println("unable to send notifications: ", err)
		return false
	}

	log.Println("influxdb notification sent.")
	return true
}

func (influxdb InfluxdbNotifier) toSeries(messages Messages) []*client.Series {

	seriesName := influxdb.SeriesName
	columns := []string{
		"node",
		"service",
		"checks",
		"notes",
		"output",
		"status",
	}

	seriesList := make([]*client.Series, len(messages))
	for index, message := range messages {

		point := []interface{}{
			message.Node,
			message.Service,
			message.Check,
			message.Notes,
			message.Output,
			message.Status,
		}

		series := &client.Series{
			Name:    seriesName,
			Columns: columns,
			Points:  [][]interface{}{point},
		}
		seriesList[index] = series
	}
	return seriesList
}
