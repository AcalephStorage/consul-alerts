package main

import (
	"encoding/json"
	log "github.com/uchiru/consul-alerts/Godeps/_workspace/src/github.com/Sirupsen/logrus"
	"net/http"
	"strings"
)

func healthWildcardHandler(w http.ResponseWriter, r *http.Request) {

	node := r.URL.Query().Get("node")
	service := r.URL.Query().Get("service")
	check := r.URL.Query().Get("check")
	status := r.URL.Query().Get("status")
	alwaysOk := r.URL.Query().Get("alwaysOk") != "" // Always return 200 code, even if failures in data
	ignoreBlacklist := r.URL.Query().Get("ignoreBlacklist") != ""

	var statuses []string
	if status != "" {
		statuses = strings.Split(status, ",")
	}
	log.Printf("Query: node: %v, service: %v, check: %v, status: %v, alwaysOk: %v, ignoreBlacklist: %v", node, service, check, status, alwaysOk, ignoreBlacklist)

	alerts := consulClient.NewAlertsWithFilter(node, service, check, statuses, ignoreBlacklist)

	code := 200

	if !alwaysOk {
		var newCode int
		for _, alert := range alerts {
			switch alert.Status {
			case "passing":
				newCode = 200
			case "warning", "critical":
				newCode = 503
			default:
				status = "unknown"
				newCode = 404
			}
			if newCode > code {
				code = newCode
			}
		}
	}

	body, _ := json.Marshal(alerts)
	w.WriteHeader(code)
	w.Write([]byte(body))
}
