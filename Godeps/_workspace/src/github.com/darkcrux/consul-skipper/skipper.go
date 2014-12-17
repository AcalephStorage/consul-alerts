/*
consul-skipper is a library for a cluster leader election using Consul KV. This needs to be attached to the Consul Agent
and then it runs leader election from there.

		import "github.com/darkcrux/consul-skipper"

		candidate := &skipper.Candidate{
			ConsulAddress:    "10.0.0.10:8500",
			ConsulDatacenter: "dc1",
			LeadershipKey:    "app/leader",
		}
		candidate.RunForElection()

Running for election runs asynchronously and will keep running as long as the main application is running. To check if
the current attached agent is the current leader, use:

		skipper.IsLeader()

It is also possible to force a leader to step down forcing a re-election.

		skipper.Resign()

*/
package skipper

import (
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/armon/consul-api"
)

type Candidate struct {
	ConsulAddress    string // The address of the consul agent. This defaults to 127.0.0.1:8500.
	ConsulDatacenter string // The datacenter to connect to. This defaults to the config used by the agent.
	LeadershipKey    string // The leadership key. This needs to be a proper Consul KV key. eg. app/leader
	session          string
	node             string
}

// RunForElection makes the candidate run for leadership against the other nodes in the cluster. This tries and acquire
// the lock of the given LeadershipKey and setting it's value with the current node. If the lock acquisition passes,
// then the node where the agent is running is now the new leader. A re-election will occur when there are changes in
// the LeadershipKey.
func (c *Candidate) RunForElection() {
	go c.campaign()
}

// IsLeader returns true if the current agent is the leader.
func (c *Candidate) IsLeader() bool {
	consul := c.consulClient()
	c.retrieveNode()
	c.retrieveSession()
	kv, _, err := consul.KV().Get(c.LeadershipKey, nil)
	if err != nil {
		logrus.Errorln("Unable to check for leadership:", err)
		return false
	}
	if kv == nil {
		logrus.Warnf("Leadership key '%s' is missing in Consuk KV.", c.LeadershipKey)
		return false
	}
	return c.node == string(kv.Value) && c.session == kv.Session
}

// Leader returns the node of the current cluster leader. This returns an empty string if there is no leader.
func (c *Candidate) Leader() string {
	consul := c.consulClient()
	kv, _, err := consul.KV().Get(c.LeadershipKey, nil)
	if kv == nil || err != nil {
		logrus.Warnln("There is no leader.")
		return ""
	}
	return string(kv.Value)
}

// Resign forces the current agent to step down as the leader forcing a re-election. Nothing happens if the agent is not
// the current leader.
func (c *Candidate) Resign() {
	if c.IsLeader() {
		consul := c.consulClient()
		kvpair := &consulapi.KVPair{
			Key:     c.LeadershipKey,
			Value:   []byte(c.node),
			Session: c.session,
		}
		success, _, err := consul.KV().Release(kvpair, nil)
		if !success || err != nil {
			logrus.Warnf("%s was unable to step down as a leader", c.node)
		} else {
			logrus.Debugf("%s is no longer the leader.", c.node)
		}
	}
}

// Campaign handles leader election. Basically this just acquires a lock on the LeadershipKey and whoever gets the lock
// is the leader. A re-election occurs when there are changes in the LeadershipKey.
func (c *Candidate) campaign() {
	c.retrieveNode()
	c.retrieveSession()
	consul := c.consulClient()

	logrus.Debugf("%s is running for election with session %s.", c.node, c.session)

	kvpair := &consulapi.KVPair{
		Key:     c.LeadershipKey,
		Value:   []byte(c.node),
		Session: c.session,
	}
	acquired, _, err := consul.KV().Acquire(kvpair, nil)
	if err != nil {
		logrus.Errorln("Failed to run Consul KV Acquire:", err)
	}

	if acquired {
		logrus.Infof("%s has become the leader.", c.node)
	}

	kv, _, _ := consul.KV().Get(c.LeadershipKey, nil)

	if kv != nil && kv.Session != "" {
		logrus.Debugf("%s is the current leader.", string(kv.Value))
		logrus.Debugf("%s is waiting for changes in '%s'.", c.node, c.LeadershipKey)
		latestIndex := kv.ModifyIndex
		options := &consulapi.QueryOptions{
			WaitIndex: latestIndex,
		}
		consul.KV().Get(c.LeadershipKey, options)
	}
	time.Sleep(15 * time.Second)
	c.campaign()
}

// RetrieveNode is a helper to retrieve the current node name of the agent.
func (c *Candidate) retrieveNode() {
	consul := c.consulClient()
	agent, err := consul.Agent().Self()
	if err != nil {
		logrus.Warnln("Unable to retrieve node name.")
	}
	c.node = agent["Config"]["NodeName"].(string)
}

// RetrieveSession retrieves the existing session needed to run leader election. If a session does not exist, a new
// session is created with the LeadershipKey as the name.
func (c *Candidate) retrieveSession() {
	consul := c.consulClient()

	if sessions, _, err := consul.Session().List(nil); err != nil {
		logrus.Warnln("Unable to retrieve list of sessions.")
	} else {
		for _, session := range sessions {
			if session.Name == c.LeadershipKey && session.Node == c.node {
				c.session = session.ID
				return
			}
		}
	}

	newSession := &consulapi.SessionEntry{
		Name: c.LeadershipKey,
	}
	if sessionId, _, err := consul.Session().Create(newSession, nil); err != nil {
		logrus.Errorln("Unable to create new sessions:", err)
	} else {
		c.session = sessionId
	}
}

// ConsulClient is a helper to create the consulapi client for access to the Consul cluster.
func (c *Candidate) consulClient() *consulapi.Client {
	config := consulapi.DefaultConfig()
	if c.ConsulAddress != "" {
		config.Address = c.ConsulAddress
	}
	if c.ConsulDatacenter != "" {
		config.Datacenter = c.ConsulDatacenter
	}
	client, _ := consulapi.NewClient(config)
	return client
}
