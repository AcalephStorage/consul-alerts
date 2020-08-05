package main

import (
	"fmt"

	"net/http"

	log "github.com/sirupsen/logrus"
)

func healthHandler(w http.ResponseWriter, r *http.Request) {

	node := r.URL.Query().Get("node")
	service := r.URL.Query().Get("service")
	check := r.URL.Query().Get("check")

	log.Println(node, service, check)

	status, output := consulClient.CheckStatus(node, service, check)

	var code int
	switch status {
	case "passing":
		code = 200
	case "warning", "critical":
		code = 503
	default:
		status = "unknown"
		code = 404
	}

	log.Printf("health status check result for node=%s,service=%s,check=%s: %d", node, service, check, code)

	var result string
	if output == "" {
		result = ""
	} else {
		result = fmt.Sprintf("output: %s\n", output)
	}
	body := fmt.Sprintf("status: %s\n%s", status, result)
	w.WriteHeader(code)
	w.Write([]byte(body))
}
