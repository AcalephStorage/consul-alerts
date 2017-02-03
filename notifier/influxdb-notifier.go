package notifier

import (
	"time"

	log "github.com/AcalephStorage/consul-alerts/Godeps/_workspace/src/github.com/Sirupsen/logrus"
	"github.com/AcalephStorage/consul-alerts/Godeps/_workspace/src/github.com/influxdb/influxdb/client/v2"
)

type InfluxdbNotifier struct {
	Host       string
	Username   string
	Password   string
	Database   string
	SeriesName string
	NotifName  string
}

// NotifierName provides name for notifier selection
func (influxdb *InfluxdbNotifier) NotifierName() string {
	return influxdb.NotifName
}

func (influxdb *InfluxdbNotifier) Notify(messages Messages) bool {

	influxdbClient, err := client.NewHTTPClient(client.HTTPConfig{
		Addr:     influxdb.Host,
		Username: influxdb.Username,
		Password: influxdb.Password,
	})

	if err != nil {
		log.Println("unable to access influxdb. can't send notification. ", err)
		return false
	}

	bp := influxdb.toSeries(messages)
	err = influxdbClient.Write(bp)

	if err != nil {
		log.Println("unable to send notifications: ", err)
		return false
	}

	log.Println("influxdb notification sent.")
	return true
}

func (influxdb *InfluxdbNotifier) toSeries(messages Messages) client.BatchPoints {

	seriesName := influxdb.SeriesName
	bp, _ := client.NewBatchPoints(client.BatchPointsConfig{
		Database:  influxdb.Database,
		Precision: "s",
	})

	for _, message := range messages {
		// Create a point and add to batch
		tags := map[string]string{"node": message.Node}
		fields := map[string]interface{}{
			"node":    message.Node,
			"service": message.Service,
			"check":   message.Check,
			"notes":   message.Notes,
			"output":  message.Output,
			"status":  message.Status,
		}
		pt, _ := client.NewPoint(seriesName, tags, fields, time.Now())
		bp.AddPoint(pt)
	}
	return bp
}
