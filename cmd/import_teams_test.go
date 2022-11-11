package cmd_test

import (
	"bytes"
	"errors"
	"io"
	"os"
	"time"

	"github.com/vatriathlon/stopwatch2/cmd"
	"gopkg.in/yaml.v3"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/types"
)

var _ = Describe("import teams", func() {
	It("should import teams of a selected race", func() {
		// given
		source := "../tmp/teams.csv"
		output, err := os.CreateTemp(os.TempDir(), "teams-*.yaml")
		Expect(err).NotTo(HaveOccurred())
		// when
		err = cmd.ImportCSV("Bike & Run XS", source, output.Name())
		// then
		Expect(err).NotTo(HaveOccurred())
		Expect(output.Name()).To(HaveTeams(
			cmd.Team{
				Name:        "Team 1",
				Gender:      "M",
				AgeCategory: "Master",
				BibNumber:   1,
				Members: []cmd.TeamMember{
					{
						FirstName:   "Élise",
						LastName:    "Bonnin",
						DateOfBirth: parseDate("1977-04-26"),
						Category:    "Master",
						Gender:      "F",
						Club:        "",
					},
					{
						FirstName:   "Bernard",
						LastName:    "Georges",
						DateOfBirth: parseDate("1975-01-26"),
						Category:    "Master",
						Gender:      "H",
						Club:        "TOBESPORT",
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
						FirstName:   "Océane",
						LastName:    "Aubert",
						DateOfBirth: parseDate("1956-07-21"),
						Category:    "Master",
						Gender:      "F",
						Club:        "",
					},
					{
						FirstName:   "Paulette",
						LastName:    "Le Gall",
						DateOfBirth: parseDate("1963-06-02"),
						Category:    "Master",
						Gender:      "F",
						Club:        "",
					},
				},
			},
			cmd.Team{
				Name:        "Team 3",
				Gender:      "F",
				AgeCategory: "Minime",
				BibNumber:   3,
				Members: []cmd.TeamMember{
					{
						FirstName:   "Margaud",
						LastName:    "Lamy",
						DateOfBirth: parseDate("2004-08-02"),
						Category:    "Minime",
						Gender:      "F",
						Club:        "VILLENEUVE D ASCQ TRIATHLON",
					},
					{
						FirstName:   "Lorraine",
						LastName:    "Poulain",
						DateOfBirth: parseDate("2004-08-30"),
						Category:    "Minime",
						Gender:      "F",
						Club:        "VILLENEUVE D ASCQ TRIATHLON",
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
		}, Equal(expected)),
	)
}

var _ = DescribeTable("age categories",

	func(dateOfBirth string, expected string) {
		// given
		pattern := "2006-01-02"
		d, err := time.Parse(pattern, dateOfBirth)
		Expect(err).NotTo(HaveOccurred())
		// when
		result := cmd.GetAgeCategory(d)
		// then
		Expect(result).To(Equal(expected))
	},
	Entry("mini poussin", "2012-02-03", cmd.MiniPoussin),
	Entry("poussin", "2010-02-03", cmd.Poussin),
	Entry("pupille", "2009-02-03", cmd.Pupille),
	Entry("benjamin", "2007-02-03", cmd.Benjamin),
	Entry("cadet", "2002-02-03", cmd.Cadet),
	Entry("junior", "2000-02-03", cmd.Junior),
	Entry("senior", "1981-02-03", cmd.Senior),
	Entry("junior", "1975-02-03", cmd.Master),
)

var _ = DescribeTable("team age categories",
	func(category1, category2 string, expected string) {
		result := cmd.GetTeamAgeCategory(category1, category2)
		// then
		Expect(result).To(Equal(expected))
	},
	Entry("mini poussin/mini poussin", cmd.MiniPoussin, cmd.MiniPoussin, cmd.MiniPoussin),
	Entry("mini poussin/poussin", cmd.MiniPoussin, cmd.Poussin, cmd.Poussin),
	Entry("poussin/poussin", cmd.Poussin, cmd.Poussin, cmd.Poussin),
	Entry("poussin/pupille", cmd.Poussin, cmd.Pupille, cmd.Pupille),
	Entry("benjamin/minime", cmd.Benjamin, cmd.Minime, cmd.Minime),
	Entry("senior/senior", cmd.Senior, cmd.Senior, cmd.Senior),
	Entry("senior/senior", cmd.Senior, cmd.Senior, cmd.Senior),
	Entry("master/senior", cmd.Master, cmd.Senior, cmd.Senior),
	Entry("master/master", cmd.Master, cmd.Master, cmd.Master),
)
