package cmd

import (
	"encoding/csv"
	"io"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

func NewImportTeamsCmd() *cobra.Command {
	importCmd := &cobra.Command{
		Use:   "import-teams <source_csv> <output_yaml>",
		Short: "Import Teams from a CSV file",
		Args:  cobra.ExactArgs(2),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if debug {
				logrus.SetLevel(logrus.DebugLevel)
			}
			if !force {
				return checkOutputFile(args[1])
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return ImportCSV(args[0], args[1])
		},
	}
	return importCmd
}

// ImportFromFile imports the data from the given file
func ImportCSV(sourceFilename, outputFilename string) error {
	var headers []string
	source, err := os.Open(sourceFilename)
	if err != nil {
		return err
	}
	defer source.Close()
	output, err := os.Create(outputFilename)
	if err != nil {
		return err
	}
	defer output.Close()
	r := csv.NewReader(source)
	undefinedMember := TeamMember{}
	member1 := undefinedMember
	member2 := undefinedMember

	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		if headers == nil {
			headers = record
			continue
		}
		if member1 == undefinedMember {
			member1, err = newTeamMember(record)
			if err != nil {
				return errors.Wrapf(err, "unable to create team member #1 from %+v", record)
			}
		} else {
			member2, err = newTeamMember(record)
			if err != nil {
				return errors.Wrapf(err, "unable to create team member #2 from %+v", record)
			}

			bibNumber, err := strconv.Atoi(record[0])
			if err != nil {
				return errors.Wrapf(err, "unable to convert bibNumber '%s' to a number", record[1])
			}
			team := Team{
				Name:        record[2], // team name
				AgeCategory: GetTeamAgeCategory(member1.Category, member2.Category),
				BibNumber:   bibNumber,
				Gender:      genderFrom(member1, member2),
				Members: []TeamMember{
					member1,
					member2,
				},
			}
			out, err := yaml.Marshal(team)
			if err != nil {
				return errors.Wrapf(err, "unable to marshall team from %+v", team)
			}
			logrus.Debugf("%s", out)
			if _, err := output.WriteString("---\n" + string(out)); err != nil {
				return errors.Wrapf(err, "unable to write team from %+v", team)
			}
			// reset
			member1 = undefinedMember
			member2 = undefinedMember
		}
	}
	return nil
}

func genderFrom(teamMember1, teamMember2 TeamMember) string {
	if teamMember1.Gender == teamMember2.Gender {
		return teamMember1.Gender
	}
	return "M"
}

func newTeamMember(record []string) (TeamMember, error) {
	dateOfBirth, err := time.Parse("02/01/2006", record[6])
	if err != nil {
		return TeamMember{}, errors.Wrapf(err, "unable to parse date '%s'", record[5])
	}
	return TeamMember{
		LastName:    record[4],
		FirstName:   record[5],
		DateOfBirth: ISO8601Date(dateOfBirth),
		Gender:      strings.ToUpper(record[7]),
		Category:    GetAgeCategory(dateOfBirth),
		Club:        record[10],
	}, nil
}
