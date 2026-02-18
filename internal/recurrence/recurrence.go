package recurrence

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// Weekday abbreviation mappings
var weekdayNames = map[string]time.Weekday{
	"sunday":    time.Sunday,
	"sun":       time.Sunday,
	"monday":    time.Monday,
	"mon":       time.Monday,
	"tuesday":   time.Tuesday,
	"tue":       time.Tuesday,
	"wednesday": time.Wednesday,
	"wed":       time.Wednesday,
	"thursday":  time.Thursday,
	"thu":       time.Thursday,
	"friday":    time.Friday,
	"fri":       time.Friday,
	"saturday":  time.Saturday,
	"sat":       time.Saturday,
}

// ParsePattern validates and normalizes a recurrence pattern string.
// Returns the normalized pattern or an error if invalid.
//
// Supported patterns:
//   - daily, weekly, monthly, yearly
//   - every <N>d, every <N>w, every <N>m, every <N>y
//   - every monday, every mon,wed,fri
func ParsePattern(pattern string) (string, error) {
	pattern = strings.TrimSpace(strings.ToLower(pattern))
	if pattern == "" {
		return "", fmt.Errorf("empty recurrence pattern")
	}

	switch pattern {
	case "daily", "weekly", "monthly", "yearly":
		return pattern, nil
	}

	if !strings.HasPrefix(pattern, "every ") {
		return "", fmt.Errorf("invalid recurrence pattern: %q (expected daily, weekly, monthly, yearly, or every ...)", pattern)
	}

	spec := strings.TrimSpace(pattern[6:])
	if spec == "" {
		return "", fmt.Errorf("invalid recurrence pattern: %q (missing interval after 'every')", pattern)
	}

	// Try interval+unit pattern: every <N>d/w/m/y
	if len(spec) >= 2 {
		unit := spec[len(spec)-1]
		numStr := spec[:len(spec)-1]
		if unit == 'd' || unit == 'w' || unit == 'm' || unit == 'y' {
			n, err := strconv.Atoi(numStr)
			if err == nil {
				if n <= 0 {
					return "", fmt.Errorf("invalid recurrence interval: %d (must be positive)", n)
				}
				return fmt.Sprintf("every %d%c", n, unit), nil
			}
		}
	}

	// Try day-of-week pattern: every mon,wed,fri
	parts := strings.Split(spec, ",")
	var days []string
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if _, ok := weekdayNames[p]; !ok {
			return "", fmt.Errorf("invalid recurrence pattern: %q (unknown day or interval: %q)", pattern, p)
		}
		days = append(days, p)
	}

	return "every " + strings.Join(days, ","), nil
}

// NextDueDate computes the next due date based on a recurrence pattern and the current due date.
// It always advances past today so late completions still get a future date.
func NextDueDate(pattern string, currentDue time.Time) (time.Time, error) {
	pattern = strings.TrimSpace(strings.ToLower(pattern))
	today := time.Now()
	today = time.Date(today.Year(), today.Month(), today.Day(), 0, 0, 0, 0, today.Location())

	var next time.Time
	var err error

	switch pattern {
	case "daily":
		next = advanceByInterval(currentDue, 1, 'd', today)
	case "weekly":
		next = advanceByInterval(currentDue, 1, 'w', today)
	case "monthly":
		next = advanceByInterval(currentDue, 1, 'm', today)
	case "yearly":
		next = advanceByInterval(currentDue, 1, 'y', today)
	default:
		if strings.HasPrefix(pattern, "every ") {
			spec := strings.TrimSpace(pattern[6:])
			next, err = parseEverySpec(spec, currentDue, today)
			if err != nil {
				return time.Time{}, err
			}
		} else {
			return time.Time{}, fmt.Errorf("invalid recurrence pattern: %q", pattern)
		}
	}

	return next, nil
}

func parseEverySpec(spec string, currentDue, today time.Time) (time.Time, error) {
	// Try interval+unit
	if len(spec) >= 2 {
		unit := spec[len(spec)-1]
		numStr := spec[:len(spec)-1]
		if unit == 'd' || unit == 'w' || unit == 'm' || unit == 'y' {
			n, err := strconv.Atoi(numStr)
			if err == nil && n > 0 {
				return advanceByInterval(currentDue, n, unit, today), nil
			}
		}
	}

	// Day-of-week pattern
	parts := strings.Split(spec, ",")
	var weekdays []time.Weekday
	for _, p := range parts {
		p = strings.TrimSpace(p)
		wd, ok := weekdayNames[p]
		if !ok {
			return time.Time{}, fmt.Errorf("unknown day: %q", p)
		}
		weekdays = append(weekdays, wd)
	}

	if len(weekdays) == 0 {
		return time.Time{}, fmt.Errorf("no weekdays specified")
	}

	return nextMatchingWeekday(currentDue, weekdays, today), nil
}

// advanceByInterval advances from currentDue by the given interval,
// repeating until the result is strictly after today.
func advanceByInterval(currentDue time.Time, n int, unit byte, today time.Time) time.Time {
	next := currentDue
	for {
		switch unit {
		case 'd':
			next = next.AddDate(0, 0, n)
		case 'w':
			next = next.AddDate(0, 0, n*7)
		case 'm':
			next = next.AddDate(0, n, 0)
		case 'y':
			next = next.AddDate(n, 0, 0)
		}
		if !next.Before(today) {
			break
		}
	}
	return next
}

// nextMatchingWeekday finds the next date after currentDue that falls on one of the given weekdays,
// ensuring it's not before today.
func nextMatchingWeekday(currentDue time.Time, weekdays []time.Weekday, today time.Time) time.Time {
	// Build a set for fast lookup
	daySet := make(map[time.Weekday]bool, len(weekdays))
	for _, wd := range weekdays {
		daySet[wd] = true
	}

	// Start from the day after currentDue
	candidate := currentDue.AddDate(0, 0, 1)
	// But if candidate is before today, start from today
	if candidate.Before(today) {
		candidate = today
	}

	// Search up to 7 days (guaranteed to find a match)
	for i := 0; i < 7; i++ {
		if daySet[candidate.Weekday()] {
			return candidate
		}
		candidate = candidate.AddDate(0, 0, 1)
	}

	// Should never reach here if weekdays is non-empty
	return candidate
}
