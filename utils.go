package main

import (
	"sort"
	"time"
)

func today() string {
	return time.Now().Format("2006-01-02")
}

func contains(slice []string, val string) bool {
	for _, item := range slice {
		if item == val {
			return true
		}
	}
	return false
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func calculateStreakFrom(startDate string, dates []string) int {
	if len(dates) == 0 {
		return 0
	}
	dateMap := make(map[string]bool)
	for _, d := range dates {
		dateMap[d] = true
	}
	t, err := time.Parse("2006-01-02", startDate)
	if err != nil {
		return 0
	}
	streak := 0
	for {
		ds := t.Format("2006-01-02")
		if dateMap[ds] {
			streak++
			t = t.AddDate(0, 0, -1)
		} else {
			break
		}
	}
	return streak
}

func calculateLongestStreak(dates []string) int {
	if len(dates) == 0 {
		return 0
	}
	var parsed []time.Time
	for _, d := range dates {
		t, err := time.Parse("2006-01-02", d)
		if err == nil {
			parsed = append(parsed, t)
		}
	}
	sort.Slice(parsed, func(i, j int) bool {
		return parsed[i].Before(parsed[j])
	})
	longest := 1
	current := 1
	for i := 1; i < len(parsed); i++ {
		diff := parsed[i].Sub(parsed[i-1]).Hours() / 24
		if diff == 1 {
			current++
			if current > longest {
				longest = current
			}
		} else if diff > 1 {
			current = 1
		}
	}
	return longest
}

