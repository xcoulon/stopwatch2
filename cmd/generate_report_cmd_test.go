package cmd_test

import (
	"fmt"
	"os"
	"time"

	"github.com/vatriathlon/stopwatch2/cmd"
	. "github.com/vatriathlon/stopwatch2/testsupport"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("generate reports", func() {

	var teamFile, timingFile *os.File
	BeforeEach(func() {
		var err error
		teamFile, err = os.CreateTemp("", "teams*.yaml")
		Expect(err).NotTo(HaveOccurred())
		teamFile.WriteString(teams)
		teamFile.Close()
		timingFile, err = os.CreateTemp("", "timing-xs*.yaml")
		Expect(err).NotTo(HaveOccurred())
		timingFile.WriteString(timings)
		timingFile.Close()
	})

	team1 := cmd.Team{
		Name:        "Team 1",
		Gender:      "H",
		AgeCategory: "Master",
		BibNumber:   1,
		Members: []cmd.TeamMember{
			{
				FirstName:   "Firstname1.1",
				LastName:    "Lastname1.1",
				DateOfBirth: parseDate("1977-01-01"),
				Category:    "Master",
				Gender:      "H",
				Club:        "",
			},
			{
				FirstName:   "Firstname1.2",
				LastName:    "Lastname1.2",
				DateOfBirth: parseDate("1977-01-02"),
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
				DateOfBirth: parseDate("1977-01-01"),
				Category:    "Master",
				Gender:      "F",
				Club:        "LILLE TRIATHLON",
			},
			{
				FirstName:   "Firstname2.2",
				LastName:    "Lastname2.2",
				DateOfBirth: parseDate("1977-01-02"),
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
				DateOfBirth: parseDate("1977-01-01"),
				Category:    "Master",
				Gender:      "F",
				Club:        "VILLENEUVE D'ASCQ TRIATHLON",
			},
			{
				FirstName:   "Firstname3.2",
				LastName:    "Lastname3.2",
				DateOfBirth: parseDate("1977-01-02"),
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
				DateOfBirth: parseDate("1977-01-01"),
				Category:    "Master",
				Gender:      "H",
				Club:        "",
			},
			{
				FirstName:   "Firstname4.2",
				LastName:    "Lastname4.2",
				DateOfBirth: parseDate("1977-01-02"),
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
				DateOfBirth: parseDate("1987-01-01"),
				Category:    "Senior",
				Gender:      "F",
				Club:        "",
			},
			{
				FirstName:   "Firstname5.2",
				LastName:    "Lastname5.2",
				DateOfBirth: parseDate("1987-01-02"),
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
				DateOfBirth: parseDate("1987-01-01"),
				Category:    "Senior",
				Gender:      "F",
				Club:        "",
			},
			{
				FirstName:   "Firstname6.2",
				LastName:    "Lastname6.2",
				DateOfBirth: parseDate("1987-01-02"),
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
				DateOfBirth: parseDate("1987-01-01"),
				Category:    "Senior",
				Gender:      "F",
				Club:        "",
			},
			{
				FirstName:   "Firstname17.2",
				LastName:    "Lastname17.2",
				DateOfBirth: parseDate("1987-01-02"),
				Category:    "Senior",
				Gender:      "H",
				Club:        "",
			},
		},
	}

	It("should compute overall results", func() {
		// when
		results, err := cmd.NewOverallResults(teamFile.Name(), timingFile.Name())

		// then
		Expect(err).NotTo(HaveOccurred())
		Expect(results).To(MatchTeamResults([]cmd.TeamResult{
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
		}))
	})

	It("should generate reports", func() {
		// given
		outputDir := os.TempDir()

		// when
		overallFilename, perCategoryFilename, err := cmd.GenerateReport("Bike & Run XS", teamFile.Name(), timingFile.Name(), outputDir)

		Expect(err).NotTo(HaveOccurred())

		// general results
		result, err := os.ReadFile(overallFilename)
		Expect(err).NotTo(HaveOccurred())
		fmt.Println(string(result))
		Expect(string(result)).To(Equal(`= Bike & Run XS - Classement Général

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
`))

		// results per category
		result, err = os.ReadFile(perCategoryFilename)
		Expect(err).NotTo(HaveOccurred())
		fmt.Println(string(result))
		Expect(string(result)).To(Equal(`= Bike & Run XS - Classement Par Catégorie

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
|===

== Master/M

[cols="2,5,5,10,10,5"]
|===
|# |Dossard |Equipe |Coureurs |Club |Temps Total

|1 |3 |Team 3 |Lastname3.1 - Lastname3.2 |VILLENEUVE D'ASCQ TRIATHLON |50m0s 
|===

`))
	})

})
