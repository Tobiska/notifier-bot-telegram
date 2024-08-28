package models

import "time"

type Detail struct {
	ID           int64
	Name         string
	softDeadline time.Time
	hardDeadline time.Time
}
