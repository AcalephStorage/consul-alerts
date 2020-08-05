package notifier

import (
	_ "github.com/influxdata/influxdb1-client"
	client "github.com/influxdata/influxdb1-client/v2"
	log "github.com/sirupsen/logrus"
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
	c, err := client.NewHTTPClient(client.HTTPConfig{
		Addr:     influxdb.Host,
		Username: influxdb.Username,
		Password: influxdb.Password,
	})

	defer c.Close()

	if err != nil {
		log.Println("invalid configure. ", err)
		return false
	}

	bp, err := client.NewBatchPoints(client.BatchPointsConfig{
		Database:  influxdb.Database,
		Precision: "s",
	})
	if err != nil {
		log.Println("invalid precision. ", err)
		return false
	}

	for _, m := range messages {
		pt, err := client.NewPoint(influxdb.SeriesName, nil, map[string]interface{}{
			"node":    m.Node,
			"service": m.Service,
			"checks":  m.Check,
			"notes":   m.Notes,
			"output":  m.Output,
			"status":  m.Status,
		}, m.Timestamp)
		if err != nil {
			log.Println("unable to send notifications: ", err)
			return false
		}

		bp.AddPoint(pt)
	}

	if err = c.Write(bp); err != nil {
		log.Println("unable to send notifications: ", err)
		return false
	}

	log.Println("influxdb notification sent.")
	return true
}
