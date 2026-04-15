package timefmt

import (
	"strings"
	"testing"
	"time"
)

func TestParseStartOfDay(t *testing.T) {
	out, err := Parse("today")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, "T00:00:00.000") {
		t.Errorf("Parse(today) should be midnight, got %s", out)
	}
}

func TestParseEndOfDay(t *testing.T) {
	out, err := ParseEnd("today")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, "T23:59:59.999") {
		t.Errorf("ParseEnd(today) should be end-of-day, got %s", out)
	}
}

func TestParseYesterday(t *testing.T) {
	start, err := Parse("yesterday")
	if err != nil {
		t.Fatal(err)
	}
	end, err := ParseEnd("yesterday")
	if err != nil {
		t.Fatal(err)
	}

	yesterday := time.Now().AddDate(0, 0, -1).Format("2006-01-02")
	if !strings.HasPrefix(start, yesterday+"T00:00:00.000") {
		t.Errorf("Parse(yesterday) wrong date prefix, got %s", start)
	}
	if !strings.HasPrefix(end, yesterday+"T23:59:59.999") {
		t.Errorf("ParseEnd(yesterday) wrong date prefix, got %s", end)
	}
}

func TestParseRelativeDays(t *testing.T) {
	start, err := Parse("7d")
	if err != nil {
		t.Fatal(err)
	}
	end, err := ParseEnd("7d")
	if err != nil {
		t.Fatal(err)
	}

	day7ago := time.Now().AddDate(0, 0, -7).Format("2006-01-02")
	if !strings.HasPrefix(start, day7ago+"T00:00:00.000") {
		t.Errorf("Parse(7d) = %s, want prefix %sT00:00:00.000", start, day7ago)
	}
	if !strings.HasPrefix(end, day7ago+"T23:59:59.999") {
		t.Errorf("ParseEnd(7d) = %s, want prefix %sT23:59:59.999", end, day7ago)
	}
}

func TestParseDateOnly(t *testing.T) {
	start, err := Parse("2026-03-15")
	if err != nil {
		t.Fatal(err)
	}
	end, err := ParseEnd("2026-03-15")
	if err != nil {
		t.Fatal(err)
	}

	if !strings.HasPrefix(start, "2026-03-15T00:00:00.000") {
		t.Errorf("Parse(date) = %s, want midnight", start)
	}
	if !strings.HasPrefix(end, "2026-03-15T23:59:59.999") {
		t.Errorf("ParseEnd(date) = %s, want end-of-day", end)
	}
}

func TestParseNowIdenticalForBoth(t *testing.T) {
	s, _ := Parse("now")
	e, _ := ParseEnd("now")
	if s[:19] != e[:19] {
		t.Errorf("now should be ~identical for Parse/ParseEnd, got %s vs %s", s, e)
	}
}

func TestParseExplicitTimestampUnchanged(t *testing.T) {
	ts := "2026-04-15T17:59:57.650+0800"
	s, _ := Parse(ts)
	e, _ := ParseEnd(ts)
	if s != e {
		t.Errorf("explicit timestamp should be identical: %s vs %s", s, e)
	}
	if s != ts {
		t.Errorf("explicit timestamp should round-trip: got %s, want %s", s, ts)
	}
}

func TestParseEmpty(t *testing.T) {
	s, err := Parse("")
	if err != nil || s != "" {
		t.Errorf("empty input should return empty: got %q, err=%v", s, err)
	}
	e, err := ParseEnd("")
	if err != nil || e != "" {
		t.Errorf("empty input should return empty: got %q, err=%v", e, err)
	}
}
