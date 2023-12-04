package cmd_test

import (
	"context"
	"log/slog"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/alecthomas/assert"
	"github.com/google/go-cmp/cmp"
	"github.com/sanity-io/litter"
	"github.com/stretchr/testify/require"
	"github.com/vatriathlon/stopwatch2/cmd"
)

func TestNewOverallResults(t *testing.T) {

	//given
	teamFilename, timingFilename := setupRawResults(t)
	logger := cmd.NewLogger(false)
	team1 := cmd.Team{
		Name:        "Team 1",
		Gender:      "H",
		AgeCategory: "Master",
		BibNumber:   1,
		Members: []cmd.TeamMember{
			{
				FirstName:   "Firstname1.1",
				LastName:    "Lastname1.1",
				DateOfBirth: parseDate(t, "1977-01-01"),
				Category:    "Master",
				Gender:      "H",
				Club:        "",
			},
			{
				FirstName:   "Firstname1.2",
				LastName:    "Lastname1.2",
				DateOfBirth: parseDate(t, "1977-01-02"),
				Category:    "Master",
				Gender:      "H",
				Club:        "",
			},
		},
	}

	team2 := cmd.Team{
		Name:        "Team 2",
		Gender:      "F",
		AgeCategory: "Master",
		BibNumber:   2,
		Members: []cmd.TeamMember{
			{
				FirstName:   "Firstname2.1",
				LastName:    "Lastname2.1",
				DateOfBirth: parseDate(t, "1977-01-01"),
				Category:    "Master",
				Gender:      "F",
				Club:        "LILLE TRIATHLON",
			},
			{
				FirstName:   "Firstname2.2",
				LastName:    "Lastname2.2",
				DateOfBirth: parseDate(t, "1977-01-02"),
				Category:    "Master",
				Gender:      "F",
				Club:        "",
			},
		},
	}

	team3 := cmd.Team{
		Name:        "Team 3",
		Gender:      "M",
		AgeCategory: "Master",
		BibNumber:   3,
		Members: []cmd.TeamMember{
			{
				FirstName:   "Firstname3.1",
				LastName:    "Lastname3.1",
				DateOfBirth: parseDate(t, "1977-01-01"),
				Category:    "Master",
				Gender:      "F",
				Club:        "VILLENEUVE D'ASCQ TRIATHLON",
			},
			{
				FirstName:   "Firstname3.2",
				LastName:    "Lastname3.2",
				DateOfBirth: parseDate(t, "1977-01-02"),
				Category:    "Master",
				Gender:      "H",
				Club:        "",
			},
		},
	}

	team4 := cmd.Team{
		Name:        "Team 4",
		Gender:      "H",
		AgeCategory: "Master",
		BibNumber:   4,
		Members: []cmd.TeamMember{
			{
				FirstName:   "Firstname4.1",
				LastName:    "Lastname4.1",
				DateOfBirth: parseDate(t, "1977-01-01"),
				Category:    "Master",
				Gender:      "H",
				Club:        "",
			},
			{
				FirstName:   "Firstname4.2",
				LastName:    "Lastname4.2",
				DateOfBirth: parseDate(t, "1977-01-02"),
				Category:    "Master",
				Gender:      "H",
				Club:        "",
			},
		},
	}
	team5 := cmd.Team{
		Name:        "Team 5",
		Gender:      "F",
		AgeCategory: "Senior",
		BibNumber:   5,
		Members: []cmd.TeamMember{
			{
				FirstName:   "Firstname5.1",
				LastName:    "Lastname5.1",
				DateOfBirth: parseDate(t, "1987-01-01"),
				Category:    "Senior",
				Gender:      "F",
				Club:        "",
			},
			{
				FirstName:   "Firstname5.2",
				LastName:    "Lastname5.2",
				DateOfBirth: parseDate(t, "1987-01-02"),
				Category:    "Senior",
				Gender:      "F",
				Club:        "",
			},
		},
	}

	team6 := cmd.Team{
		Name:        "Team 6",
		Gender:      "M",
		AgeCategory: "Senior",
		BibNumber:   6,
		Members: []cmd.TeamMember{
			{
				FirstName:   "Firstname6.1",
				LastName:    "Lastname6.1",
				DateOfBirth: parseDate(t, "1987-01-01"),
				Category:    "Senior",
				Gender:      "F",
				Club:        "",
			},
			{
				FirstName:   "Firstname6.2",
				LastName:    "Lastname6.2",
				DateOfBirth: parseDate(t, "1987-01-02"),
				Category:    "Senior",
				Gender:      "H",
				Club:        "",
			},
		},
	}

	team17 := cmd.Team{
		Name:        "Team 17",
		Gender:      "M",
		AgeCategory: "Senior",
		BibNumber:   17,
		Members: []cmd.TeamMember{
			{
				FirstName:   "Firstname17.1",
				LastName:    "Lastname17.1",
				DateOfBirth: parseDate(t, "1987-01-01"),
				Category:    "Senior",
				Gender:      "F",
				Club:        "",
			},
			{
				FirstName:   "Firstname17.2",
				LastName:    "Lastname17.2",
				DateOfBirth: parseDate(t, "1987-01-02"),
				Category:    "Senior",
				Gender:      "H",
				Club:        "",
			},
		},
	}

	t.Run("overall results", func(t *testing.T) {
		// when
		actual, err := cmd.NewOverallResults(teamFilename, timingFilename)

		// then
		require.NoError(t, err)
		expected := []cmd.TeamResult{
			{
				Rank:      1,
				Team:      team3,
				TotalTime: 50 * time.Minute,
			},
			{
				Rank:      2,
				Team:      team4,
				TotalTime: 51 * time.Minute,
			},
			{
				Rank:      3,
				Team:      team2,
				TotalTime: 51*time.Minute + 30*time.Second,
			},
			{
				Rank:      4,
				Team:      team1,
				TotalTime: 52*time.Minute + 55*time.Second,
			},
			{
				Rank:      5,
				Team:      team5,
				TotalTime: 53*time.Minute + 30*time.Second,
			},
			{
				Rank:      6,
				Team:      team17,
				TotalTime: 54*time.Minute + 30*time.Second,
			},
			{
				Rank:      7,
				Team:      team6,
				TotalTime: 55*time.Minute + 30*time.Second,
			},
		}
		assert.Condition(t, matchTeamResults(logger, expected, actual))
	})

}

