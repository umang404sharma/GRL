package models

type HostReport struct {
	Host string `json:"host"`
	RPS  int64  `json:"rps"`
}
