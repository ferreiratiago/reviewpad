// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file

package utils

import (
	"fmt"
	"strings"
	"time"
)

const (
	SecondsInAYear   = 31536000
	SecondsInAMonth  = 2628288
	SecondsInADay    = 86400
	SecondsInAnHour  = 3600
	SecondsInAMinute = 60
)

func pluralize(part int, label string) string {
	if part > 1 {
		return label + "s"
	}

	return label
}

// taken from https://stackoverflow.com/a/36531443/5288071
func diff(a, b time.Time) (year, month, day, hour, min, sec int) {
	if a.Location() != b.Location() {
		b = b.In(a.Location())
	}
	if a.After(b) {
		a, b = b, a
	}
	y1, M1, d1 := a.Date()
	y2, M2, d2 := b.Date()

	h1, m1, s1 := a.Clock()
	h2, m2, s2 := b.Clock()

	year = int(y2 - y1)
	month = int(M2 - M1)
	day = int(d2 - d1)
	hour = int(h2 - h1)
	min = int(m2 - m1)
	sec = int(s2 - s1)

	// Normalize negative values
	if sec < 0 {
		sec += 60
		min--
	}
	if min < 0 {
		min += 60
		hour--
	}
	if hour < 0 {
		hour += 24
		day--
	}
	if day < 0 {
		// days in month:
		t := time.Date(y1, M1, 32, 0, 0, 0, 0, time.UTC)
		day += 32 - t.Day()
		month--
	}
	if month < 0 {
		month += 12
		year--
	}

	return
}

func ReadableTimeDiff(x time.Time, y time.Time) string {
	years, months, days, hours, minutes, seconds := diff(x, y)
	parts := make([]string, 0, 6)

	if years > 0 {
		parts = append(parts, fmt.Sprintf("%d %s", years, pluralize(years, "year")))
	}

	if months > 0 {
		parts = append(parts, fmt.Sprintf("%d %s", months, pluralize(months, "month")))
	}

	if days > 0 {
		parts = append(parts, fmt.Sprintf("%d %s", days, pluralize(days, "day")))
	}

	if hours > 0 {
		parts = append(parts, fmt.Sprintf("%d %s", hours, pluralize(hours, "hour")))
	}

	if minutes > 0 {
		parts = append(parts, fmt.Sprintf("%d %s", minutes, pluralize(minutes, "minute")))
	}

	if seconds > 0 {
		parts = append(parts, fmt.Sprintf("%d %s", seconds, pluralize(seconds, "second")))
	}

	l := len(parts)

	if l == 1 {
		return parts[0]
	} else if l == 2 {
		return parts[0] + " and " + parts[1]
	}

	return strings.Join(parts[:l-1], ", ") + " and " + parts[l-1]
}