func TestGenerateReports(t *testing.T) {

	// given
	teamFilename, timingFilename := setupRawResults(t)
	outputDir, err := os.MkdirTemp(os.TempDir(), "bikerun2023-")
	require.NoError(t, err)

	logger := cmd.NewLogger(true)

	// when
	overallFilename, perCategoryFilename, err := cmd.GenerateReport(logger, "Bike & Run XS", teamFilename, timingFilename, outputDir)

	//then
	require.NoError(t, err)

	t.Run("overall", func(t *testing.T) {
		// general results
		actual, err := os.ReadFile(overallFilename)
		require.NoError(t, err)
		logger.Debug("overall results", "contents", string(actual))

		expected := `= Bike & Run XS - Classement Général

[cols="2,5,5,5,10,10,5"]
|===
|# |Dossard |Equipe |Catégorie |Coureurs |Club |Temps Total

|1 |3 |Team 3 |Master/M |Lastname3.1 - Lastname3.2 |VILLENEUVE D'ASCQ TRIATHLON |50m0s 
|2 |4 |Team 4 |Master/H |Lastname4.1 - Lastname4.2 | |51m0s 
|3 |2 |Team 2 |Master/F |Lastname2.1 - Lastname2.2 |LILLE TRIATHLON |51m30s 
|4 |1 |Team 1 |Master/H |Lastname1.1 - Lastname1.2 | |52m55s 
|5 |5 |Team 5 |Senior/F |Lastname5.1 - Lastname5.2 | |53m30s 
|6 |17 |Team 17 |Senior/M |Lastname17.1 - Lastname17.2 | |54m30s 
|7 |6 |Team 6 |Senior/M |Lastname6.1 - Lastname6.2 | |55m30s 
|===
`
		assert.Equal(t, expected, string(actual))
	})

	t.Run("per category", func(t *testing.T) {
		// results per category
		actual, err := os.ReadFile(perCategoryFilename)
		require.NoError(t, err)
		logger.Debug("results per category", "contents", string(actual))
		expected := `= Bike & Run XS - Classement Par Catégorie

== Senior/F

[cols="2,5,5,10,10,5"]
|===
|# |Dossard |Equipe |Coureurs |Club |Temps Total

|5 |5 |Team 5 |Lastname5.1 - Lastname5.2 | |53m30s 
|===

== Senior/M

[cols="2,5,5,10,10,5"]
|===
|# |Dossard |Equipe |Coureurs |Club |Temps Total

|6 |17 |Team 17 |Lastname17.1 - Lastname17.2 | |54m30s 
|7 |6 |Team 6 |Lastname6.1 - Lastname6.2 | |55m30s 
|===

== Master/F

[cols="2,5,5,10,10,5"]
|===
|# |Dossard |Equipe |Coureurs |Club |Temps Total

|3 |2 |Team 2 |Lastname2.1 - Lastname2.2 |LILLE TRIATHLON |51m30s 
|===

== Master/H

[cols="2,5,5,10,10,5"]
|===
|# |Dossard |Equipe |Coureurs |Club |Temps Total

|2 |4 |Team 4 |Lastname4.1 - Lastname4.2 | |51m0s 
|4 |1 |Team 1 |Lastname1.1 - Lastname1.2 | |52m55s 
|===

== Master/M

[cols="2,5,5,10,10,5"]
|===
|# |Dossard |Equipe |Coureurs |Club |Temps Total

|1 |3 |Team 3 |Lastname3.1 - Lastname3.2 |VILLENEUVE D'ASCQ TRIATHLON |50m0s 
|===

`
		assert.Equal(t, expected, string(actual))
	})
}

func setupRawResults(t *testing.T) (string, string) {
	outputDir, err := os.MkdirTemp(os.TempDir(), "bikerun2023-")
	require.NoError(t, err)
	teamFilename := filepath.Join(outputDir, "teams.yaml")
	timingFilename := filepath.Join(outputDir, "timing-xs.yaml")
	err = os.WriteFile(teamFilename, []byte(teams), 0755)
	require.NoError(t, err)
	err = os.WriteFile(timingFilename, []byte(timings), 0755)
	require.NoError(t, err)
	return teamFilename, timingFilename
}

func matchTeamResults(logger *slog.Logger, expected, actual []cmd.TeamResult) assert.Comparison {
	return func() bool {
		if diff := cmp.Diff(expected, actual); diff != "" {
			if logger.Enabled(context.Background(), slog.LevelDebug) {
				logger.Debug("actual team results", "contents", litter.Sdump(actual))
				logger.Debug("expected team results", "contents", litter.Sdump(expected))
			}
			return false
		}
		return true
	}
}
