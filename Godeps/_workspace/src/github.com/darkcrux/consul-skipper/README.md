consul-skipper
==============

[![GoDoc](https://godoc.org/github.com/darkcrux/consul-skipper?status.png)](https://godoc.org/github.com/darkcrux/consul-skipper)
[![Build Status](https://travis-ci.org/darkcrux/consul-skipper.png)](https://travis-ci.org/darkcrux/consul-skipper)


consul-skipper is a library for a cluster leader election using Consul KV. This needs to be attached to the Consul Agent
and then it runs leader election from there.

```
import "github.com/darkcrux/consul-skipper"

candidate := &skipper.Candidate{
    ConsulAddress:    "10.0.0.10:8500",
    ConsulDatacenter: "dc1",
    LeadershipKey:    "app/leader",
    ConsulAclToken:   "",               // Optional
}
candidate.RunForElection()
```

Running for election runs asynchronously and will keep running as long as the main application is running. To check if
the current attached agent is the current leader, use:

```
skipper.IsLeader()
```

It is also possible to force a leader to step down forcing a re-election.

```
skipper.Resign()
```
