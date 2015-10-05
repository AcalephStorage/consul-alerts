package client

import (
	"errors"
	"fmt"
	alerts "github.com/AcalephStorage/consul-alerts/Godeps/_workspace/src/github.com/opsgenie/opsgenie-go-sdk/alerts"
	hb "github.com/AcalephStorage/consul-alerts/Godeps/_workspace/src/github.com/opsgenie/opsgenie-go-sdk/heartbeat"
	yaml "gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"testing"
)

type ClientTestConfig struct {
	Alert struct {
		ApiKey  string   `yaml:"apiKey"`
		User    string   `yaml:"user"`
		Team    string   `yaml:"team"`
		Actions []string `yaml:"actions"`
	} `yaml:"alert"`

	Heartbeat struct {
		ApiKey string `yaml:"apiKey"`
	} `yaml:"heartbeat"`
}

// common globals
var cli *OpsGenieAlertClient
var CONFIG_FILE_NAME string = "client_at_test_cfg.yaml"
var testCfg ClientTestConfig

// global for alert id
var alertId string

// global for heartbeat id
var hbId string
var hbName string
var hbCli *OpsGenieHeartbeatClient

func TestCreateAlert(t *testing.T) {
	req := alerts.CreateAlertRequest{
		Message: "[AT] Test Alert",
		Note:    "Created for testing purposes",
		User:    testCfg.Alert.User,
		Actions: testCfg.Alert.Actions,
	}
	response, alertErr := cli.Create(req)

	if alertErr != nil {
		t.Errorf(alertErr.Error())
	}

	if response.Code >= 400 {
		t.Errorf(fmt.Sprintf("Creating alert failed with response code: %d", response.Code))
	}

	alertId = response.AlertId
	t.Log("[OK] alert created")
}

func TestAcknowledgeAlert(t *testing.T) {
	ackReq := alerts.AcknowledgeAlertRequest{AlertId: alertId}
	ackResponse, alertErr := cli.Acknowledge(ackReq)

	if alertErr != nil {
		t.Errorf(alertErr.Error())
	}

	if ackResponse.Code >= 400 {
		t.Errorf(fmt.Sprintf("Acknowledge alert failed with response code: %d", ackResponse.Code))
	}
	t.Log("[OK] alert acked")
}

func TestAddNoteAlert(t *testing.T) {
	addnotereq := alerts.AddNoteAlertRequest{}
	// add alert ten notes
	for i := 0; i < 10; i++ {
		addnotereq.AlertId = alertId
		addnotereq.Note = fmt.Sprintf("Alert note # %d", i)
		addnoteresp, alertErr := cli.AddNote(addnotereq)
		if alertErr != nil {
			t.Errorf(alertErr.Error())
		}
		if addnoteresp.Code >= 400 {
			t.Errorf(fmt.Sprintf("Add alert note failed with response code: %d", addnoteresp.Code))
		}
	}
	t.Log("[OK] notes added to alert")
}

func TestListNotes(t *testing.T) {
	listNotesReq := alerts.ListAlertNotesRequest{Id: alertId}
	listNotesResponse, alertErr := cli.ListNotes(listNotesReq)
	if alertErr != nil {
		t.Errorf(alertErr.Error())
	}

	alertNotes := listNotesResponse.Notes
	if len(alertNotes) != 11 {
		t.Errorf("Retrieving all alert notes failed")
	}
	t.Log("[OK] alert notes listed")
}

func TestAddTeam(t *testing.T) {
	addTeamReq := alerts.AddTeamAlertRequest{AlertId: alertId, Team: testCfg.Alert.Team}
	addTeamResponse, alertErr := cli.AddTeam(addTeamReq)

	if alertErr != nil {
		t.Errorf(alertErr.Error())
	}

	if addTeamResponse.Code >= 400 {
		t.Errorf(fmt.Sprintf("Add team request failed with response code: %d", addTeamResponse.Code))
	}
	t.Log("[OK] team added to alert")
}

func TestAssignOwner(t *testing.T) {
	assignOwnerReq := alerts.AssignOwnerAlertRequest{AlertId: alertId, Owner: testCfg.Alert.User}
	assignOwnerResponse, alertErr := cli.AssignOwner(assignOwnerReq)

	if alertErr != nil {
		t.Errorf(alertErr.Error())
	}

	if assignOwnerResponse.Code >= 400 {
		t.Errorf(fmt.Sprintf("Assign owner request failed with response code: %d", assignOwnerResponse.Code))
	}
	t.Log("[OK] owner assigned to alert")
}

func TestExecuteAction(t *testing.T) {
	execActionReq := alerts.ExecuteActionAlertRequest{AlertId: alertId,
		Action: testCfg.Alert.Actions[0],
		Note:   fmt.Sprintf("Action <b>%s</b> executed by the Go API", testCfg.Alert.Actions[0]),
	}
	execActionResponse, alertErr := cli.ExecuteAction(execActionReq)

	if alertErr != nil {
		t.Errorf(alertErr.Error())
	}

	if execActionResponse.Code >= 400 {
		t.Errorf(fmt.Sprintf("Executing action failed with response code: %d", execActionResponse.Code))
	}
	t.Log("[OK] action %s executed on alert", testCfg.Alert.Actions[0])
}

