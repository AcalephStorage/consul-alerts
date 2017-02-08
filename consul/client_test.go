package consul

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"testing"

	consulapi "github.com/AcalephStorage/consul-alerts/Godeps/_workspace/src/github.com/hashicorp/consul/api"
)

func testClient() (*ConsulAlertClient, error) {
	return NewClient("192.168.10.10:8500", "dc1", "")
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
	checkProfileInfo := client.GetProfileInfo("node", "serviceID", "checkID")
	if !reflect.DeepEqual(checkProfileInfo, defaultProfileInfo) {
		t.Error("Default profile info is loaded incorrectly")
	}

	// test that profiles based on checks, nodes and services are loaded correctly
	var testCombinations = []struct {
		Interval       int
		NotifiersList  map[string]bool
		NotifProfile   string
		NotifSelection string
	}{
		{
			Interval:       2,
			NotifiersList:  map[string]bool{"email": true},
			NotifProfile:   "nodes",
			NotifSelection: "hosts/node",
		},
		{
			Interval:       5,
			NotifiersList:  map[string]bool{"slack": true},
			NotifProfile:   "checks",
			NotifSelection: "checks/checkID",
		},
		{
			Interval:       99,
			NotifiersList:  map[string]bool{"influxdb": true},
			NotifProfile:   "services",
			NotifSelection: "services/serviceID",
		},
	}
	for _, s := range testCombinations {
		profileInfo := ProfileInfo{Interval: s.Interval, NotifList: s.NotifiersList}
		data, err = json.Marshal(profileInfo)
		if err != nil {
			t.Error(err.Error())
		}
		client.api.KV().Put(&consulapi.KVPair{
			Key:   fmt.Sprintf("consul-alerts/config/notif-profiles/%s", s.NotifProfile),
			Value: data}, nil)
		client.api.KV().Put(&consulapi.KVPair{
			Key:   fmt.Sprintf("consul-alerts/config/notif-selection/%s", s.NotifSelection),
			Value: []byte(s.NotifProfile)}, nil)
		checkProfileInfo := client.GetProfileInfo("node", "serviceID", "checkID")
		if !reflect.DeepEqual(checkProfileInfo, profileInfo) {
			t.Error("Profile info is loaded incorrectly")
		}
	}
}

func TestIsBlacklisted(t *testing.T) {
	client, err := testClient()
	if err != nil {
		t.Error(err.Error())
	}
	clearKVPath(t, client, "consul-alerts/config/checks/blacklist/")
	node := "test-node"
	checkID := "test-check"
	serviceID := "test-service"
	check := Check{Node: node, CheckID: checkID, ServiceID: serviceID}
	isBlackListed := client.IsBlacklisted(&check)
	if isBlackListed {
		t.Error("isBlackListed should be false if there is no corresponding entry in the blacklist")
	}

	testCombinations := []map[string]string{
		{"type": "node",
			"key": fmt.Sprintf("consul-alerts/config/checks/blacklist/nodes/%s", node)},
		{"type": "service",
			"key": fmt.Sprintf("consul-alerts/config/checks/blacklist/services/%s", serviceID)},
		{"type": "check",
			"key": fmt.Sprintf("consul-alerts/config/checks/blacklist/checks/%s", checkID)},
		{"type": "node-service-check combination",
			"key": fmt.Sprintf("consul-alerts/config/checks/blacklist/single/%s/%s/%s",
				node, serviceID, checkID)},
	}

	// test that blacklisting the exact key works
	for _, m := range testCombinations {
		clearKVPath(t, client, "consul-alerts/config/checks/blacklist/")
		client.api.KV().Put(&consulapi.KVPair{
			Key:   m["key"],
			Value: []byte{}}, nil)
		isBlackListed = client.IsBlacklisted(&check)
		if !isBlackListed {
			t.Errorf("isBlackListed should be true if the %s is blacklisted", m["type"])
		}
	}

	// test that blacklisting by regexp works
	testCombinations = []map[string]string{
		{"type": "node",
			"regexp": `["test-.*", "111"]`},
		{"type": "service",
			"regexp": `["test-.*"]`},
		{"type": "check",
			"regexp": `["test-.*", ""]`},
	}
	for _, m := range testCombinations {
		clearKVPath(t, client, "consul-alerts/config/checks/blacklist/")
		client.api.KV().Put(&consulapi.KVPair{
			Key:   fmt.Sprintf("consul-alerts/config/checks/blacklist/%ss", m["type"]),
			Value: []byte(m["regexp"])}, nil)
		isBlackListed = client.IsBlacklisted(&check)
		if !isBlackListed {
			t.Errorf("isBlackListed should be true if there is a regexp for %s matching the key",
				m["type"])
		}
	}
}
