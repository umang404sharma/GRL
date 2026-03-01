package models

type ZoneReport struct {
	Zone     string `json:"zone"`
	TotalRPS int64  `json:"total_rps"`
}
