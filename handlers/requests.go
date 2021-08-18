package handlers

// FetchDataReq represents a HTTP request body for Fetch Data requests.
type FetchDataReq struct {
	StartDate string `json:"startDate"`
	EndDate   string `json:"endDate"`
	MinCount  int    `json:"minCount"`
	MaxCount  int    `json:"maxCount"`
}

// Validate validates given request.
func (f *FetchDataReq) Validate() error {
	if f.StartDate == "" || f.EndDate == "" {
		return ErrMalformedEntity
	}
	return nil
}

// InMemoryCreateReq represents a HTTP request body for in-memory
// data creation through HTTP.
type InMemoryCreateReq struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}
