package entities

import "time"

type MigrationStatus struct {
	Status    string    `json:"status"`
	Error     error     `json:"error"`
	StartTime time.Time `json:"startTime"`
	EndTime   time.Time `json:"endTime"`
	Total     int       `json:"total"`
	N         int       `json:"n"`
}
