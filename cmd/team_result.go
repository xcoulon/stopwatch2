package cmd

import "time"

type TeamResult struct {
	Team
	Rank      int
	TotalTime time.Duration
}

func NewTeamResult(team Team, rank int, totalTime time.Duration) TeamResult {
	return TeamResult{
		Team:      team,
		Rank:      rank,
		TotalTime: totalTime,
	}
}
