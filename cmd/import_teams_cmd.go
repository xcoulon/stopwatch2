package cmd

import (
	"encoding/csv"
	"io"
	"log/slog"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

func NewImportTeamsCmd() *cobra.Command {
	importCmd := &cobra.Command{
		Use:   "import-teams <teams.csv> <teams.yaml>",
		Short: "Import Teams from a CSV file",
		Args:  cobra.ExactArgs(2),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if !force {
				return checkOutputFile(args[1])
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			logger := NewLogger(debug)
			return ImportCSV(logger, args[0], args[1])
		},
	}
	return importCmd
}

// ImportFromFile imports the data from the given file
func ImportCSV(logger *slog.Logger, sourceFilename, outputFilename string) error {
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
	// r.Comma = ';'
	undefinedMember := TeamMember{}
	member1 := undefinedMember
	var bibNumber int
	var teamName string
	line := 0
	for {
		line++
		logger.Info("reading record", "line", line)
		record, err := r.Read()
		if errors.Is(err, io.EOF) {
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
			member1, teamName, err = newTeamMember(record)
			if err != nil {
				return errors.Wrapf(err, "unable to create team member #1 from %+v", record)
			}
			bibNumber, err = strconv.Atoi(record[0])
			if err != nil {
				return errors.Wrapf(err, "unable to convert bibNumber '%s' to a number", record[1])
			}
		} else {
			// in some cases of export from .xls, only the member1 has the bibnumber
			if len(record) < len(headers) {
				record2 := make([]string, len(record)+1)
				copy(record2[1:], record)
				record = record2
			}
			member2, _, err := newTeamMember(record)
			if err != nil {
				return errors.Wrapf(err, "unable to create team member #2 from %+v", record)
			}

			team := Team{
				Name:        teamName, // record[2], // team name
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
			logger.Debug("team", "yaml", string(out))
			if _, err := output.WriteString("---\n" + string(out)); err != nil {
				return errors.Wrapf(err, "unable to write team from %+v", team)
			}
			// reset
			member1 = undefinedMember
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

func newTeamMember(record []string) (TeamMember, string, error) {
	dateOfBirth, err := time.Parse("02/01/2006", record[5])
	if err != nil {
		return TeamMember{}, "", errors.Wrapf(err, "unable to parse date '%s'", record[5])
	}
	return TeamMember{
		LastName:    record[3],
		FirstName:   record[4],
		DateOfBirth: ISO8601Date(dateOfBirth),
		Gender:      strings.ToUpper(record[6]),
		Category:    GetAgeCategory(dateOfBirth),
		Club:        record[8],
	}, record[2], nil
}
