package models

import "time"

type Detail struct {
	ID           int64
	ChatID       int64
	Name         string
	SoftDeadline *time.Time
	HardDeadline *time.Time
}

func (d *Detail) IsFilled() bool {
	return d.SoftDeadline != nil && d.HardDeadline != nil
}
