package cmd_test

import (
	"bytes"
	"errors"
	"io"
	"os"
	"time"

	"github.com/vatriathlon/stopwatch2/cmd"
	. "github.com/vatriathlon/stopwatch2/testsupport"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/types"
	"gopkg.in/yaml.v3"
)

var _ = Describe("import teams", func() {

	It("should import teams of a selected race", func() {
		// given
		source, err := os.Open("teams.csv")
		Expect(err).NotTo(HaveOccurred())
		output, err := os.CreateTemp(os.TempDir(), "teams-*.yaml")
		Expect(err).NotTo(HaveOccurred())
		// when
		err = cmd.ImportCSV(source.Name(), output.Name())
		// then
		Expect(err).NotTo(HaveOccurred())
		Expect(output.Name()).To(HaveTeams(
			cmd.Team{
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
			},
			cmd.Team{
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
			},
			cmd.Team{
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
			},
			cmd.Team{
				Name:        "Team 101",
				Gender:      "H",
				AgeCategory: "Minime",
				BibNumber:   101,
				Members: []cmd.TeamMember{
					{
						FirstName:   "Firstname101.1",
						LastName:    "Lastname101.1",
						DateOfBirth: parseDate("2007-01-01"),
						Category:    "Minime",
						Gender:      "H",
						Club:        "VILLENEUVE D'ASCQ TRIATHLON",
					},
					{
						FirstName:   "Firstname101.2",
						LastName:    "Lastname101.2",
						DateOfBirth: parseDate("2007-01-02"),
						Category:    "Minime",
						Gender:      "H",
						Club:        "VILLENEUVE D'ASCQ TRIATHLON",
					},
				},
			},
			cmd.Team{
				Name:        "Team 201",
				Gender:      "H",
				AgeCategory: "Poussin",
				BibNumber:   201,
				Members: []cmd.TeamMember{
					{
						FirstName:   "Firstname201.1",
						LastName:    "Lastname201.1",
						DateOfBirth: parseDate("2014-01-01"),
						Category:    "Poussin",
						Gender:      "H",
						Club:        "LILLE TRIATHLON",
					},
					{
						FirstName:   "Firstname201.2",
						LastName:    "Lastname201.2",
						DateOfBirth: parseDate("2014-01-02"),
						Category:    "Poussin",
						Gender:      "H",
						Club:        "LILLE TRIATHLON",
					},
				},
			},
		))
	})
})

func parseDate(d string) cmd.ISO8601Date {
	r, err := time.Parse("2006-01-02", d)
	Expect(err).NotTo(HaveOccurred())
	return cmd.ISO8601Date(r)
}

func HaveTeams(expected ...cmd.Team) types.GomegaMatcher {
	return And(
		WithTransform(func(filename string) ([]cmd.Team, error) {
			content, err := os.ReadFile(filename)
			if err != nil {
				return nil, err
			}
			teams := []cmd.Team{}
			decoder := yaml.NewDecoder(bytes.NewReader(content))
			// decode 1 team at a time
			for {
				team := cmd.Team{}
				if err := decoder.Decode(&team); errors.Is(err, io.EOF) {
					break
				} else if err != nil {
					return nil, err
				}
				teams = append(teams, team)
			}
			return teams, nil
		}, MatchTeams(expected)),
	)
}
