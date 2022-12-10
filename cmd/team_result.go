package cmd

import "time"

type TeamResult struct {
	Team
	Rank      int
	TotalTime time.Duration
}
