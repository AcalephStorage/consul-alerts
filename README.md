consul-alerts
=============

[![Join the chat at https://gitter.im/AcalephStorage/consul-alerts](https://badges.gitter.im/Join%20Chat.svg)](https://gitter.im/AcalephStorage/consul-alerts?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge&utm_content=badge)

A highly available daemon for sending notifications and reminders based on Consul health checks.

Under the covers, consul-alerts leverages Consul's own leadership election and KV store to provide automatic failover and seamless operation in the case of a consul-alerts node failure and ensures that your notifications are still sent.

consul-alerts provides a high degree of configuration including:

-  Several built-in [Notifiers](#notifiers) for distribution of health check alerts (email, sns, pagerduty, etc.)
-  The ability to create Notification Profiles, sets of Notifiers which will respond to the given alert when a configurable threshold is exceeded
-  Multiple degrees of customization for Notifiers and Blacklisting of alerts (service, check id or host)


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

By default, this runs the daemon and API at localhost:9000 and connects to the local consul agent (localhost:8500) and default datacenter (dc1). These can be overriden by the following flags:

```
$ consul-alerts start --alert-addr=localhost:9000 --consul-addr=localhost:8500 --consul-dc=dc1 --consul-acl-token=""
```

Once the daemon is running, it can act as a handler for consul watches. At the moment only checks and events are supported.

```
$ consul watch -type checks consul-alerts watch checks --alert-addr=localhost:9000
$ consul watch -type event consul-alerts watch event --alert-addr=localhost:9000
```

or run the watchers on the agent the daemon connects by adding the following flags during consul-alerts run:

```
$ consul-alerts start --watch-events --watch-checks
```

Usage - Docker
--------------

There are a few options for running in docker.

First option is using the consul agent built into the container. This option requires overriding the default entry point and running an exec to launch consul alerts.

Start consul:

```
docker run -ti \
  --rm -p 9000:9000 \
  --hostname consul-alerts \
  --name consul-alerts \
  --entrypoint=/bin/consul \
  acaleph/consul-alerts \
  agent -data-dir /data -server -bootstrap -client=0.0.0.0
```
Then in a separate terminal start consul-alerts:

```
$ docker exec -ti consul-alerts /bin/consul-alerts start --alert-addr=0.0.0.0:9000 --log-level=info --watch-events --watch-checks
```

The second option is to link to an existing consul container through docker networking and --link option.  This method can more easily
share the consul instance with other containers such as vault.

First launch consul container:

```
$ docker run \
  -p 8400:8400 \
  -p 8500:8500 \
  -p 8600:53/udp \
  --hostname consul \
  --name consul \
  progrium/consul \
  -server -bootstrap -ui-dir /ui
```

Then run consul alerts container:

```
$ docker run -ti \
  -p 9000:9000 \
  --hostname consul-alerts \
  --name consul-alerts \
  --link consul:consul \
  acaleph/consul-alerts start \
  --consul-addr=consul:8500 \
  --log-level=info --watch-events --watch-checks
```

Last option is to launch the container and point at a remote consul instance:

```
$ docker run -ti \
  -p 9000:9000 \
  --hostname consul-alerts \
  --name consul-alerts \
  acaleph/consul-alerts start \
  --consul-addr=remote-consul-server.domain.tdl:8500 \
  --log-level=info --watch-events --watch-checks
```

**NOTE:** Don't change --alert-addr when using the docker container.

Configuration
-------------

To assure consistency between instances, configuration is stored in Consul's KV with the prefix: `consul-alerts/config/`. consul-alerts works out of the box without any customizations by using the defaults documented below and leverages the KV settings as overrides.

A few suggestions on operating and bootstrapping your consul-alerts configuration via the KV store are located in the [Operations](#operations) section below.

If **ACL**s are enabled the folowing policy should be configured for consul-alerts token:

```
key "consul-alerts" {
  policy = "write"
}

service "" {
  policy = "read"
}

event "" {
  policy = "read"
}

session "" {
  policy = "write"
}
```

### Health Checks

Health checking is enabled by default and is at the core what consul-alerts provides. The Health Check functionality is responsible for triggering a notification when the given consul check has changed status. To prevent flapping, notifications are only sent when a check status has been consistent for a specific time in seconds (60 by default). The threshold can be set globally or for a particular node, check, service and/or all of them.

**Configuration Options:**
The default Health Check configuration can be customized by setting kv with the prefix: `consul-alerts/config/checks/`

| key              | description                                                                                        |
|------------------|----------------------------------------------------------------------------------------------------|
| enabled          | Globally enable the Health Check functionality. [Default: true]                                    |
| change-threshold | The time, in seconds, that a check must be in a given status before an alert is sent [Default: 60] |
| single/{{ node }}/{{ serviceID }}/{{ checkID }}/change-threshold | Overrides `change-threshold` for a specific check associated with a particular service running on a particular node |
| check/{{ checkID }}/change-threshold | Overrides `change-threshold` for a specific check |
| service/{{ serviceID }}/change-threshold | Overrides `change-threshold` for a specific service |
| node/{{ node }}/change-threshold | Overrides `change-threshold` for a specific node |

When `change-threshold` is overridden multiple times, the most specific condition will be used based on the following order: (most specific) `single` > `check` > `service` > `node` > `global settings` > `default settings` (least specific).

### Notification Profiles

Notification Profiles allow the operator the ability to customize how often and to which Notifiers alerts will be sent via the Interval and NotifList attributes described below.

Profiles are configured as keys with the prefix: `consul-alerts/config/notif-profiles/`.

#### Notification Profile Specification

**Key:** The name of the Notification Profile

Ex. `emailer_only` would be located at `consul-alerts/config/notif-profiles/emailer_only`

**Value:** A JSON object adhering to the schema shown below.

```
{
  "$schema": "http://json-schema.org/draft-04/schema#",
  "type": "object",
  "title": "Notifier Profile Schema.",
  "description": "Defines a given Notifier Profile's configuration",
  "properties": {
    "Interval": {
      "type": "integer",
      "title": "Reminder Interval.",
      "description": "Defines the Interval (in minutes) which Reminders should be sent to the given Notifiers.  Should be a multiple of 5."
    },
    "NotifList": {
      "type": "object",
      "title": "Hash of Notifiers to configure.",
      "description": "A listing of Notifier names with a boolean value indicating if it should be enabled or not.",
      "patternProperties" : {
        ".{1,}" : { "type" : "string" }
      }
    },
    "VarOverrides": {
      "type": "object",
      "title": "Hash of Notifier variables to override.",
      "description": "A listing of Notifier names with hash values containing the parameters to be overridden",
      "patternProperties" : {
        ".{1,}" : { "type" : "object" }
      }
    }
  },
  "required": [
    "Interval",
    "NotifList"
  ]
}
```

#### Notification Profile Examples
**Notification Profile to only send Emails with reminders every 10 minutes:**

**Key:** `consul-alerts/config/notif-profiles/emailer_only`

**Value:**
```
{
  "Interval": 10,
  "NotifList": {
    "log":false,
    "email":true
  }
}
```

**NOTE:** While it is completely valid to explicitly disable a Notifier in a Notifier Profile, it is not necessary.  In the event that a Notification Profile is used, only Notifiers which are explicitly defined and enabled will be used.  In the example above then, we could have omitted the `"log": false` in the `NotifList` and achieved the same results.

**Example - Notification Profile to only send to PagerDuty but never send reminders:**

**Key:** `consul-alerts/config/notif-profiles/pagerduty_no_reminders`

**Value:**
```
{
  "Interval": 0,
  "NotifList": {
    "pagerduty":true
  }
}
```

**NOTE:** The Interval being set to 0 **disables** Reminders from being sent for a given alert.  If the service stays in a critical status for an extended period, only that first notification will be sent.

**Example - Notification Profile to only send Emails to the overridden receivers:**

**Key:** `consul-alerts/config/notif-profiles/emailer_overridden`

**Value:**
```
{
  "Interval": 10,
  "NotifList": {
    "email":true
  },
  "VarOverrides": {
    "email": {
      "receivers": ["my-team@company.com"]
    }
  }
}
```

#### Notification Profile Activation

It is possible to activate Notification Profiles in 2 ways - for a specific entity or for a set of entities matching a regular expression.
For a specific item the selection is done by setting keys in `consul-alerts/config/notif-selection/services/`, `consul-alerts/config/notif-selection/checks/`, or `consul-alerts/config/notif-selection/hosts/` with the appropriate service, check, or host name as the key and the desired Notification Profile name as the value.
To activate a Notification Profile for a set of entities matching a regular expression, create a json map of type `regexp->notification-profile` as a value for the keys `consul-alerts/config/notif-selection/services`, `consul-alerts/config/notif-selection/checks`, or `consul-alerts/config/notif-selection/hosts`.

**Example - Notification Profile activated for all the services which names start with infra-**

**Key:** `consul-alerts/config/notif-selection/services`

**Value:**
```
{
  "^infra-.*$": "infra-support-profile"
}
```

In addition to the service, check and host specific Notification Profiles, the operator can setup a default Notification Profile by creating a Notification Profile kv `consul-alerts/config/notif-profiles/default`, which acts as a fallback in the event a specific Notification Profile is not found.  If there are no Notification Profiles matching the criteria, consul-alerts will send the notification to the full list of enabled Notifiers and no reminders will be sent.

As consul-alerts attempts to process a given notification, it has a series of lookups it does to associate an event with a given Notification Profile by matching on:

- Service
- Check
- Host
- Default

**NOTE:** An event will only trigger notification for the FIRST Notification Profile that meets it's criteria.

Reminders resend the notifications at programmable intervals until they are resolved or added to the blacklist. Reminders are processed every five minutes therefore Interval values should be a multiple of five.  If the Interval value is 0 or not set then reminders will not be sent.

#### Enable/Disable Specific Health Checks

There are multiple ways to enable/disable health check notifications: mark them by node, serviceID, checkID, regular expression, or mark individually by node/serviceID/checkID. This is done by adding a KV entry in `consul-alerts/config/checks/blacklist/...`. Removing the entry will re-enable the check notifications.

##### Disable all notification by node

Add a KV entry with the key `consul-alerts/config/checks/blacklist/nodes/{{ nodeName }}`. This will disable notifications for the specified `nodeName`.

##### Disable all notifications for the nodes matching regular expressions

Add a KV entry with the key `consul-alerts/config/checks/blacklist/nodes` and the value containing a list of regular expressions. This will disable notifications for all the nodes, which names match at least one of the regular expressions.

##### Disable all notification by service

Add a KV entry with the key `consul-alerts/config/checks/blacklist/services/{{ serviceId }}`. This will disable notifications for the specified `serviceId`.

##### Disable all notifications for the services matching regular expressions

Add a KV entry with the key `consul-alerts/config/checks/blacklist/services` and the value containing a list of regular expressions. This will disable notifications for all the services, which names match at least one of the regular expressions.

##### Disable all notification by healthCheck

Add a KV entry with the key `consul-alerts/config/checks/blacklist/checks/{{ checkId }}`. This will disable notifications for the specified `checkId`.

##### Disable all notifications for the healthChecks matching regular expressions

Add a KV entry with the key `consul-alerts/config/checks/blacklist/checks` and the value containing a list of regular expressions. This will disable notifications for all the healthchecks, which names match at least one of the regular expressions.

##### Disable a single health check

Add a KV entry with the key `consul-alerts/config/checks/blacklist/single/{{ node }}/{{ serviceId }}/{{ checkId }}`. This will disable the specific health check. If the health check is not associated with a service, use the `_` as the serviceId.

### Events

Event handling is enabled by default. This delegates any consul event received by the agent to the list of handlers configured. To disable event handling, set `consul-alerts/config/events/enabled` to `false`.

Handlers can be configured by adding them to `consul-alerts/config/events/handlers`. This should be a JSON array of string. Each string should point to any executable. The event data should be read from `stdin`.

### Notifiers

There are several built-in notifiers. Only the *Log* notifier is enabled by default. Details on enabling and configuration these are documented for each Notifier.

#### Custom Notifiers
It is also possible to add custom notifiers similar to custom event handlers. Custom notifiers can be added as keys with command path string values in `consul-alerts/config/notifiers/custom/`. The keys will be used as notifier names in the profiles.

#### Logger

This logs any health check notification to a file. To disable this notifier, set `consul-alerts/config/notifiers/log/enabled` to `false`.

The log file is set to `/tmp/consul-notifications.log` by default. This can be changed by changing `consul-alerts/config/notifiers/log/path`.

#### Email

This emails the health notifications. To enable this, set `consul-alerts/config/notifiers/email/enabled` to `true`.

The email and smtp details needs to be configured:

prefix: `consul-alerts/config/notifiers/email/`

| key          | description                                                                      |
|--------------|----------------------------------------------------------------------------------|
| enabled      | Enable the email notifier. [Default: false]                                      |
| cluster-name | The name of the cluster. [Default: "Consul Alerts"]                              |
| url          | The SMTP server url                                                              |
| port         | The SMTP server port                                                             |
| username     | The SMTP username                                                                |
| password     | The SMTP password                                                                |
| sender-alias | The sender alias. [Default: "Consul Alerts"]                                     |
| sender-email | The sender email                                                                 |
| receivers    | The emails of the receivers. JSON array of string                                |
| template     | Path to custom email template. [Default: internal template]                      |
| one-per-alert| Whether to send one email per alert [Default: false]                             |
| one-per-node | Whether to send one email per node [Default: false] (overriden by one-per-alert) |

The template can be any go html template. An `TemplateData` instance will be passed to the template.

#### InfluxDB

This sends the notifications as series points in influxdb. Set `consul-alerts/config/notifiers/influxdb/enabled` to `true` to enabled. InfluxDB details need to be set too.

prefix: `consul-alerts/config/notifiers/influxdb/`

| key         | description                                    |
|-------------|------------------------------------------------|
| enabled     | Enable the influxdb notifier. [Default: false] |
| host        | The influxdb host. (eg. localhost:8086)        |
| username    | The influxdb username                          |
| password    | The influxdb password                          |
| database    | The influxdb database name                     |
| series-name | The series name for the points                 |

#### Slack

Slack integration is also supported. To enable, set
`consul-alerts/config/notifiers/slack/enabled` to `true`. Slack details needs to
be configured.

prefix: `consul-alerts/config/notifiers/slack/`

| key          | description                                                                                |
|--------------|-----------------------------------------------------                                       |
| enabled      | Enable the Slack notifier. [Default: false]                                                |
| cluster-name | The name of the cluster. [Default: `Consul Alerts`]                                        |
| url          | The incoming-webhook url (mandatory) [eg: `https://hooks.slack.com...`]                    |
| channel      | The channel to post the notification (mandatory) [eg: `#consul-alerts` or `@consul-alerts`]|
| username     | The username to appear on the post [eg: `Consul Alerts`]                                   |
| icon-url     | URL of a custom image for the notification [eg: `http://someimage.com/someimage.png`]      |
| icon-emoji   | Emoji (if not using icon-url) for the notification [eg: `:ghost:`]                         |
| detailed     | Enable "pretty" Slack notifications [Default: false]                                       |

In order to enable slack integration, you have to create a new
[_Incoming WebHooks_](https://my.slack.com/services/new/incoming-webhook). Then use the
token created by the previous action.

#### Mattermost

Mattermost integration is also supported. To enable, set
`consul-alerts/config/notifiers/mattermost/enabled` to `true`. Mattermost details needs to
be configured.

prefix: `consul-alerts/config/notifiers/mattermost/`

| key          | description                                         |
|--------------|-----------------------------------------------------|
| enabled      | Enable the Mattermost notifier. [Default: false]    |
| cluster-name | The name of the cluster. [Default: "Consul Alerts"] |
| url          | The mattermost url (mandatory)                      |
| username     | The mattermost username (mandatory)                 |
| password     | The mattermost password (mandatory)                 |
| team         | The mattermost team (mandatory)                     |
| channel      | The channel to post the notification (mandatory)    |

#### PagerDuty

To enable PagerDuty built-in notifier, set `consul-alerts/config/notifiers/pagerduty/enabled` to `true`. This is disabled by default. Service key and client details also needs to be configured.

prefix: `consul-alerts/config/notifiers/pagerduty/`

| key         | description                                     |
|-------------|-------------------------------------------------|
| enabled     | Enable the PagerDuty notifier. [Default: false] |
| service-key | Service key to access PagerDuty                 |
| client-name | The monitoring client name                      |
| client-url  | The monitoring client url                       |

#### HipChat

To enable HipChat built-in notifier, set
`consul-alerts/config/notifiers/hipchat/enabled` to `true`. Hipchat details
needs to be configured.

prefix: `consul-alerts/config/notifiers/hipchat/`

| key          | description                                               |
|--------------|-----------------------------------------------------------|
| enabled      | Enable the HipChat notifier. [Default: false]             |
| from         | The name to send notifications as  (optional)             |
| cluster-name | The name of the cluster. [Default: "Consul Alerts"]       |
| base-url     | HipChat base url [Default: `https://api.hipchat.com/v2/`] |
| room-id      | The room to post to                  (mandatory)          |
| auth-token   | Authentication token                 (mandatory)          |

The `auth-token` needs to be a room notification token for the `room-id`
being posted to.
See [HipChat API docs](https://developer.atlassian.com/hipchat/guide/hipchat-rest-api).

The default `base-url` works for HipChat-hosted rooms. You only need to
override it if you are running your own server.

#### OpsGenie

To enable OpsGenie built-in notifier, set
`consul-alerts/config/notifiers/opsgenie/enabled` to `true`. OpsGenie details
needs to be configured.

prefix: `consul-alerts/config/notifiers/opsgenie/`

| key          | description                                         |
|--------------|-----------------------------------------------------|
| enabled      | Enable the OpsGenie notifier. [Default: false]      |
| cluster-name | The name of the cluster. [Default: "Consul Alerts"] |
| api-key      | API Key                              (mandatory)    |

#### Amazon Web Services Simple Notification Service ("SNS")

To enable AWS SNS built-in notifier, set
`consul-alerts/config/notifiers/awssns/enabled` to `true`. AWS SNS details
needs to be configured.

prefix: `consul-alerts/config/notifiers/awssns/`

| key          | description                                                  |
|--------------|--------------------------------------------------------------|
| enabled      | Enable the AWS SNS notifier.   [Default: false]              |
| cluster-name | The name of the cluster.       [Default: "Consul Alerts"]    |
| region       | AWS Region                     (mandatory)                   |
| topic-arn    | Topic ARN to publish to.       (mandatory)                   |
| template     | Path to custom template.       [Default: internal template]  |
#### VictorOps

To enable the VictorOps built-in notifier, set
`consul-alerts/config/notifiers/victorops/enabled` to `true`. VictorOps details
needs to be configured.

prefix: `consul-alerts/config/notifiers/victorops/`

| key          | description                                         |
|--------------|-----------------------------------------------------|
| enabled      | Enable the VictorOps notifier. [Default: false]     |
| api-key      | API Key                              (mandatory)    |
| routing-key  | Routing Key                          (mandatory)    |


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

`http://consul-alerts:9000/v1/health/wildcard?node=<node>&service=<serviceId>&check=<checkId>&status=<status>&alwaysOk=true&ignoreBlacklist=true`
`v1/health/wildcard` is similar to `v1/health` but returns all matched checks (omitted service/node/check params assumed as any) . Values returned in JSON form, status code 503 if one of services in critical state.

Additional params are ignoreBlacklist and alwaysOk which forces status code to 200 regardless of checks status.

## Operations
Configuration may be set manually through consul UI or API, using configuration management tools such as [chef](https://github.com/dpetzel/consul_alerts-cookbook), puppet or Ansible, or backed up and restored using [consulate](https://github.com/gmr/consulate).

Consulate Example:
```
consulate kv backup consul-alerts/config -f consul-alerts-config.json
consulate kv restore consul-alerts/config -f consul-alerts-config.json --prune
```

Contribution
------------

PRs are more than welcome. Just fork, create a feature branch, and open a PR. We love PRs. :)

TODO
----

Needs better doc and some cleanup too. :)
