package utils

import (
	"time"
)

type Retry struct {
	countAttempts         int
	timeoutBetweenAttempt time.Duration
}

func NewRetry(countAttempts int, timeoutBetweenAttempt time.Duration) *Retry {
	return &Retry{
		countAttempts:         countAttempts,
		timeoutBetweenAttempt: timeoutBetweenAttempt,
	}
}

func (r *Retry) DoRetry(retryable func() error) error {
	var err error
	for i := 0; i < r.countAttempts; i++ {
		err = retryable()
		if err == nil {
			break
		}
		time.Sleep(r.timeoutBetweenAttempt)
	}
	return err
}
