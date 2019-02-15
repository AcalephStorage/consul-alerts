package notificationv2

// Repeat defines the amount of time in minutes that notification steps will be repeatedly apply.
type Repeat struct {
	LoopAfter int  `json:"loopAfter"`
	Enabled   bool `json:"enabled"`
}
