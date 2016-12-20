package consul

import (
	"encoding/json"
	"reflect"
	"strconv"
	"testing"

	consulapi "github.com/AcalephStorage/consul-alerts/Godeps/_workspace/src/github.com/hashicorp/consul/api"
)

func testClient() (*ConsulAlertClient, error) {
	return NewClient("localhost:8500", "dc1", "")
}

func clearKVPath(t *testing.T, c *ConsulAlertClient, path string) {
	_, err := c.api.KV().DeleteTree(path, nil)
	if err != nil {
		t.Error(err.Error())
	}
}

func TestLoadCustomValueForString(t *testing.T) {
	var strVar string
	input := "test-data"
	data := []byte(input)
	loadCustomValue(&strVar, data, ConfigTypeString)
	if strVar != "test-data" {
		t.Errorf("unable to parse %s to string", input)
	}
}

func TestLoadCustomValueForBool(t *testing.T) {
	var boolVar bool
	input := []string{
		"true",
		"false",
		"True",
		"False",
		"TRUE",
		"FALSE",
	}

	for i, in := range input {
		data := []byte(in)
		loadCustomValue(&boolVar, data, ConfigTypeBool)
		if i%2 == 0 && !boolVar {
			t.Errorf("unable to parse %s to boolean", in)
		}
	}

}

func TestLoadCustomValueForInt(t *testing.T) {
	var intVar int
	input := "235"
	data := []byte(input)
	loadCustomValue(&intVar, data, ConfigTypeInt)
	if in, _ := strconv.Atoi(input); in != intVar {
		t.Errorf("unable to parse %s to int", input)
	}
}

func TestGetProfileForEntity(t *testing.T) {
	client, err := testClient()
	if err != nil {
		t.Error(err.Error())
	}

	clearKVPath(t, client, "consul-alerts/config/notif-selection/")
	dataMap := make(map[string]string)
	dataMap["^_nomad-.*$"] = "client-profile"
	data, err := json.Marshal(dataMap)
	if err != nil {
		t.Error(err.Error())
	}
	client.api.KV().Put(&consulapi.KVPair{
		Key:   "consul-alerts/config/notif-selection/services",
		Value: data}, nil)
	profile := client.getProfileForEntity("service", "_nomad-client")
	if profile != "client-profile" {
		t.Error("getProfileForEntity must have matched client-profile")
	}

	client.api.KV().Put(&consulapi.KVPair{
		Key:   "consul-alerts/config/notif-selection/services/_nomad-server",
		Value: []byte("server-profile")}, nil)
	profile = client.getProfileForEntity("service", "_nomad-server")
	if profile != "server-profile" {
		t.Error("getProfileForEntity must have matched server-profile")
	}
}

func TestGetProfileInfo(t *testing.T) {
	client, err := testClient()
	if err != nil {
		t.Error(err.Error())
	}
	clearKVPath(t, client, "consul-alerts/config/notif-selection/")

	// test the default profile
	notifiersList := map[string]bool{"log": true}
	interval := 10
	defaultProfileInfo := ProfileInfo{Interval: interval, NotifList: notifiersList}
	data, err := json.Marshal(defaultProfileInfo)
	if err != nil {
		t.Error(err.Error())
	}
	client.api.KV().Put(&consulapi.KVPair{
		Key:   "consul-alerts/config/notif-profiles/default",
		Value: data}, nil)
	profileNotifiersList, profileInterval := client.GetProfileInfo("node", "serviceID", "checkID")
	if !reflect.DeepEqual(notifiersList, profileNotifiersList) {
		t.Error("Default notifiers list is loaded incorrectly")
	}
	if interval != profileInterval {
		t.Error("Default interval is loaded incorrectly")
	}

	// test notifier-selection based on nodes
	notifiersList = map[string]bool{"email": true}
	interval = 2
	nodeProfileInfo := ProfileInfo{Interval: interval, NotifList: notifiersList}
	data, err = json.Marshal(nodeProfileInfo)
	if err != nil {
		t.Error(err.Error())
	}
	client.api.KV().Put(&consulapi.KVPair{
		Key:   "consul-alerts/config/notif-profiles/node-profile",
		Value: data}, nil)

	client.api.KV().Put(&consulapi.KVPair{
		Key:   "consul-alerts/config/notif-selection/hosts/node",
		Value: []byte("node-profile")}, nil)
	profileNotifiersList, profileInterval = client.GetProfileInfo("node", "serviceID", "checkID")
	if !reflect.DeepEqual(notifiersList, profileNotifiersList) || interval != profileInterval {
		t.Error("notif-selection based on nodes loaded an incorrect profile")
	}

	// test notifier-selection based on checks
	notifiersList = map[string]bool{"influxdb": true}
	interval = 99
	checkProfileInfo := ProfileInfo{Interval: interval, NotifList: notifiersList}
	data, err = json.Marshal(checkProfileInfo)
	if err != nil {
		t.Error(err.Error())
	}
	client.api.KV().Put(&consulapi.KVPair{
		Key:   "consul-alerts/config/notif-profiles/check-profile",
		Value: data}, nil)

	client.api.KV().Put(&consulapi.KVPair{
		Key:   "consul-alerts/config/notif-selection/checks/checkID",
		Value: []byte("check-profile")}, nil)
	profileNotifiersList, profileInterval = client.GetProfileInfo("node", "serviceID", "checkID")
	if !reflect.DeepEqual(notifiersList, profileNotifiersList) || interval != profileInterval {
		t.Error("notif-selection based on checks loaded an incorrect profile")
	}

	// test notifier-selection based on services
	notifiersList = map[string]bool{"slack": true}
	interval = 5
	serviceProfileInfo := ProfileInfo{Interval: interval, NotifList: notifiersList}
	data, err = json.Marshal(serviceProfileInfo)
	if err != nil {
		t.Error(err.Error())
	}
	client.api.KV().Put(&consulapi.KVPair{
		Key:   "consul-alerts/config/notif-profiles/service-profile",
		Value: data}, nil)

	client.api.KV().Put(&consulapi.KVPair{
		Key:   "consul-alerts/config/notif-selection/services/serviceID",
		Value: []byte("service-profile")}, nil)
	profileNotifiersList, profileInterval = client.GetProfileInfo("node", "serviceID", "checkID")
	if !reflect.DeepEqual(notifiersList, profileNotifiersList) || interval != profileInterval {
		t.Error("notif-selection based on services loaded an incorrect profile")
	}
}
