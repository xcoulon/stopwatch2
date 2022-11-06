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

// importTeamsCmd represents the importTeams command
var importTeamsCmd = &cobra.Command{
	Use:   "import-teams --race=<race> --source=<source_csv> --output=<output_yaml>",
	Short: "Import Teams from a CSV file",
	Args:  cobra.ExactArgs(0),
	RunE: func(cmd *cobra.Command, args []string) error {
		if debug {
			logrus.SetLevel(logrus.DebugLevel)
		}
		if !force {
			if err := checkOutputFile(output); err != nil {
				return err
			}
		}
		return ImportCSV(race, source, output)
	},
}

var race, source string

func init() {
	rootCmd.AddCommand(importTeamsCmd)
	importTeamsCmd.Flags().StringVar(&race, "race", "", "Race name")
	importTeamsCmd.Flags().StringVar(&source, "source", "", "Source CSV file")
	importTeamsCmd.Flags().StringVar(&output, "output", "", "Output YAML file")
	importTeamsCmd.Flags().BoolVar(&debug, "debug", false, "Display debug logs")
	importTeamsCmd.Flags().BoolVar(&force, "force", false, "Force-write in output file even if it exists (existing content will be lost)")
}

// ImportFromFile imports the data from the given file
func ImportCSV(race, source, output string) error {

	var headers []string
	sourceFile, err := os.Open(source)
	if err != nil {
		return err
	}
	defer sourceFile.Close()
	outputFile, err := os.Create(output)
	if err != nil {
		return err
	}
	defer outputFile.Close()
	r := csv.NewReader(sourceFile)
	undefinedMember := TeamMember{}
	teamMember1 := undefinedMember
	teamMember2 := undefinedMember

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
		if record[0] != race {
			// skip
			continue
		}
		if teamMember1 == undefinedMember {
			teamMember1, err = newTeamMember(record)
			if err != nil {
				return errors.Wrapf(err, "unable to create team member #1 from %+v", record)
			}
		} else {
			teamMember2, err = newTeamMember(record)
			if err != nil {
				return errors.Wrapf(err, "unable to create team member #2 from %+v", record)
			}

			bibNumber, err := strconv.Atoi(record[1])
			if err != nil {
				return errors.Wrapf(err, "unable to convert bibnumber '%s' to a number", record[1])
			}
			team := Team{
				Name:      record[2], // team name
				Category:  GetTeamAgeCategory(teamMember1.Category, teamMember2.Category),
				BibNumber: bibNumber,
				Gender:    genderFrom(teamMember1, teamMember2),
				Members: []TeamMember{
					teamMember1,
					teamMember2,
				},
			}
			out, err := yaml.Marshal(team)
			if err != nil {
				return errors.Wrapf(err, "unable to marshall team from %+v", team)
			}
			logrus.Debugf("%s", out)
			if _, err := outputFile.WriteString("---\n" + string(out)); err != nil {
				return errors.Wrapf(err, "unable to write team from %+v", team)
			}
			// reset
			teamMember1 = undefinedMember
			teamMember2 = undefinedMember
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
	// MiniPoussin 2012 à 2013
	MiniPoussin = "Mini-poussin"
	// Poussin 		2010 à 2011
	Poussin = "Poussin"
	// Pupille 		2008 à 2009
	Pupille = "Pupille"
	// Benjamin 	2006 à 2007
	Benjamin = "Benjamin"
	// Minime 		2004 à 2005
	Minime = "Minime"
	// Cadet 		2002 à 2003
	Cadet = "Cadet"
	// Junior 		2000 à 2001
	Junior = "Junior"
	// Senior 	 	1980 à 1999
	Senior = "Senior"
	// Master 		1955 à 1979
	Master = "Master"
)

// GetAgeCategory gets the age category associated with the given date of birth
func GetAgeCategory(dateOfBirth time.Time) string {
	yearOfBirth := dateOfBirth.Year()
	logrus.WithField("year_of_birth", yearOfBirth).Debug("computing age category")
	switch {
	case yearOfBirth == 2012 || yearOfBirth == 2013:
		return MiniPoussin
	case yearOfBirth == 2010 || yearOfBirth == 2011:
		return Poussin
	case yearOfBirth == 2008 || yearOfBirth == 2009:
		return Pupille
	case yearOfBirth == 2006 || yearOfBirth == 2007:
		return Benjamin
	case yearOfBirth == 2004 || yearOfBirth == 2005:
		return Minime
	case yearOfBirth == 2002 || yearOfBirth == 2003:
		return Cadet
	case yearOfBirth == 2000 || yearOfBirth == 2001:
		return Junior
	case yearOfBirth >= 1980 && yearOfBirth <= 1999:
		return Senior
	default:
		return Master
	}

}

var ageCategories map[string]int

func init() {
	ageCategories = map[string]int{
		MiniPoussin: 1,
		Poussin:     2,
		Pupille:     3,
		Benjamin:    4,
		Minime:      5,
		Cadet:       6,
		Junior:      7,
		Master:      8,
		Senior:      9,
	}
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
