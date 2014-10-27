package notifier

import (
	"log"

	"github.com/influxdb/influxdb/client"
)

type InfluxdbNotifier struct {
	Host       string
	Username   string
	Password   string
	Database   string
	SeriesName string
}

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

func (influxdb *InfluxdbNotifier) toSeries(messages Messages) []*client.Series {

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
