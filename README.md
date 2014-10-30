consul-alerts
=============

A simple daemon to send notifications based on Consul health checks.

## Requirement

1. Consul 0.4+. Get it [here](http://consul.io).
2. Configured `GOPATH`.

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
$ consul-alerts start --alert-addr=localhost:9000 --consul-addr=localhost:8500 --consul-dc=dc1
```

Note: Don't change --alert-addr when using the docker container.

Once the daemon is running, it can act as a handler for consul watches. At the moment only checks and events are supported.

```
$ consul watch -type checks consul-alerts watch checks [--alert-addr=localhost:9000]
$ consul watch -type event consul-alerts watch event [--alert-addr=localhost:9000]
```

or run the watchers on the agent the daemon connects by adding the following flags during consul-alerts run:

```
$ consul-alerts start --watch-events --watch-checks
```

Configuration
-------------

All configurations are stored in consul's KV with the prefix: `consul-alerts/config/`. The daemon is using default values and the KV entries will only override the defaults.

### Health Checks

Health checking is enabled by default. This also triggers the notification when a check has changed status for a configured duration. Health checks can be disabled by setting the kv`consul-alert/config/checks/enabled` to `false`.

To prevent flapping, notifications are only sent when a check status has been stable for a specific time in seconds (60 by default). this value can be changed by adding/changing the kv `consul-alert/config/checks/change-threshold` to an integer greater than and divisible by 10.

eg. `consul-alert/config/checks/change-threshold` = `30`

### Events

Event handling is enabled by default. This delegates any consul event received by the agent to the list of handlers configured. To disable event handling, set `consul-alert/config/events/enabled` to `false`.

Handlers can be configured by adding them to `consul-alert/config/events/handlers`. This should be a JSON array of string. Each string should point to any executable. The event data should be read from `stdin`.

### Notifiers

There are four builtin notifiers. Only the *Log* notifier is enabled by default. It is also possible to add custom notifiers similar to custom event handlers. Custom notifiers can be added in `consul-alerts/config/notifiers/custom`.

#### Logger

This logs any health check notification to a file. To disable this notifier, set `consul-alert/config/notifiers/log/enabled` to `false`.

The log file is set to `/tmp/consul-notifications.log` by default. This can be changed by changing `consul-alert/config/notifiers/log/path`.

#### Email

This emails the health notifications. To enable this, set `consul-alert/config/notifiers/email/enabled` to `true`.

The email and smtp details needs to be configured:

prefix: `consul-alert/config/notifiers/email/`

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

This sends the notifications as series points in influxdb. Set `consul-alert/config/notifiers/influxdb/enabled` to `true` to enabled. InfluxDB details need to be set too.

prefix: `consul-alert/config/notifiers/influxdb/`

| key         | description                                    |
|-------------|------------------------------------------------|
| enabled     | Enable the influxdb notifier. [Default: false] |
| host        | The influxdb host. (eg. localhost:8086)        |
| username    | The influxdb username                          |
| password    | The influxdb password                          |
| database    | The influxdb database name                     |
| series-name | The series name for the points                 |

#### Slack

Slack integration is also supported. To enable, set `consul-alert/config/notifiers/slack/enabled` to `true`. Slack details needs to be configured.

prefix: `consul-alert/config/notifiers/slack/`

| key          | description                                         |
|--------------|-----------------------------------------------------|
| enabled      | Enable the Slack notifier. [Default: false]         |
| cluster-name | The name of the cluster. [Default: "Consul Alerts"] |
| team         | The slack team name                                 |
| token        | The token for the webhook api                       |
| channel      | The channel to post the notification                |
| username     | The username to appear on the post                  |
| icon-url     | URL of a custom image for the notification          |
| icon-emoji   | Emoji (if not using icon-url) for the notification  |

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

This is a port from a tool we developed recently, there are still a few things missing like loading a custom configuration via command/api instead of manually editing consul's KV. Also need to set up a reminder feature. Needs better doc and some cleanup too. :)
