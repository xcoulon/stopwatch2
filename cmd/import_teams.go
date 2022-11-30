package cmd

import (
	"encoding/csv"
	"io"
	"math"
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
	var sourceFilename string
	var outputFilename string
	importCmd := &cobra.Command{
		Use:   "import-teams --race=<race> --source=<source_csv> --output=<output_yaml>",
		Short: "Import Teams from a CSV file",
		Args:  cobra.ExactArgs(0),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if debug {
				logrus.SetLevel(logrus.DebugLevel)
			}
			if !force {
				return checkOutputFile(outputFilename)
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return ImportCSV(sourceFilename, outputFilename)
		},
	}
	importCmd.Flags().StringVar(&sourceFilename, "source", "", "Source file (CSV)")
	importCmd.Flags().StringVar(&outputFilename, "output", "", "Output file (YAML)")
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

			bibNumber, err := strconv.Atoi(record[1])
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
	dateOfBirth, err := time.Parse("02/01/2006", record[5])
	if err != nil {
		return TeamMember{}, errors.Wrapf(err, "unable to parse date '%s'", record[5])
	}
	return TeamMember{
		LastName:    record[3],
		FirstName:   record[4],
		DateOfBirth: ISO8601Date(dateOfBirth),
		Gender:      strings.ToUpper(record[6]),
		Category:    GetAgeCategory(dateOfBirth),
		Club:        record[9],
	}, nil
}

const (
	// MiniPoussin 2015 à 2016
	MiniPoussin = "Mini-poussin"
	// Poussin 		2013 à 2014
	Poussin = "Poussin"
	// Pupille 		2011 à 2012
	Pupille = "Pupille"
	// Benjamin 	2009 à 2010
	Benjamin = "Benjamin"
	// Minime 		2007 à 2008
	Minime = "Minime"
	// Cadet 		2005 à 2006
	Cadet = "Cadet"
	// Junior 		2003 à 2004
	Junior = "Junior"
	// Senior 	 	1983 à 2002
	Senior = "Senior"
	// Master 		1955 à 1982
	Master = "Master"
)

// GetAgeCategory gets the age category associated with the given date of birth
func GetAgeCategory(dateOfBirth time.Time) string {
	yearOfBirth := dateOfBirth.Year()
	logrus.WithField("year_of_birth", yearOfBirth).Debug("computing age category")
	switch {
	case yearOfBirth == 2015 || yearOfBirth == 2016:
		return MiniPoussin
	case yearOfBirth == 2013 || yearOfBirth == 2014:
		return Poussin
	case yearOfBirth == 2011 || yearOfBirth == 2012:
		return Pupille
	case yearOfBirth == 2009 || yearOfBirth == 2010:
		return Benjamin
	case yearOfBirth == 2007 || yearOfBirth == 2008:
		return Minime
	case yearOfBirth == 2005 || yearOfBirth == 2006:
		return Cadet
	case yearOfBirth == 2003 || yearOfBirth == 2004:
		return Junior
	case yearOfBirth >= 1983 && yearOfBirth <= 2002:
		return Senior
	default:
		return Master
	}

}

var ageCategories = map[string]int{
	MiniPoussin: 1,
	Poussin:     2,
	Pupille:     3,
	Benjamin:    4,
	Minime:      5,
	Cadet:       6,
	Junior:      7,
	Master:      8,
	Senior:      9, // so that max(master, senior) -> senior
}

// GetTeamAgeCategory computes the age category for the team
func GetTeamAgeCategory(ageCategory1, ageCategory2 string) string {
	cat1 := ageCategories[ageCategory1]
	cat2 := ageCategories[ageCategory2]
	// assign to senior if 1 veteran + 1 under senior
	if (ageCategory1 == Master && cat2 <= ageCategories[Junior]) || (ageCategory2 == Master && cat1 <= ageCategories[Junior]) {
		return Senior
	}
	teamAgeCategoryValue := math.Max(float64(cat1), float64(cat2))
	logrus.WithField("team_age_category_value", teamAgeCategoryValue).Debugf("computing team age category...")
	//

	for k, v := range ageCategories {
		if float64(v) == teamAgeCategoryValue {
			return k
		}
	}
	return ""

}
