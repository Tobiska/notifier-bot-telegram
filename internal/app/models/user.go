package models

import "time"

type Status string

const (
	Created         Status = "created"
	Wait            Status = "waiting"
	StartAdded             = "add_start"
	AddName         Status = "add_name"
	AddSoftDeadline Status = "add_soft_deadline"
	AddHardDeadline Status = "add_hard_deadline"
)

type User struct {
	ChatID    int64
	UserID    int64
	Status    Status
	Username  string
	Details   []Detail
	CreatedAt time.Time
	UpdatedAt time.Time
}
