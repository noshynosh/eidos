package sleep

import (
	"context"
	"fmt"
	"time"
)

type Sleeper struct{}

func (s *Sleeper) Description() string {
	return "Sleep for a given amount of time"
}

func (s *Sleeper) Manual() string {
	return `{"indicator": "sleep", "duration": "<duration>", "description": "Sleep for a given amount of time where <duration> is in seconds"}`
}

type sleepRequest struct {
	Duration time.Duration `json:"duration"`
}

func (s *Sleeper) Use(ctx context.Context, req any) (string, error) {
	reqObj, ok := req.(sleepRequest)
	if !ok {
		return "", fmt.Errorf("expected duration to be a time.Duration")
	}

	time.Sleep(reqObj.Duration)

	return "Slept for " + reqObj.Duration.String(), nil
}
