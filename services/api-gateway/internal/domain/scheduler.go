package domain

import (
	"math"
	"time"
)

// ReviewInput drives SM-2 scheduling from a single review outcome.
type ReviewInput struct {
	Correct          bool
	Confidence       int
	Confused         bool
	ResponseTimeMs   int
	CurrentEase      float64
	CurrentInterval  int
	LapseCount       int
}

// SchedulerUpdate is the new scheduler state after applying SM-2 rules.
type SchedulerUpdate struct {
	NewEase     float64
	NewInterval int
	NewDueAt    time.Time
	NewLapses   int
}

// ComputeSM2 applies SM-2-style scheduling with confidence, confusion, and latency adjustments.
func ComputeSM2(input ReviewInput) SchedulerUpdate {
	quality := float64(input.Confidence)
	quality = quality * 1.25

	if input.Confused && input.Correct {
		quality = max(quality-1.5, 0)
	} else if input.Confused && !input.Correct {
		quality = 0
	}

	if input.ResponseTimeMs > 8000 {
		quality = max(quality-0.5, 0)
	}

	newEase := input.CurrentEase + (0.1 - (5-quality)*(0.08+(5-quality)*0.02))
	if newEase < 1.3 {
		newEase = 1.3
	}

	lapses := input.LapseCount
	var newInterval int
	if quality < 3 {
		newInterval = 1
		lapses++
	} else if input.CurrentInterval == 0 {
		newInterval = 1
	} else if input.CurrentInterval == 1 {
		newInterval = 6
	} else {
		newInterval = int(math.Round(float64(input.CurrentInterval) * newEase))
	}

	return SchedulerUpdate{
		NewEase:     newEase,
		NewInterval: newInterval,
		NewDueAt:    time.Now().Add(time.Duration(newInterval) * 24 * time.Hour),
		NewLapses:   lapses,
	}
}
