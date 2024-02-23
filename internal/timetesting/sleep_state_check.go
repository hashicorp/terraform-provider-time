package timetesting

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/statecheck"
)

var _ statecheck.StateCheck = &sleep{}

type sleep struct {
	seconds int
}

func (s *sleep) CheckState(ctx context.Context, req statecheck.CheckStateRequest, resp *statecheck.CheckStateResponse) {
	time.Sleep(time.Duration(s.seconds) * time.Second)
}

// Sleep returns a state check that sleep for the given number of seconds. This state check can be used
// for certain tests that require time differences that are too fast for unit testing.
func Sleep(seconds int) statecheck.StateCheck {
	return &sleep{
		seconds: seconds,
	}
}
