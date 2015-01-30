package main

import (
	"io"
	"os"
	"syscall"

	"encoding/json"
	"io/ioutil"
	"os/exec"

	log "github.com/AcalephStorage/consul-alerts/Godeps/_workspace/src/github.com/Sirupsen/logrus"
)

func runWatcher(address, datacenter, watchType string) {
	consulAlert := os.Args[0]
	cmd := exec.Command(
		"consul", "watch",
		"-http-addr", address,
		"-datacenter", datacenter,
		"-type", watchType,
		consulAlert, "watch", watchType)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		exitError, _ := err.(*exec.ExitError)
		exitCode := 1
		if exitError == nil {
			log.Println("Shutting down watcher --> ", err.Error())
		} else {
			status, _ := exitError.Sys().(syscall.WaitStatus)
			exitCode := status.ExitStatus()
			log.Println("Shutting down watcher --> Exit Code: ", exitCode)
		}
		os.Exit(exitCode)
	} else {
		log.Printf("Execution complete.")
	}
}

func toWatchObject(reader io.Reader, v interface{}) {
	data, err := ioutil.ReadAll(reader)
	if err != nil {
		log.Println("stdin read error: ", err)
		// todo: what to do when can't read?
	}
	err = json.Unmarshal(data, v)
	if err != nil {
		log.Println("json unmarshall error: ", err)
		// todo: what if we can't serialise?
	}
}
