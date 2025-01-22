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
	return `
		You will be need to provide a duration to sleep for. The duration will be in seconds. It should be a whole number,
		between 1 and 10.

		Do not include quotes in your response.
	`
}

type sleepRequest struct {
	Duration time.Duration `json:"duration"`
}

func (s *Sleeper) Use(ctx context.Context, req any) error {
	reqObj, ok := req.(sleepRequest)
	if !ok {
		return fmt.Errorf("expected duration to be a time.Duration")
	}

	fmt.Println("Sleeping for", reqObj.Duration)
	time.Sleep(reqObj.Duration)

	return nil
}
