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

const csv = `Numéro de dossard,COMPÉTITION,Nom de l'équipe,Identifiant d'inscription,Nom,Prénom,Date de naissance,Sexe,Êtes vous licencié(e) ?,Numéro de licence,Numéro de licence,Pass C,Statut du dossier,Paiement
1,Bike & Run XS,Team 1,unused,Lastname1.1,Firstname1.1,01/01/1977,h,not_member,A1,,1,Complet,Payé
1,Bike & Run XS,Team 1,unused,Lastname1.2,Firstname1.2,02/01/1977,h,not_member,,,1,Complet,Payé
2,Bike & Run XS,Team 2,unused,Lastname2.1,Firstname2.1,01/01/1977,f,not_member,A2,LILLE TRIATHLON,1,Complet,Payé
2,Bike & Run XS,Team 2,unused,Lastname2.2,Firstname2.2,02/01/1977,f,not_member,,,1,Complet,Payé
3,Bike & Run XS,Team 3,unused,Lastname3.1,Firstname3.1,01/01/1977,f,not_member,A3,VILLENEUVE D'ASCQ TRIATHLON,1,Complet,Payé
3,Bike & Run XS,Team 3,unused,Lastname3.2,Firstname3.2,02/01/1977,h,not_member,,,1,Complet,Payé
101,Bike & Run Jeunes 12-15,Team 101,unused,Lastname101.1,Firstname101.1,01/01/2007,h,fftri,A9,VILLENEUVE D'ASCQ TRIATHLON,,Complet,Payé
101,Bike & Run Jeunes 12-15,Team 101,unused,Lastname101.2,Firstname101.2,02/01/2007,h,fftri,A9,VILLENEUVE D'ASCQ TRIATHLON,,Complet,Payé
201,Bike & Run Jeunes 6-11,Team 201,unused,Lastname201.1,Firstname201.1,01/01/2014,h,fftri,A9,LILLE TRIATHLON,,Complet,Payé
201,Bike & Run Jeunes 6-11,Team 201,unused,Lastname201.2,Firstname201.2,02/01/2014,h,fftri,A9,LILLE TRIATHLON,,Complet,Payé`

var _ = Describe("import teams", func() {

	It("should import teams of a selected race", func() {
		// given
		source, err := os.CreateTemp(os.TempDir(), "teams*.csv")
		Expect(err).NotTo(HaveOccurred())
		source.WriteString(csv)
		source.Close()
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
