package grafana

import (
	"testing"
	"time"
)

type FakeClock struct {
}

func (f FakeClock) Now() time.Time {
	return time.Date(2024, time.February, 25, 13, 25, 0, 0, time.UTC)
}

func TestParseGrafanaRelativeTime(t *testing.T) {
	clock := FakeClock{}
	tests := []struct {
		input         string
		expectedError bool
		expected      time.Time
	}{
		{"now", false, clock.Now()},
		{"now-1h", false, clock.Now().Add(-1 * time.Hour)},
		{"now-30m", false, clock.Now().Add(-30 * time.Minute)},
		{"now-7d", false, clock.Now().Add(-7 * 24 * time.Hour)},
		{"now-3M", false, clock.Now().Add(-3 * 30 * 24 * time.Hour)},
		{"now-1y", false, clock.Now().Add(-365 * 24 * time.Hour)},
		{"invalid", true, clock.Now()},
		{"now-5x", true, clock.Now()},
	}

	p := RelativeTimeParser{
		Clock: FakeClock{},
	}
	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			parsedTime, err := p.ParseGrafanaRelativeTime(test.input)

			if test.expectedError {
				if err == nil {
					t.Errorf("expected error but got none for input %s", test.input)
				}
				return
			}

			if err != nil {
				t.Errorf("did not expect error but got %v for input %s", err, test.input)
				return
			}

			if parsedTime != test.expected {
				t.Errorf("expected %v but got %v for input %s", test.expected, parsedTime, test.input)
			}
		})
	}
}
