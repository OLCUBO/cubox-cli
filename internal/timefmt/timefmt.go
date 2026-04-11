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

// Parse converts a human-friendly time string to the Cubox API format.
//
// Accepted inputs:
//
//	"today"                            → start of today
//	"yesterday"                        → start of yesterday
//	"7d" / "-7d"                       → 7 days ago at 00:00
//	"2026-04-09"                       → date only, midnight local time
//	"2026-04-09 17:12"                 → date + hour:minute
//	"2026-04-09 17:12:48"              → date + time
//	"2026-04-09T17:12:48"              → ISO-style date + time
//	"2026-04-09T17:12:48.695+0800"     → full precision, passed through
//	"2026-04-09T17:12:48.695+08:00"    → full precision with colon tz
//
// Returns the original string unchanged if it already looks like a full
// API timestamp (contains 'T' and a timezone offset).
func Parse(input string) (string, error) {
	input = strings.TrimSpace(input)
	if input == "" {
		return "", nil
	}

	now := time.Now()

	switch strings.ToLower(input) {
	case "today":
		return formatMidnight(now), nil
	case "yesterday":
		return formatMidnight(now.AddDate(0, 0, -1)), nil
	case "now":
		return now.Format(APILayout), nil
	}

	if m := relDaysRe.FindStringSubmatch(input); m != nil {
		days, _ := strconv.Atoi(m[1])
		return formatMidnight(now.AddDate(0, 0, -days)), nil
	}

	tryLayouts := []string{
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

	for _, layout := range tryLayouts {
		if t, err := time.ParseInLocation(layout, input, now.Location()); err == nil {
			return t.Format(APILayout), nil
		}
	}

	return "", fmt.Errorf("unrecognized time format: %q\nAccepted: today, yesterday, 7d, 2006-01-02, 2006-01-02 15:04:05, or full ISO timestamp", input)
}

func formatMidnight(t time.Time) string {
	y, m, d := t.Date()
	return time.Date(y, m, d, 0, 0, 0, 0, t.Location()).Format(APILayout)
}
