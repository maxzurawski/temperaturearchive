package crosscutting

import "time"

type SearchDTO struct {
	Uuid           string     `json:"uuid"`
	ProcessId      string     `json:"processId"`
	ValueFrom      *float64   `json:"valueFrom"`
	ValueTo        *float64   `json:"valueTo"`
	ReportedAtFrom *time.Time `json:"reportedAtFrom"`
	ReportedAtTo   *time.Time `json:"reportedAtTo"`
	LastLimited    *int       `json:"lastLimited"`
	OrderDesc      *bool      `json:"orderDesc"`
}
