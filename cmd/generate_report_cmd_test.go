package cmd_test

import (
	"context"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
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
	logger := cmd.NewLogger(true)
	team1 := cmd.NewTeam("Team 1", "H", "Master", 1,
		cmd.NewTeamMember("Firstname1.1", "Lastname1.1", parseDate(t, "1977-01-01"), "Master", "H", ""),
		cmd.NewTeamMember("Firstname1.2", "Lastname1.2", parseDate(t, "1977-01-02"), "Master", "H", ""),
	)
	team2 := cmd.NewTeam("Team 2", "F", "Master", 2,
		cmd.NewTeamMember("Firstname2.1", "Lastname2.1", parseDate(t, "1977-01-01"), "Master", "F", "LILLE TRIATHLON"),
		cmd.NewTeamMember("Firstname2.2", "Lastname2.2", parseDate(t, "1977-01-02"), "Master", "F", ""),
	)
	team3 := cmd.NewTeam("Team 3", "M", "Master", 3,
		cmd.NewTeamMember("Firstname3.1", "Lastname3.1", parseDate(t, "1977-01-01"), "Master", "F", "VILLENEUVE D'ASCQ TRIATHLON"),
		cmd.NewTeamMember("Firstname3.2", "Lastname3.2", parseDate(t, "1977-01-02"), "Master", "H", ""),
	)
	team4 := cmd.NewTeam("Team 4", "H", "Master", 4,
		cmd.NewTeamMember("Firstname4.1", "Lastname4.1", parseDate(t, "1977-01-01"), "Master", "H", ""),
		cmd.NewTeamMember("Firstname4.2", "Lastname4.2", parseDate(t, "1977-01-02"), "Master", "H", ""),
	)
	team5 := cmd.NewTeam("Team 5", "F", "Senior", 5,
		cmd.NewTeamMember("Firstname5.1", "Lastname5.1", parseDate(t, "1987-01-01"), "Senior", "F", ""),
		cmd.NewTeamMember("Firstname5.2", "Lastname5.2", parseDate(t, "1987-01-02"), "Senior", "F", ""),
	)
	team6 := cmd.NewTeam("Team 6", "M", "Senior", 6,
		cmd.NewTeamMember("Firstname6.1", "Lastname6.1", parseDate(t, "1987-01-01"), "Senior", "F", ""),
		cmd.NewTeamMember("Firstname6.2", "Lastname6.2", parseDate(t, "1987-01-02"), "Senior", "H", ""),
	)
	team7 := cmd.NewTeam("Team 7", "M", "Senior", 7,
		cmd.NewTeamMember("Firstname7.1", "Lastname7.1", parseDate(t, "1987-01-01"), "Senior", "F", ""),
		cmd.NewTeamMember("Firstname7.2", "Lastname7.2", parseDate(t, "1987-01-02"), "Senior", "H", ""),
	)
	team8 := cmd.NewTeam("Team 8", "M", "Senior", 8,
		cmd.NewTeamMember("Firstname8.1", "Lastname8.1", parseDate(t, "1987-01-01"), "Senior", "F", ""),
		cmd.NewTeamMember("Firstname8.2", "Lastname8.2", parseDate(t, "1987-01-02"), "Senior", "H", ""),
	)
	team9 := cmd.NewTeam("Team 9", "M", "Senior", 9,
		cmd.NewTeamMember("Firstname9.1", "Lastname9.1", parseDate(t, "1987-01-01"), "Senior", "F", ""),
		cmd.NewTeamMember("Firstname9.2", "Lastname9.2", parseDate(t, "1987-01-02"), "Senior", "H", ""),
	)
	team10 := cmd.NewTeam("Team 10", "M", "Senior", 10,
		cmd.NewTeamMember("Firstname10.1", "Lastname10.1", parseDate(t, "1987-01-01"), "Senior", "F", ""),
		cmd.NewTeamMember("Firstname10.2", "Lastname10.2", parseDate(t, "1987-01-02"), "Senior", "H", ""),
	)
	team11 := cmd.NewTeam("Team 11", "M", "Senior", 11,
		cmd.NewTeamMember("Firstname11.1", "Lastname11.1", parseDate(t, "1987-01-01"), "Senior", "F", ""),
		cmd.NewTeamMember("Firstname11.2", "Lastname11.2", parseDate(t, "1987-01-02"), "Senior", "H", ""),
	)
	team12 := cmd.NewTeam("Team 12", "M", "Senior", 12,
		cmd.NewTeamMember("Firstname12.1", "Lastname12.1", parseDate(t, "1987-01-01"), "Senior", "F", ""),
		cmd.NewTeamMember("Firstname12.2", "Lastname12.2", parseDate(t, "1987-01-02"), "Senior", "H", ""),
	)
	team17 := cmd.NewTeam("Team 17", "M", "Senior", 17,
		cmd.NewTeamMember("Firstname17.1", "Lastname17.1", parseDate(t, "1987-01-01"), "Senior", "F", ""),
		cmd.NewTeamMember("Firstname17.2", "Lastname17.2", parseDate(t, "1987-01-02"), "Senior", "H", ""),
	)
	t.Run("overall results", func(t *testing.T) {
		// when
		actual, err := cmd.NewOverallResults(teamFilename, timingFilename)

		// then
		require.NoError(t, err)
		expected := []cmd.TeamResult{
			cmd.NewTeamResult(team3, 1, 50*time.Minute),
			cmd.NewTeamResult(team4, 2, 51*time.Minute),
			cmd.NewTeamResult(team2, 3, 51*time.Minute+30*time.Second),
			cmd.NewTeamResult(team1, 4, 52*time.Minute+55*time.Second),
			cmd.NewTeamResult(team5, 5, 53*time.Minute+30*time.Second),
			cmd.NewTeamResult(team17, 6, 54*time.Minute+30*time.Second),
			cmd.NewTeamResult(team6, 7, 56*time.Minute+30*time.Second),
			cmd.NewTeamResult(team7, 8, 57*time.Minute+30*time.Second),
			cmd.NewTeamResult(team8, 9, 58*time.Minute+30*time.Second),
			cmd.NewTeamResult(team9, 10, 59*time.Minute+30*time.Second),
			cmd.NewTeamResult(team10, 11, 60*time.Minute+30*time.Second),
			cmd.NewTeamResult(team11, 12, 61*time.Minute+30*time.Second),
			cmd.NewTeamResult(team12, 13, 62*time.Minute+30*time.Second),
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

		expected := `== Bike & Run XS - Classement Général

[cols="1,5,5,10,10,5"]
|===
|# |Equipe |Catégorie |Coureur 1 |Coureur 2 |Temps

|1 |3 Team 3 |Master/M |Firstname3.1 Lastname3.1 VILLENEUVE D'ASCQ TRIATHLON |Firstname3.2 Lastname3.2  |50m0s 
|2 |4 Team 4 |Master/H |Firstname4.1 Lastname4.1  |Firstname4.2 Lastname4.2  |51m0s 
|3 |2 Team 2 |Master/F |Firstname2.1 Lastname2.1 LILLE TRIATHLON |Firstname2.2 Lastname2.2  |51m30s 
|4 |1 Team 1 |Master/H |Firstname1.1 Lastname1.1  |Firstname1.2 Lastname1.2  |52m55s 
|5 |5 Team 5 |Senior/F |Firstname5.1 Lastname5.1  |Firstname5.2 Lastname5.2  |53m30s 
|6 |17 Team 17 |Senior/M |Firstname17.1 Lastname17.1  |Firstname17.2 Lastname17.2  |54m30s 
|7 |6 Team 6 |Senior/M |Firstname6.1 Lastname6.1  |Firstname6.2 Lastname6.2  |56m30s 
|8 |7 Team 7 |Senior/M |Firstname7.1 Lastname7.1  |Firstname7.2 Lastname7.2  |57m30s 
|9 |8 Team 8 |Senior/M |Firstname8.1 Lastname8.1  |Firstname8.2 Lastname8.2  |58m30s 
|10 |9 Team 9 |Senior/M |Firstname9.1 Lastname9.1  |Firstname9.2 Lastname9.2  |59m30s 
|11 |10 Team 10 |Senior/M |Firstname10.1 Lastname10.1  |Firstname10.2 Lastname10.2  |1h0m30s 
|12 |11 Team 11 |Senior/M |Firstname11.1 Lastname11.1  |Firstname11.2 Lastname11.2  |1h1m30s 
|13 |12 Team 12 |Senior/M |Firstname12.1 Lastname12.1  |Firstname12.2 Lastname12.2  |1h2m30s 
|===
`
		assert.Equal(t, expected, string(actual))
	})

	t.Run("per category", func(t *testing.T) {
		// results per category
		actual, err := os.ReadFile(perCategoryFilename)
		require.NoError(t, err)
		logger.Debug("results per category", "contents", string(actual))
		expected := `== Bike & Run XS - Classement Par Catégorie

=== Senior/F

[cols="1,10,10,10,5"]
|===
|# |Equipe |Coureur 1 |Coureur 2 |Temps

|5 |5 Team 5 |Firstname5.1 Lastname5.1  |Firstname5.2 Lastname5.2  |53m30s
|===

=== Senior/M

[cols="1,10,10,10,5"]
|===
|# |Equipe |Coureur 1 |Coureur 2 |Temps

|6 |17 Team 17 |Firstname17.1 Lastname17.1  |Firstname17.2 Lastname17.2  |54m30s
|7 |6 Team 6 |Firstname6.1 Lastname6.1  |Firstname6.2 Lastname6.2  |56m30s
|8 |7 Team 7 |Firstname7.1 Lastname7.1  |Firstname7.2 Lastname7.2  |57m30s
|===

=== Master/F

[cols="1,10,10,10,5"]
|===
|# |Equipe |Coureur 1 |Coureur 2 |Temps

|3 |2 Team 2 |Firstname2.1 Lastname2.1 LILLE TRIATHLON |Firstname2.2 Lastname2.2  |51m30s
|===

=== Master/H

[cols="1,10,10,10,5"]
|===
|# |Equipe |Coureur 1 |Coureur 2 |Temps

|2 |4 Team 4 |Firstname4.1 Lastname4.1  |Firstname4.2 Lastname4.2  |51m0s
|4 |1 Team 1 |Firstname1.1 Lastname1.1  |Firstname1.2 Lastname1.2  |52m55s
|===

=== Master/M

[cols="1,10,10,10,5"]
|===
|# |Equipe |Coureur 1 |Coureur 2 |Temps

|1 |3 Team 3 |Firstname3.1 Lastname3.1 VILLENEUVE D'ASCQ TRIATHLON |Firstname3.2 Lastname3.2  |50m0s
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
	err = os.WriteFile(teamFilename, []byte(teams), 0600)
	require.NoError(t, err)
	err = os.WriteFile(timingFilename, []byte(timings), 0600)
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
			diff = strings.ReplaceAll(diff, "\u00a0", "")
			diff = strings.ReplaceAll(diff, "\t", "  ")
			logger.Info("team results are not equal", "diff", diff)
			return false
		}
		return true
	}
}

const teams = `---
name: Team 1
gender: H
ageCategory: Master
bibNumber: 1
members:
    - firstName: Firstname1.1
      lastName: Lastname1.1
      dateOfBirth: "1977-01-01"
      category: Master
      gender: H
      club: ""
    - firstName: Firstname1.2
      lastName: Lastname1.2
      dateOfBirth: "1977-01-02"
      category: Master
      gender: H
      club: ""
---
name: Team 2
gender: F
ageCategory: Master
bibNumber: 2
members:
    - firstName: Firstname2.1
      lastName: Lastname2.1
      dateOfBirth: "1977-01-01"
      category: Master
      gender: F
      club: LILLE TRIATHLON
    - firstName: Firstname2.2
      lastName: Lastname2.2
      dateOfBirth: "1977-01-02"
      category: Master
      gender: F
---
name: Team 3
gender: M
ageCategory: Master
bibNumber: 3
members:
    - firstName: Firstname3.1
      lastName: Lastname3.1
      dateOfBirth: "1977-01-01"
      category: Master
      gender: F
      club: VILLENEUVE D'ASCQ TRIATHLON
    - firstName: Firstname3.2
      lastName: Lastname3.2
      dateOfBirth: "1977-01-02"
      category: Master
      gender: H
---
name: Team 4
gender: H
ageCategory: Master
bibNumber: 4
members:
    - firstName: Firstname4.1
      lastName: Lastname4.1
      dateOfBirth: "1977-01-01"
      category: Master
      gender: H
      club: ""
    - firstName: Firstname4.2
      lastName: Lastname4.2
      dateOfBirth: "1977-01-02"
      category: Master
      gender: H
      club: ""
---
name: Team 5
gender: F
ageCategory: Senior
bibNumber: 5
members:
    - firstName: Firstname5.1
      lastName: Lastname5.1
      dateOfBirth: "1987-01-01"
      category: Senior
      gender: F
      club: ""
    - firstName: Firstname5.2
      lastName: Lastname5.2
      dateOfBirth: "1987-01-02"
      category: Senior
      gender: F
      club: ""
---
name: Team 6
gender: M
ageCategory: Senior
bibNumber: 6
members:
    - firstName: Firstname6.1
      lastName: Lastname6.1
      dateOfBirth: "1987-01-01"
      category: Senior
      gender: F
      club: ""
    - firstName: Firstname6.2
      lastName: Lastname6.2
      dateOfBirth: "1987-01-02"
      category: Senior
      gender: H
      club: ""
---
name: Team 7
gender: M
ageCategory: Senior
bibNumber: 7
members:
    - firstName: Firstname7.1
      lastName: Lastname7.1
      dateOfBirth: "1987-01-01"
      category: Senior
      gender: F
      club: ""
    - firstName: Firstname7.2
      lastName: Lastname7.2
      dateOfBirth: "1987-01-02"
      category: Senior
      gender: H
      club: ""
---
name: Team 8
gender: M
ageCategory: Senior
bibNumber: 8
members:
    - firstName: Firstname8.1
      lastName: Lastname8.1
      dateOfBirth: "1987-01-01"
      category: Senior
      gender: F
      club: ""
    - firstName: Firstname8.2
      lastName: Lastname8.2
      dateOfBirth: "1987-01-02"
      category: Senior
      gender: H
      club: ""
---
name: Team 9
gender: M
ageCategory: Senior
bibNumber: 9
members:
    - firstName: Firstname9.1
      lastName: Lastname9.1
      dateOfBirth: "1987-01-01"
      category: Senior
      gender: F
      club: ""
    - firstName: Firstname9.2
      lastName: Lastname9.2
      dateOfBirth: "1987-01-02"
      category: Senior
      gender: H
      club: ""
---
name: Team 10
gender: M
ageCategory: Senior
bibNumber: 10
members:
    - firstName: Firstname10.1
      lastName: Lastname10.1
      dateOfBirth: "1987-01-01"
      category: Senior
      gender: F
      club: ""
    - firstName: Firstname10.2
      lastName: Lastname10.2
      dateOfBirth: "1987-01-02"
      category: Senior
      gender: H
      club: ""
---
name: Team 11
gender: M
ageCategory: Senior
bibNumber: 11
members:
    - firstName: Firstname11.1
      lastName: Lastname11.1
      dateOfBirth: "1987-01-01"
      category: Senior
      gender: F
      club: ""
    - firstName: Firstname11.2
      lastName: Lastname11.2
      dateOfBirth: "1987-01-02"
      category: Senior
      gender: H
      club: ""
---
name: Team 12
gender: M
ageCategory: Senior
bibNumber: 12
members:
    - firstName: Firstname12.1
      lastName: Lastname12.1
      dateOfBirth: "1987-01-01"
      category: Senior
      gender: F
      club: ""
    - firstName: Firstname12.2
      lastName: Lastname12.2
      dateOfBirth: "1987-01-02"
      category: Senior
      gender: H
      club: ""
---
name: Team 17
gender: M
ageCategory: Senior
bibNumber: 17
members:
    - firstName: Firstname17.1
      lastName: Lastname17.1
      dateOfBirth: "1987-01-01"
      category: Senior
      gender: F
      club: ""
    - firstName: Firstname17.2
      lastName: Lastname17.2
      dateOfBirth: "1987-01-02"
      category: Senior
      gender: H
      club: ""
---
name: Team 101
gender: H
ageCategory: Minime
bibNumber: 101
members:
    - firstName: Firstname101.1
      lastName: Lastname101.1
      dateOfBirth: 2007-01-01
      category: Minime
      gender: H
      club: VILLENEUVE D'ASCQ TRIATHLON
    - firstName: Firstname101.2
      lastName: Lastname101.2
      dateOfBirth: 2007-01-02
      category: Minime
      gender: H
      club: VILLENEUVE D'ASCQ TRIATHLON
---
name:  Team 201
gender: H
ageCategory: Poussin
bibNumber: 201
members:
    - firstName: Firstname201.1
      lastName: Lastname201.1
      dateOfBirth: 2014-01-01
      category: Minime
      gender: H
      club: VILLENEUVE D'ASCQ TRIATHLON
    - firstName: Firstname201.2
      lastName: Lastname201.2
      dateOfBirth: 2014-01-02
      category: Minime
      gender: H
      club: VILLENEUVE D'ASCQ TRIATHLON
`

const timings = `- 10:00:00: start
- 10:50:00: 3
- 10:51:00: 4
- 10:51:30: 2
- 10:52:55: 1
- 10:53:30: 5
- 10:54:30: 17
- 10:56:30: 6
- 10:57:30: 7
- 10:58:30: 8
- 10:59:30: 9
- 11:00:30: 10
- 11:01:30: 11
- 11:02:30: 12
`
