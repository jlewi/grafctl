package grafana

import (
	"regexp"
	"strconv"
	"time"

	"github.com/pkg/errors"
)

type Clock interface {
	Now() time.Time
}

type RealClock struct{}

func (RealClock) Now() time.Time {
	return time.Now()
}

func NewRelativeTimeParser() *RelativeTimeParser {
	return &RelativeTimeParser{Clock: RealClock{}}
}

type RelativeTimeParser struct {
	Clock Clock
}

// ParseGrafanaRelativeTime converts a Grafana-style relative time string to a time.Time object.
// Example inputs: "now", "now-1h", "now-30m", "now-7d"
func (p RelativeTimeParser) ParseGrafanaRelativeTime(relativeTime string) (time.Time, error) {
	if relativeTime == "" {
		return time.Time{}, errors.New("empty relative time; use 'now' for the current time")
	}
	// Get the current time
	now := p.Clock.Now()

	if relativeTime == "now" {
		return now, nil
	}

	// Define regex to match Grafana relative time format (e.g., now-1h, now-7d)
	re := regexp.MustCompile(`^now-(\d+)([smhdwMy])$`)
	matches := re.FindStringSubmatch(relativeTime)

	if len(matches) != 3 {
		return time.Time{}, errors.New("invalid relative time format")
	}

	// Extract the amount and unit from the matches
	amount, err := strconv.Atoi(matches[1])
	if err != nil {
		return time.Time{}, errors.New("invalid time amount")
	}

	unit := matches[2]

	// Calculate the offset based on the unit
	var duration time.Duration
	switch unit {
	case "s":
		duration = time.Duration(amount) * time.Second
	case "m":
		duration = time.Duration(amount) * time.Minute
	case "h":
		duration = time.Duration(amount) * time.Hour
	case "d":
		duration = time.Duration(amount) * 24 * time.Hour
	case "w":
		duration = time.Duration(amount) * 7 * 24 * time.Hour
	case "M":
		duration = time.Duration(amount) * 30 * 24 * time.Hour
	case "y":
		duration = time.Duration(amount) * 365 * 24 * time.Hour
	default:
		return time.Time{}, errors.New("unknown time unit")
	}

	// Subtract the duration from now
	return now.Add(-duration), nil
}
