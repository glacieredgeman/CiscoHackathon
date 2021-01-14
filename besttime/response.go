package besttime

type NewForecastResponse struct {
	Status    string    `json:"status"`
	VenueInfo VenueInfo `json:"venue_info"`
	Analysis  []Day     `json:"analysis"`
}

type VenueInfo struct {
	VenueId      string `json:"venue_id"`
	VenueName    string `json:"venue_name"`
	VenueAddress string `json:"venue_address"`
}

type Day struct {
	HourAnalysis []HourInfo `json:"hour_analysis"`
}

type HourInfo struct {
	Intensity string `json:"intensity_txt"`
}

type QuietHoursResponse struct {
	Analysis QuietHourAnalysis `json:"analysis"`
}

type QuietHourAnalysis struct {
	QuietHoursList []string `json:"quiet_hours_list_12h"`
}
