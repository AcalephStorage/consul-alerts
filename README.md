consul-alerts
=============

[![Join the chat at https://gitter.im/AcalephStorage/consul-alerts](https://badges.gitter.im/Join%20Chat.svg)](https://gitter.im/AcalephStorage/consul-alerts?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge&utm_content=badge)

A highly available daemon to send notifications and reminders based on Consul health checks.  Including profile selection based on service, check, or host that enables specific handlers and reminder intervals.

## Requirement

1. Consul 0.4+. Get it [here](http://consul.io).
2. Configured `GOPATH`.

Releases
--------

Stable release are [here](https://github.com/AcalephStorage/consul-alerts/releases).

Latest release are found here:
 - [darwin-amd64](https://bintray.com/artifact/download/darkcrux/generic/consul-alerts-latest-darwin-amd64.tar)
 - [FreeBSD-amd64](https://bintray.com/artifact/download/darkcrux/generic/consul-alerts-latest-FreeBSD-amd64.tar)
 - [linux-386](https://bintray.com/artifact/download/darkcrux/generic/consul-alerts-latest-linux-386.tar)
 - [linux-amd64](https://bintray.com/artifact/download/darkcrux/generic/consul-alerts-latest-linux-amd64.tar)

Installation
------------

```
$ go get github.com/AcalephStorage/consul-alerts
$ go install
```

This should install consul-alerts to `$GOPATH/bin`

or pull the image from `docker`:

```
$ docker pull acaleph/consul-alerts
```

Usage
-----

```
$ consul-alerts start
```

or using docker:

```
$ docker run acaleph/consul-alerts start
```

By default, this runs the daemon and API at localhost:9000 and connects to the local consul agent (localhost:8500) and default datacenter (dc1). These can be overriden by the following flags:

```
$ consul-alerts start --alert-addr=localhost:9000 --consul-addr=localhost:8500 --consul-dc=dc1 --consul-acl-token=""
```

Note: Don't change --alert-addr when using the docker container.

Once the daemon is running, it can act as a handler for consul watches. At the moment only checks and events are supported.

```
$ consul watch -type checks consul-alerts watch checks --alert-addr=localhost:9000
$ consul watch -type event consul-alerts watch event --alert-addr=localhost:9000
```

or run the watchers on the agent the daemon connects by adding the following flags during consul-alerts run:

```
$ consul-alerts start --watch-events --watch-checks
```

Configuration
-------------

All configurations are stored in consul's KV with the prefix: `consul-alerts/config/`. The daemon is using default values and the KV entries will only override the defaults.

### Health Checks

Health checking is enabled by default. This also triggers the notification when a check has changed status for a configured duration. Health checks can be disabled by setting the kv`consul-alerts/config/checks/enabled` to `false`.

To prevent flapping, notifications are only sent when a check status has been stable for a specific time in seconds (60 by default). this value can be changed by adding/changing the kv `consul-alerts/config/checks/change-threshold` to an integer greater than and divisible by 10.

eg. `consul-alerts/config/checks/change-threshold` = `30`

#### Enable Profiles Selection
Profiles may be configured as keys in consul-alerts/config/notif-profiles/.  The key name is the name of the profile and the value should be a JSON object with an "Interval" key set to an int in minutes and a key "NotifList" that should be an object of profile names as keys and true for the value. 

Profile selection is done by setting keys in consul-alerts/config/notif-selection/services/, consul-alerts/config/notif-selection/checks/, or consul-alerts/config/notif-selection/hosts/ with the appropriate service, check, or host name as the key and the selected profile name as the value.

Reminders are processed every five minutes.  Interval values should be a multiple of five.  If the Interval value is 0 or not set then reminders will not be set.

The default profile may be set as the fallback to any checks that do not match a selection.  If there is no default profile set then the full list of enabled notifiers will be used and no reminders.

#### Enable/Disable Specific Health Checks

There are four ways to enable/disable health check notifications: mark them by node, serviceID, checkID, or mark individually by node/serviceID/checkID. This is done by adding a KV entry in `consul-alerts/config/checks/blacklist/...`. Removing the entry will re-enable the check notifications.

##### Disable all notification by node

Add a KV entry with the key `consul-alerts/config/checks/blacklist/nodes/{{ nodeName }}`. This will disable notifications for the specified `nodeName`.

##### Disable all notification by service

Add a KV entry with the key `consul-alerts/config/checks/blacklist/services/{{ serviceId }}`. This will disable notifications for the specified `serviceId`.

##### Disable all notification by healthCheck

Add a KV entry with the key `consul-alerts/config/checks/blacklist/checks/{{ checkId }}`. This will disable notifications for the specified `checkId`.

##### Disable a single health check

Add a KV entry with the key `consul-alerts/config/checks/blacklist/single/{{ node }}/{{ serviceId }}/{{ checkId }}`. This will disable the specific health check. If the health check is not associated with a service, use the `_` as the serviceId.

### Events

Event handling is enabled by default. This delegates any consul event received by the agent to the list of handlers configured. To disable event handling, set `consul-alerts/config/events/enabled` to `false`.

Handlers can be configured by adding them to `consul-alerts/config/events/handlers`. This should be a JSON array of string. Each string should point to any executable. The event data should be read from `stdin`.

### Notifiers

There are four builtin notifiers. Only the *Log* notifier is enabled by default. It is also possible to add custom notifiers similar to custom event handlers. Custom notifiers can be added as keys with command path string values in `consul-alerts/config/notifiers/custom/`. The keys will be used as notifier names in the profiles.

#### Logger

This logs any health check notification to a file. To disable this notifier, set `consul-alerts/config/notifiers/log/enabled` to `false`.

The log file is set to `/tmp/consul-notifications.log` by default. This can be changed by changing `consul-alerts/config/notifiers/log/path`.

#### Email

This emails the health notifications. To enable this, set `consul-alerts/config/notifiers/email/enabled` to `true`.

The email and smtp details needs to be configured:

prefix: `consul-alerts/config/notifiers/email/`

| key          | description                                                 |
|--------------|-------------------------------------------------------------|
| enabled      | Enable the email notifier. [Default: false]                 |
| cluster-name | The name of the cluster. [Default: "Consul Alerts"]         |
| url          | The SMTP server url                                         |
| port         | The SMTP server port                                        |
| username     | The SMTP username                                           |
| password     | The SMTP password                                           |
| sender-alias | The sender alias. [Default: "Consul Alerts"]                |
| sender-email | The sender email                                            |
| receivers    | The emails of the receivers. JSON array of string           |
| template     | Path to custom email template. [Default: internal template] |

The template can be any go html template. An `EmailData` instance will be passed to the template.

#### InfluxDB

This sends the notifications as series points in influxdb. Set `consul-alerts/config/notifiers/influxdb/enabled` to `true` to enabled. InfluxDB details need to be set too.

prefix: `consul-alerts/config/notifiers/influxdb/`

| key         | description                                    |
|-------------|------------------------------------------------|
| enabled     | Enable the influxdb notifier. [Default: false] |
| host        | The influxdb host. (eg. http://localhost:8086) |
| username    | The influxdb username                          |
| password    | The influxdb password                          |
| database    | The influxdb database name                     |
| series-name | The series name for the points                 |

#### Slack

Slack integration is also supported. To enable, set
`consul-alerts/config/notifiers/slack/enabled` to `true`. Slack details needs to
be configured.

prefix: `consul-alerts/config/notifiers/slack/`

| key          | description                                         |
|--------------|-----------------------------------------------------|
| enabled      | Enable the Slack notifier. [Default: false]         |
| cluster-name | The name of the cluster. [Default: "Consul Alerts"] |
| url          | The incoming-webhook url (mandatory)                |
| channel      | The channel to post the notification (mandatory)    |
| username     | The username to appear on the post                  |
| icon-url     | URL of a custom image for the notification          |
| icon-emoji   | Emoji (if not using icon-url) for the notification  |

In order to enable slack integration, you have to create a new
[_Incoming WebHooks_](https://my.slack.com/services/new/incoming-webhook). Then use the
token created by the previous action.

#### PagerDuty

To enable PagerDuty built-in notifier, set `consul-alerts/config/notifiers/pagerduty/enabled` to `true`. This is disabled by default. Service key and client details also needs to be configured.

prefix: `consul-alerts/config/notifiers/pagerduty/`

| key         | description                                     |
|-------------|-------------------------------------------------|
| enabled     | Enable the PagerDuty notifier. [Default: false] |
| service-key | Service key to access PagerDuty                 |
| client-name | The monitoring client name                      |
| client-url  | The monitoring client url                       |

#### Hipchat

To enable HipChat builtin notifier, set
`consul-alerts/config/notifiers/hipchat/enabled` to `true`. Hipchat details
needs to be configured.

prefix: `consul-alerts/config/notifiers/hipchat/`

| key          | description                                         |
|--------------|-----------------------------------------------------|
| enabled      | Enable the Hipchat notifier. [Default: false]       |
| from         | The name to send notifications as                   |
| cluster-name | The name of the cluster. [Default: "Consul Alerts"] |
| base-url     | HipChat base url                                    |
| room-id      | The room to post to                  (mandatory)    |
| auth-token   | Authentication token                 (mandatory)    |

#### OpsGenie

To enable OpsGenie builtin notifier, set
`consul-alerts/config/notifiers/opsgenie/enabled` to `true`. OpsGenie details
needs to be configured.

prefix: `consul-alerts/config/notifiers/opsgenie/`

| key          | description                                         |
|--------------|-----------------------------------------------------|
| enabled      | Enable the OpsGenie notifier. [Default: false]      |
| cluster-name | The name of the cluster. [Default: "Consul Alerts"] |
| api-key      | API Key                              (mandatory)    |

#### Amazon Web Services Simple Notification Service ("SNS")

To enable AWS SNS builtin notifier, set
`consul-alerts/config/notifiers/awssns/enabled` to `true`. AWS SNS details
needs to be configured.

prefix: `consul-alerts/config/notifiers/awssns/`

| key          | description                                         |
|--------------|-----------------------------------------------------|
| enabled      | Enable the AWS SNS notifier. [Default: false]       |
| region       | AWS Region                           (mandatory)    |
| topic-arn    | Topic ARN to publish to.             (mandatory)    |

Health Check via API
--------------------

Health status can also be queried via the API. This can be used for compatibility with nagios, sensu, or other monitoring tools. To get the status of a specific check, use the following entrypoint.

`http://consul-alerts:9000/v1/health?node=<node>&service=<serviceId>&check=<checkId>`

This will return the output of the check and the following HTTP codes:

| Status   | Code |
|----------|------|
| passing  | 200  |
| warning  | 503  |
| critical | 503  |
| unknown  | 404  |

Contribution
------------

PRs are more than welcome. Just fork, create a feature branch, and open a PR. We love PRs. :)

TODO
----

Needs better doc and some cleanup too. :)
