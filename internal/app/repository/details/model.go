package details

import "time"

type UpdateModel struct {
	Name         *string
	SoftDeadline *time.Time
	HardDeadline *time.Time
}

type Filter struct {
	ChatID   *int64
	Name     *string
	IsFilled *bool
}
