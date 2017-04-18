gopherduty [![Build Status](https://travis-ci.org/darkcrux/gopherduty.png)](https://travis-ci.org/darkcrux/gopherduty)
==========



A simple Go client for PagerDuty. This includes a retry feature when sending to PagerDuty

# Usage

#### Get library
```
$ go get github.com/darkcrux/gopherduty
```

#### Use library
```
import "github.com/darkcrux/gopherduty"
```

#### Create client
```
client := gopherduty.NewClient("e93facc04764012d7bfb002500d5d1a6")
```

#### Configure client
```
client.MaxRetries = 5 // set max retries to 5 before failing, Defaults to 0.
client.RetryBaseInterval = 5 // set first retry to 5s. Defaults to 10s.
```

#### Trigger an incident
```
response := client.Trigger("check-01", "something failed", "my-monitoring-client", "http://my.url.com", details)
```

#### Acknowledge an incident
```
response := client.Acknowledge("check1", "haxxor is fixing it naw", details)
```

#### Resolve an incident
```
response := client.Resolve("check1", "haxxor has fixed. Can haxxor has cheezburger", details)
```

#### Verify response
```
response.HasErrors() // true if there were errors even after all the retries. :(
response.Status // the status code
response.Message // the return message
response.IncidentKey // the incident key of the request
response.Errors // list of errors
```

# More Info

More info can be found [here](http://godoc.org/github.com/darkcrux/gopherduty).
