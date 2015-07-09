package heartbeat

type AddHeartbeatResponse struct {
	Id		string	`json:"id"`
	Status	string	`json:"status"`
	Code	int		`json:"code"`
}

type UpdateHeartbeatResponse struct {
	Id		string	`json:"id"`
	Status	string	`json:"status"`
	Code	int		`json:"code"`
}

type EnableHeartbeatResponse struct {
	Status 	string	`json:"status"`	
	Code	int	`json:"code"`
}

type DisableHeartbeatResponse struct {
	Status 	string	`json:"status"`
	Code	int	`json:"code"`
}

type DeleteHeartbeatResponse struct {
	Status 	string	`json:"status"`
	Code	int	`json:"code"`
}

type GetHeartbeatResponse struct {
	Id				string	`json:"id"`
	Name			string	`json:"name"`
	Status			string	`json:"status"`
	Description		string	`json:"description"`
	Enabled 		bool	`json:"enabled"`
	LastHeartbeat	uint64	`json:"lastHeartBeat"`
	Interval		int		`json:"interval"`
	IntervalUnit	string	`json:"intervalUnit"`
	Expired			bool	`json:"expired"`
}

type ListHeartbeatsResponse struct {
	Heartbeats []struct {
			Id				string	`json:"id"`
			Name			string	`json:"name"`
			Status			string	`json:"status"`
			Description		string	`json:"description"`
			Enabled 		bool	`json:"enabled"`
			LastHeartbeat	uint64	`json:"lastHeartBeat"`
			Interval		int		`json:"interval"`
			IntervalUnit	string 	`json:"intervalUnit"`
			Expired			bool	`json:"expired"`
	}	`json:"heartbeats"`
}

type SendHeartbeatResponse struct {
	WillExpireAt 	uint64		`json:"willExpireAt"`
	Status			string 		`json:"status"`
	Heartbeat 		uint64 		`json:"heartbeat"`
	Took			int 		`json:"took"`	
	Code			int 		`json:"code"`	
}
