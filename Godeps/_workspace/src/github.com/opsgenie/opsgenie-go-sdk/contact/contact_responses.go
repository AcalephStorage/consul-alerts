package contact

// Create contact response structure
type CreateContactResponse struct {
	Id     string `json:"id"`
	Status string `json:"status"`
	Code   int    `json:"code"`
}

// Update contact response structure
type UpdateContactResponse struct {
	Id     string `json:"id"`
	Status string `json:"status"`
	Code   int    `json:"code"`
}

// Delete contact response structure
type DeleteContactResponse struct {
	Status string `json:"status"`
	Code   int    `json:"code"`
}

// Disable contact response structure
type DisableContactResponse struct {
	Status string `json:"status"`
}

// Enable contact response structure
type EnableContactResponse struct {
	Status string `json:"status"`
}

// Get contact response structure
type GetContactResponse struct {
	Id             string `json:"id,omitempty"`
	Method         string `json:"method,omitempty"`
	To             string `json:"to,omitempty"`
	DisabledReason string `json:"disabledReason, omitempty"`
	Enabled        bool   `json:"enabled, omitempty"`
}