func TestDeleteAlert(t *testing.T) {
	delreq := alerts.DeleteAlertRequest{AlertId: alertId, Source: "Go API Test"}
	delResp, alertErr := cli.Delete(delreq)
	if alertErr != nil {
		t.Errorf(alertErr.Error())
	}
	if delResp.Code >= 400 {
		t.Errorf(fmt.Sprintf("Delete alert request failed with response code: %d", delResp.Code))
	}
	t.Log("[OK] alert deleted")
}

//
// Heartbeat tests
//
func TestAddHeartbeat(t *testing.T) {
	req := hb.AddHeartbeatRequest{Name: "[AT] Test Heartbeat"}
	response, hbErr := hbCli.Add(req)

	if hbErr != nil {
		t.Errorf("Add heartbeat request failed: " + hbErr.Error())
	}

	if response == nil {
		t.FailNow()
	}
	if response.Code >= 400 {
		t.Errorf(fmt.Sprintf("Add heartbeat failed with error code %d", response.Code))
	}
	// set the heartbeat id to be used for the following testing functions
	hbId = response.Id
}

func TestUpdateHeartbeat(t *testing.T) {
	updateReq := hb.UpdateHeartbeatRequest{Id: hbId, Name: "[AT] Test Heartbeat Updated", Description: "Some description"}
	updateResp, updateErr := hbCli.Update(updateReq)

	if updateErr != nil {
		t.Errorf("Update heartbeat request failed")
	}
	if updateResp == nil {
		t.FailNow()
	}
	if updateResp.Code >= 400 {
		t.Errorf("Update heartbeat request failed with error code %d", updateResp.Code)
	}

	hbName = "[AT] Test Heartbeat Updated"
}

func TestEnableHeartbeat(t *testing.T) {
	enableReq := hb.EnableHeartbeatRequest{Id: hbId}
	enableResp, enableErr := hbCli.Enable(enableReq)

	if enableErr != nil {
		t.Errorf("Enable heartbeat request failed")
	}
	if enableResp == nil {
		t.FailNow()
	}
	if enableResp.Code >= 400 {
		t.Errorf("Enable heartbeat request failed with error code %d", enableResp.Code)
	}
}

func TestSendHeartbeat(t *testing.T) {
	sendReq := hb.SendHeartbeatRequest{Name: hbName}
	sendResp, sendErr := hbCli.Send(sendReq)

	if sendErr != nil {
		t.Errorf("Send heartbeat request failed")
	}
	if sendResp == nil {
		t.FailNow()
	}
	if sendResp.Code >= 400 {
		t.Errorf("Send heartbeat request failed with error code %d", sendResp.Code)
	}
}

func TestDisableHeartbeat(t *testing.T) {
	disableReq := hb.DisableHeartbeatRequest{Id: hbId}
	disableResp, disableErr := hbCli.Disable(disableReq)

	if disableErr != nil {
		t.Errorf("Disable heartbeat request failed")
	}
	if disableResp == nil {
		t.FailNow()
	}
	if disableResp.Code >= 400 {
		t.Errorf("Disable heartbeat request failed with error code %d", disableResp.Code)
	}
}

func TestListHeartbeats(t *testing.T) {
	listReq := hb.ListHeartbeatsRequest{}
	listResp, listErr := hbCli.List(listReq)
	if listErr != nil {
		t.Errorf("List heartbeat request failed")
	}
	if listResp == nil {
		t.FailNow()
	}
	beats := listResp.Heartbeats

	if len(beats) == 0 {
		t.Errorf("No heartbeat found, expecting at least 1 heartbeat")
	}

	found := false

	for _, beat := range beats {
		if beat.Id == hbId {
			found = true
		}
	}

	if found == false {
		t.Errorf(fmt.Sprintf("Newly created heartbeat with id %s not found in the heartbeats list", hbId))
	}
}

func TestDeleteHeartbeat(t *testing.T) {
	deleteReq := hb.DeleteHeartbeatRequest{Id: hbId}
	deleteResp, deleteErr := hbCli.Delete(deleteReq)
	if deleteErr != nil {
		t.Errorf("Delete heartbeat request failed")
	}
	if deleteResp == nil {
		t.FailNow()
	}
	if deleteResp.Code >= 400 {
		t.Errorf("Delete heartbeat request failed with error code %d", deleteResp.Code)
	}
}

// utility function
func readSettingsFromConfigFile() error {
	cfgData, err := ioutil.ReadFile(CONFIG_FILE_NAME)
	if err != nil {
		return errors.New("Can not read from the configuration file: " + err.Error())
	}
	err = yaml.Unmarshal(cfgData, &testCfg)
	if err != nil {
		return errors.New("Can not parse the configuration file: " + err.Error())
	}
	return nil
}

// setup the test suite
func TestMain(m *testing.M) {
	// read the settings
	err := readSettingsFromConfigFile()
	if err != nil {
		panic(err)
	}
	// create an opsgenie client
	opsGenieClient := new(OpsGenieClient)
	opsGenieClient.SetApiKey(testCfg.Alert.ApiKey)
	// create the alerting client
	var cliErr error
	cli, cliErr = opsGenieClient.Alert()

	if cliErr != nil {
		panic(cliErr)
	}

	// create the heartbeat client
	// Api Key should be switched in order to send heartbeat requests
	opsGenieClient.SetApiKey(testCfg.Heartbeat.ApiKey)
	hbCli, cliErr = opsGenieClient.Heartbeat()

	if cliErr != nil {
		panic(cliErr)
	}
	os.Exit(m.Run())
}
