package recurrence

import (
	"testing"
	"time"
)

func TestParsePattern(t *testing.T) {
	tests := []struct {
		input   string
		want    string
		wantErr bool
	}{
		// Simple patterns
		{"daily", "daily", false},
		{"weekly", "weekly", false},
		{"monthly", "monthly", false},
		{"yearly", "yearly", false},
		{"Daily", "daily", false},
		{"WEEKLY", "weekly", false},

		// Interval patterns
		{"every 2d", "every 2d", false},
		{"every 3w", "every 3w", false},
		{"every 6m", "every 6m", false},
		{"every 1y", "every 1y", false},
		{"every 14d", "every 14d", false},

		// Day-of-week patterns
		{"every monday", "every monday", false},
		{"every mon", "every mon", false},
		{"every mon,wed,fri", "every mon,wed,fri", false},
		{"every tuesday", "every tuesday", false},

		// Invalid patterns
		{"", "", true},
		{"biweekly", "", true},
		{"every", "", true},
		{"every 0d", "", true},
		{"every -1w", "", true},
		{"every funday", "", true},
		{"every 2x", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got, err := ParsePattern(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParsePattern(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ParsePattern(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestNextDueDate(t *testing.T) {
	// Use a fixed "today" by testing with dates far in the future
	// so advanceByInterval stops at the first step
	date := func(y, m, d int) time.Time {
		return time.Date(y, time.Month(m), d, 0, 0, 0, 0, time.Local)
	}

	tests := []struct {
		name       string
		pattern    string
		currentDue time.Time
		wantAfter  time.Time // result should be >= this
	}{
		{
			name:       "daily from future date",
			pattern:    "daily",
			currentDue: date(2099, 1, 1),
			wantAfter:  date(2099, 1, 2),
		},
		{
			name:       "weekly from future date",
			pattern:    "weekly",
			currentDue: date(2099, 1, 1),
			wantAfter:  date(2099, 1, 8),
		},
		{
			name:       "monthly from future date",
			pattern:    "monthly",
			currentDue: date(2099, 1, 15),
			wantAfter:  date(2099, 2, 15),
		},
		{
			name:       "yearly from future date",
			pattern:    "yearly",
			currentDue: date(2099, 6, 15),
			wantAfter:  date(2100, 6, 15),
		},
		{
			name:       "every 2w from future date",
			pattern:    "every 2w",
			currentDue: date(2099, 1, 1),
			wantAfter:  date(2099, 1, 15),
		},
		{
			name:       "every 3m from future date",
			pattern:    "every 3m",
			currentDue: date(2099, 1, 1),
			wantAfter:  date(2099, 4, 1),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NextDueDate(tt.pattern, tt.currentDue)
			if err != nil {
				t.Fatalf("NextDueDate(%q, %v) error = %v", tt.pattern, tt.currentDue, err)
			}
			if got.Before(tt.wantAfter) {
				t.Errorf("NextDueDate(%q, %v) = %v, want >= %v", tt.pattern, tt.currentDue, got, tt.wantAfter)
			}
		})
	}
}

func TestNextDueDateAdvancesPastToday(t *testing.T) {
	// Use a past date - should advance past today
	pastDate := time.Date(2020, 1, 1, 0, 0, 0, 0, time.Local)
	today := time.Now()
	today = time.Date(today.Year(), today.Month(), today.Day(), 0, 0, 0, 0, today.Location())

	got, err := NextDueDate("weekly", pastDate)
	if err != nil {
		t.Fatalf("NextDueDate error = %v", err)
	}
	if got.Before(today) {
		t.Errorf("NextDueDate with past date = %v, want >= today (%v)", got, today)
	}
}

func TestNextDueDateWeekdays(t *testing.T) {
	// Start from a Monday far in the future
	// 2099-01-05 is a Monday
	monday := time.Date(2099, 1, 5, 0, 0, 0, 0, time.Local)

	got, err := NextDueDate("every wed,fri", monday)
	if err != nil {
		t.Fatalf("NextDueDate error = %v", err)
	}

	// Next wed after Monday Jan 5 should be Jan 7
	if got.Weekday() != time.Wednesday {
		t.Errorf("Expected Wednesday, got %v (%v)", got.Weekday(), got)
	}
}

func TestNextDueDateInvalidPattern(t *testing.T) {
	_, err := NextDueDate("invalid", time.Now())
	if err == nil {
		t.Error("expected error for invalid pattern")
	}
}
