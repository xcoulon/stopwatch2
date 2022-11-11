/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"bufio"
	"bytes"
	"fmt"
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

func newGenerateReportCmd() *cobra.Command {
	var race string
	var teams string
	var results string
	var output string
	cmd := &cobra.Command{
		Use:   "generate-report --race=<race_name> --teams=<teams> --results=<results> --output=<output>",
		Short: "Generate a race report",
		Args:  cobra.ExactArgs(4),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if debug {
				logrus.SetLevel(logrus.DebugLevel)
			}
			if !force {
				return checkOutputFile(output)
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return GenerateReport(race, teams, results, output)
		},
	}
	cmd.Flags().StringVar(&race, "race", "", "Race name")
	cmd.Flags().StringVar(&teams, "teams", "", "File describing the teams (YAML)")
	cmd.Flags().StringVar(&results, "results", "", "File containing raw results (YAML")
	cmd.Flags().StringVar(&output, "output", "", "Output file (AsciiDoc)")
	return cmd
}

type TeamResult struct {
	bibNumber string
	name      string
	category  string
	members   string
	club      string
	totalTime time.Duration
}

func GenerateReport(raceName, teamFilename, raceFilename, outputFilename string) error {
	results, err := loadResults(raceFilename, teamFilename)
	if err != nil {
		return err
	}
	return generateAsciidoc(raceName, results, outputFilename)
	// // by age and gender
	// ageCategories := []string{Poussin, Pupille, Benjamin, Minime, Cadet, Junior, Senior, Master}
	// genders := []string{"H", "F", "M"}
	// for _, ageCategory := range ageCategories {
	// 	for _, gender := range genders {
	// 		categoryRows, err := s.baseService.db.Raw(byGenderAndAgeQuery, race.ID, ageCategory, gender).Rows()
	// 		defer categoryRows.Close()
	// 		if err != nil {
	// 			return nil, errors.Wrap(err, "unable to generate results")
	// 		}
	// 		file, err = generateAsciidoc(outputDir, race, categoryRows, ageCategory, gender, false)
	// 		if err != nil {
	// 			return nil, errors.Wrap(err, "unable to generate results")
	// 		}
	// 		if file != "" {
	// 			files = append(files, file)
	// 		}
	// 	}
	// }
	// return files, nil
}

func loadTeams(filename string) (map[int]Team, error) {
	content, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	teams := map[int]Team{}
	decoder := yaml.NewDecoder(bytes.NewReader(content))
	// decode 1 team at a time
	for {
		team := Team{}
		if err := decoder.Decode(&team); errors.Is(err, io.EOF) {
			break
		} else if err != nil {
			return nil, err
		}
		teams[team.BibNumber] = team
	}
	return teams, nil
}

func loadResults(raceFilename, teamFilename string) ([]TeamResult, error) {
	teams, err := loadTeams(teamFilename)
	if err != nil {
		return nil, err
	}

	raceFile, err := os.Open(raceFilename)
	if err != nil {
		return nil, err
	}
	raceResults := RaceResults{}
	decoder := yaml.NewDecoder(raceFile)
	if err := decoder.Decode(&raceResults); err != nil {
		return nil, err
	}
	results := []TeamResult{}
	for _, t := range raceResults.Teams {
		team, found := teams[t.BibNumber]
		if !found {
			return nil, fmt.Errorf("no team with number %d", t.BibNumber)
		}
		results = append(results, TeamResult{
			bibNumber: strconv.Itoa(t.BibNumber),
			name:      team.Name,
			category:  getCategory(team.AgeCategory, team.Gender),
			club:      getMemberClubs(team.Members[0], team.Members[1]),
			members:   getMemberNames(team.Members[0], team.Members[1]),
			totalTime: t.ScratchTime.Sub(raceResults.StartTime),
		})
	}
	return results, nil
}

func getCategory(ageCategory, gender string) string {
	return fmt.Sprintf("%s/%s", string([]rune(ageCategory)[0]), string([]rune(gender)[0]))
}

func getMemberNames(member1, member2 TeamMember) string {
	return fmt.Sprintf("%s - %s", member1.LastName, member2.LastName)
}

func getMemberClubs(member1, member2 TeamMember) string {
	if member1.Club == member2.Club {
		return member1.Club
	}
	return strings.TrimSpace(fmt.Sprintf("%s %s", member1.Club, member2.Club))
}

func generateAsciidoc(raceName string, results []TeamResult, outputFilename string) error {
	outputFile, err := os.Create(outputFilename)
	if err != nil {
		return err
	}
	defer outputFile.Close()

	if len(results) == 0 {
		return fmt.Errorf("empty results?")
	}

	logrus.WithField("race_name", raceName).
		Info("generating results...")

	adocWriter := bufio.NewWriter(outputFile)
	adocWriter.WriteString(fmt.Sprintf("= Classement %s\n\n", raceName))
	adocWriter.WriteString(fmt.Sprintf("== Classement %s\n\n", "scratch"))
	// table header
	adocWriter.WriteString("[cols=\"2,5,5,")
	adocWriter.WriteString("5,")
	adocWriter.WriteString("8,8,4\"]\n")
	adocWriter.WriteString("|===\n")
	adocWriter.WriteString("|# |Dossard ")
	adocWriter.WriteString("|Equipe ")
	adocWriter.WriteString("|Catégorie ")
	adocWriter.WriteString("|Coureurs |Club |Temps Total\n\n")

	// table rows
	for i, r := range results {
		adocWriter.WriteString(fmt.Sprintf("|%d |%s |%s ",
			i+1,
			r.bibNumber,
			r.name))
		adocWriter.WriteString(fmt.Sprintf("|%s ",
			r.category))
		adocWriter.WriteString(fmt.Sprintf("|%s |%s |%s \n",
			r.members,
			r.club,
			r.totalTime.Round(time.Second).String()))
	}
	// close table
	adocWriter.WriteString("|===\n")
	err = adocWriter.Flush()
	if err != nil {
		return errors.Wrap(err, "unable to generate results in Asciidoc")
	}
	return nil
}

func label(cat1, cat2 string) string {
	// "Scratch" and "Challenge Entreprise"
	if cat2 == "" {
		return cat1
	}
	// other: age / gender
	switch cat2 {
	case "M":
		return fmt.Sprintf("%ss / Mixte", cat1)
	case "F":
		return fmt.Sprintf("%ss / Femmes", cat1)
	default:
		return fmt.Sprintf("%ss / Hommes", cat1)
	}
}
