package timefmt

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// APILayout is the time format the Cubox API expects.
const APILayout = "2006-01-02T15:04:05.000-0700"

var relDaysRe = regexp.MustCompile(`^-?(\d+)d$`)

// dateOnlyLayouts match inputs that specify only a date (no hour/minute).
var dateOnlyLayouts = []string{
	"2006-01-02",
}

var tryLayouts = []string{
	"2006-01-02",
	"2006-01-02 15:04",
	"2006-01-02 15:04:05",
	"2006-01-02T15:04:05",
	"2006-01-02T15:04:05.000-0700",
	"2006-01-02T15:04:05.000-07:00",
	"2006-01-02T15:04:05-0700",
	"2006-01-02T15:04:05-07:00",
	"2006-01-02T15:04:05.000Z",
}

// Parse converts a human-friendly time string to the Cubox API format.
// Day-level shortcuts (today, yesterday, 7d, date-only) resolve to the
// start of the day (00:00:00.000). Use ParseEnd for end-time semantics.
//
// Accepted inputs:
//
//	"today"                            → start of today (00:00:00.000)
//	"yesterday"                        → start of yesterday
//	"now"                              → current instant
//	"7d" / "-7d"                       → 7 days ago at 00:00
//	"2026-04-09"                       → date only, midnight local time
//	"2026-04-09 17:12"                 → date + hour:minute
//	"2026-04-09 17:12:48"              → date + time
//	"2026-04-09T17:12:48"              → ISO-style date + time
//	"2026-04-09T17:12:48.695+0800"     → full precision
//	"2026-04-09T17:12:48.695+08:00"    → full precision with colon tz
func Parse(input string) (string, error) {
	return parse(input, false)
}

// ParseEnd is like Parse but resolves day-level inputs to end-of-day
// (23:59:59.999) instead of midnight. Use this for --end-time parameters
// so that "today" covers the entire day rather than just its first instant.
func ParseEnd(input string) (string, error) {
	return parse(input, true)
}

func parse(input string, endOfDay bool) (string, error) {
	input = strings.TrimSpace(input)
	if input == "" {
		return "", nil
	}

	now := time.Now()
	snap := snapFn(endOfDay)

	switch strings.ToLower(input) {
	case "today":
		return snap(now), nil
	case "yesterday":
		return snap(now.AddDate(0, 0, -1)), nil
	case "now":
		return now.Format(APILayout), nil
	}

	if m := relDaysRe.FindStringSubmatch(input); m != nil {
		days, _ := strconv.Atoi(m[1])
		return snap(now.AddDate(0, 0, -days)), nil
	}

	if endOfDay {
		for _, layout := range dateOnlyLayouts {
			if t, err := time.ParseInLocation(layout, input, now.Location()); err == nil {
				return formatEndOfDay(t), nil
			}
		}
	}

	for _, layout := range tryLayouts {
		if t, err := time.ParseInLocation(layout, input, now.Location()); err == nil {
			return t.Format(APILayout), nil
		}
	}

	return "", fmt.Errorf("unrecognized time format: %q\nAccepted: today, yesterday, 7d, 2006-01-02, 2006-01-02 15:04:05, or full ISO timestamp", input)
}

func snapFn(endOfDay bool) func(time.Time) string {
	if endOfDay {
		return formatEndOfDay
	}
	return formatMidnight
}

func formatMidnight(t time.Time) string {
	y, m, d := t.Date()
	return time.Date(y, m, d, 0, 0, 0, 0, t.Location()).Format(APILayout)
}

func formatEndOfDay(t time.Time) string {
	y, m, d := t.Date()
	return time.Date(y, m, d, 23, 59, 59, 999_000_000, t.Location()).Format(APILayout)
}
