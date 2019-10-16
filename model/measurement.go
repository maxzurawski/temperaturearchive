package model

import "time"

type Measurement struct {
	ID         *uint      `gorm:"primary_key"`
	Uuid       *string    `gorm:"column:uuid"`
	ProcessId  *string    `gorm:"column:process_id"`
	Value      *float64   `gorm:"column:value"`
	ReportedAt *time.Time `gorm:"column:reported_at"` // when it was received in dispatcher
	ReceivedAt *time.Time `gorm:"column:received_at"` // when it was received in temperaturearchive
}
