package responseMessage

type ResponseUnderlyingList struct {
	RequestID      string         `json:"request_id"`
	Name           string         `json:"name"`
	UnderlyingData UnderlyingData `json:"msg"`
	Status         int            `json:"status"`
}
type Schedule struct {
	Open  int `json:"open"`
	Close int `json:"close"`
}
type Tags struct {
}
type Underlying struct {
	ActiveID        int        `json:"active_id"`
	ActiveGroupID   int        `json:"active_group_id"`
	ActiveType      string     `json:"active_type"`
	Underlying      string     `json:"underlying"`
	Schedule        []Schedule `json:"schedule"`
	IsEnabled       bool       `json:"is_enabled"`
	Name            string     `json:"name"`
	LocalizationKey string     `json:"localization_key"`
	Image           string     `json:"image"`
	ImagePrefix     string     `json:"image_prefix"`
	Precision       int        `json:"precision"`
	StartTime       int64      `json:"start_time"`
	RegulationMode  string     `json:"regulation_mode"`
	Tags            Tags       `json:"tags"`
	IsSuspended     bool       `json:"is_suspended"`
}
type UnderlyingData struct {
	Type        string       `json:"type"`
	UserGroupID int          `json:"user_group_id"`
	Underlying  []Underlying `json:"underlying"`
}
