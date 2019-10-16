package dto

import "time"

type MeasurementDTO struct {
	ID         *uint     `json:"id"`
	Uuid       string    `json:"uuid"`
	ProcessId  string    `json:"processId"`
	Value      float64   `json:"value"`
	ReportedAt time.Time `json:"reportedAt"`
	ReceivedAt time.Time `json:"receivedAt"`
}
