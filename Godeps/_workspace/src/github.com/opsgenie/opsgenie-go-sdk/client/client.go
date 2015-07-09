// Copyright 2015 OpsGenie. All rights reserved.
// Use of this source code is governed by a Apache Software 
// license that can be found in the LICENSE file.

/*
	Package client manages the creation of API clients. 
	API user first creates a pointer of type OpsGenieClient. Following that
	he/she can set some configurations for HTTP communication layer by setting 
	a proxy definition and/or transport layer options. 

	Introduction

	The most fundamental and general use case is being able to access the 
	OpsGenie Web API by coding a Go program.
	The program -by mean of a client application- can send OpsGenie Web API 
	the requests using the 'client' package in a higher level. For the programmer 
	of the client application, that reduces the number of LoCs.
	Besides it will result a less error-prone application and reduce 
	the complexity by hiding the low-level networking, error-handling and 
	byte-processing calls.

	Package client has ports for all entry points to the Web API. 
	The OpsGenie Web API is structured in JSON-bodied 
	calls (except the file attachment).
*/
package client

import (
	"errors"
	"time"
	"runtime"
	"fmt"
)

// OpsGenie Go SDK performs HTTP calls to the Web API.
// The Web API is designated by a URL so called an endpoint
const ENDPOINT_URL string = "https://api.opsgenie.com" 

const DEFAULT_CONNECTION_TIMEOUT_IN_SECONDS time.Duration = 1
const DEFAULT_MAX_RETRY_ATTEMPTS int = 1
const TIME_SLEEP_BETWEEN_REQUESTS time.Duration = 500 * time.Millisecond

// User-Agent values tool/version (OS;GO_Version;language)
type RequestHeaderUserAgent struct {
	sdkName 	string
	version 	string
	os 			string
	goVersion 	string
	timezone	string
}

func (p RequestHeaderUserAgent) ToString() string {
	return fmt.Sprintf("%s/%s (%s;%s;%s)", p.sdkName, p.version, p.os, p.goVersion, p.timezone)	
}

var userAgentParam RequestHeaderUserAgent

// OpsGenieClient is a general data type used for:
// 	- authenticating callers through their api keys and 
// 	- instanciating "alert" and "heartbeat" clients
//	- setting HTTP transport layer configurations
type OpsGenieClient struct {
	proxy *ClientProxyConfiguration
	httpTransportSettings *HttpTransportSettings
	apiKey string
}
// Setters:
//	- proxy
//	- http transport layer conf
//	- api key
func (cli *OpsGenieClient) SetClientProxyConfiguration(conf *ClientProxyConfiguration) {
	cli.proxy = conf
}

func (cli *OpsGenieClient) SetHttpTransportSettings(settings *HttpTransportSettings) {
	cli.httpTransportSettings = settings
}

func (cli *OpsGenieClient) SetApiKey(key string) error {
	if key == "" {
		return errors.New("API Key can not be empty")
	}
	cli.apiKey = key
	return nil
}

// Instanciates a new OpsGenieAlertClient
// and sets the api key to be used alongside the execution.
func (cli *OpsGenieClient) Alert() (*OpsGenieAlertClient, error) {
	if cli.apiKey == "" {
		return nil, errors.New("API Key should be set first")
	}
	alertClient := new (OpsGenieAlertClient)
	alertClient.apiKey = cli.apiKey
	if cli.proxy != nil {
		alertClient.proxy = cli.proxy.ToString()	
	}
	alertClient.SetConnectionTimeout( DEFAULT_CONNECTION_TIMEOUT_IN_SECONDS * time.Second )
	alertClient.SetMaxRetryAttempts( DEFAULT_MAX_RETRY_ATTEMPTS )

	if cli.httpTransportSettings != nil {
		if cli.httpTransportSettings.ConnectionTimeout > 0 {
			alertClient.SetConnectionTimeout( cli.httpTransportSettings.ConnectionTimeout )			
		}
		if cli.httpTransportSettings.MaxRetryAttempts > 0 {
			alertClient.SetMaxRetryAttempts(cli.httpTransportSettings.MaxRetryAttempts)
		}
	}
	return alertClient, nil
}
// Instanciates a new OpsGenieHeartbeatClient
// and sets the api key to be used alongside the execution.
func (cli *OpsGenieClient) Heartbeat() (*OpsGenieHeartbeatClient, error) {
	if cli.apiKey == "" {
		return nil, errors.New("API Key should be set first")
	}
	heartbeatClient := new (OpsGenieHeartbeatClient)
	heartbeatClient.apiKey = cli.apiKey
	if cli.proxy != nil {
		heartbeatClient.proxy = cli.proxy.ToString()	
	}
	heartbeatClient.SetConnectionTimeout(DEFAULT_CONNECTION_TIMEOUT_IN_SECONDS * time.Second)
	heartbeatClient.SetMaxRetryAttempts(DEFAULT_MAX_RETRY_ATTEMPTS)
	if cli.httpTransportSettings != nil {
		if cli.httpTransportSettings.ConnectionTimeout > 0 {
			heartbeatClient.SetConnectionTimeout(cli.httpTransportSettings.ConnectionTimeout)			
		}
	}

	return heartbeatClient, nil
}
// Instanciates a new OpsGenieIntegrationClient
// and sets the api key to be used alongside the execution.
func (cli *OpsGenieClient) Integration() (*OpsGenieIntegrationClient, error) {
	if cli.apiKey == "" {
		return nil, errors.New("API Key should be set first")
	}
	integrationClient := new (OpsGenieIntegrationClient)
	integrationClient.apiKey = cli.apiKey
	if cli.proxy != nil {
		integrationClient.proxy = cli.proxy.ToString()	
	}
	integrationClient.SetConnectionTimeout(DEFAULT_CONNECTION_TIMEOUT_IN_SECONDS * time.Second)
	integrationClient.SetMaxRetryAttempts(DEFAULT_MAX_RETRY_ATTEMPTS)	
	if cli.httpTransportSettings != nil {
		if cli.httpTransportSettings.ConnectionTimeout > 0 {
			integrationClient.SetConnectionTimeout(cli.httpTransportSettings.ConnectionTimeout)			
		}
	}	
	return integrationClient, nil
}

// Instanciates a new OpsGeniePolicyClient
// and sets the api key to be used alongside the execution.
func (cli *OpsGenieClient) Policy() (*OpsGeniePolicyClient, error) {
	if cli.apiKey == "" {
		return nil, errors.New("API Key should be set first")
	}
	policyClient := new (OpsGeniePolicyClient)
	policyClient.apiKey = cli.apiKey
	if cli.proxy != nil {
		policyClient.proxy = cli.proxy.ToString()	
	}
	policyClient.SetConnectionTimeout(DEFAULT_CONNECTION_TIMEOUT_IN_SECONDS * time.Second)
	policyClient.SetMaxRetryAttempts(DEFAULT_MAX_RETRY_ATTEMPTS)	
	if cli.httpTransportSettings != nil {
		if cli.httpTransportSettings.ConnectionTimeout > 0 {
			policyClient.SetConnectionTimeout(cli.httpTransportSettings.ConnectionTimeout)	
		}
	}	
	return policyClient, nil
}

// Initializer for the package client
// Initializes the User-Agent parameter of the requests.
// TODO version information must be read from a MANIFEST file
func init() {
	userAgentParam.sdkName = "opsgenie-go-sdk"
	userAgentParam.version = "1.0.0"
	userAgentParam.os = runtime.GOOS
	userAgentParam.goVersion = runtime.Version()
	userAgentParam.timezone = time.Local.String()
}
