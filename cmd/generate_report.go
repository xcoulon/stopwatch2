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
	"path/filepath"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

func NewGenerateReportCmd() *cobra.Command {
	var race string
	var teams string
	var timings string
	var outputDir string
	generateCmd := &cobra.Command{
		Use:   "generate-report --race-name=<race_name> --teams=<teams.yaml> --race-results=<results.yaml> --output-dir=<directory>",
		Short: "Generate a race report",
		Args:  cobra.ExactArgs(0),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if debug {
				logrus.SetLevel(logrus.DebugLevel)
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			_, _, err := GenerateReport(race, teams, timings, outputDir)
			return err
		},
	}
	generateCmd.Flags().StringVar(&race, "race-name", "", "Race name")
	generateCmd.Flags().StringVar(&teams, "teams", "", "File describing the teams (YAML)")
	generateCmd.Flags().StringVar(&timings, "timings", "", "File containing the timings (YAML)")
	generateCmd.Flags().StringVar(&outputDir, "output-dir", "", "Output dir (AsciiDoc)")
	return generateCmd
}

func GenerateReport(raceName, teamFilename, timingFilename, outputDir string) (string, string, error) {
	results, err := NewTeamResults(teamFilename, timingFilename)
	if err != nil {
		return "", "", err
	}
	base := filepath.Base(timingFilename)[:len(filepath.Base(timingFilename))-len(filepath.Ext(timingFilename))]
	overallFilename := filepath.Join(outputDir, base+"-overall.adoc")
	if err := checkOutputFile(overallFilename); err != nil {
		return "", "", err
	}
	if err := generateOverallResults(raceName, results, overallFilename); err != nil {
		return "", "", err
	}
	perCategoryFilename := filepath.Join(outputDir, base+"-per-category.adoc")
	if err := checkOutputFile(perCategoryFilename); err != nil {
		return "", "", err
	}

	if err := generateResultsPerCategory(raceName, results, perCategoryFilename); err != nil {
		return "", "", err
	}
	return overallFilename, perCategoryFilename, nil
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
			return nil, fmt.Errorf("unable to decode teams: %v", err)
		}
		teams[team.BibNumber] = team
	}
	return teams, nil
}

func NewTeamResults(teamFilename, timingFilename string) ([]TeamResult, error) {
	teams, err := loadTeams(teamFilename)
	if err != nil {
		return nil, err
	}

	timingFile, err := os.Open(timingFilename)
	if err != nil {
		return nil, err
	}
	timings := Timings{}
	decoder := yaml.NewDecoder(timingFile)
	if err := decoder.Decode(&timings); err != nil {
		return nil, fmt.Errorf("unable to decode timings: %v", err)
	}
	// start time: assume it's the first entry
	if !timings[0].IsStart() {
		return nil, fmt.Errorf("missing start time")
	}
	startTime, err := timings[0].Time()
	if err != nil {
		return nil, fmt.Errorf("invalid start time: %v", err)
	}
	results := []TeamResult{}
	for _, t := range timings[1:] {
		bib, err := t.BibNumber()
		if err != nil {
			return nil, fmt.Errorf("invalid bib number: %v", err)
		}
		arrivalTime, err := t.Time()
		if err != nil {
			return nil, fmt.Errorf("invalid arrival time: %v", err)
		}
		team, found := teams[bib]
		if !found {
			return nil, fmt.Errorf("no team with number %d", bib)
		}
		results = append(results, TeamResult{
			Team:      team,
			TotalTime: arrivalTime.Sub(startTime),
		})
	}
	return results, nil
}

func getCategory(ageCategory, gender string) string {
	return fmt.Sprintf("%s/%s", ageCategory, gender)
}

func getMemberNames(members []TeamMember) string {
	return fmt.Sprintf("%s - %s", members[0].LastName, members[1].LastName)
}

func getMemberClubs(members []TeamMember) string {

	if members[0].Club == members[1].Club {
		return members[0].Club
	}
	return strings.TrimSpace(fmt.Sprintf("%s %s", members[0].Club, members[1].Club))
}

func generateOverallResults(raceName string, results []TeamResult, outputFilename string) error {
	outputFile, err := os.Create(outputFilename)
	if err != nil {
		return err
	}
	defer outputFile.Close()

	if len(results) == 0 {
		return fmt.Errorf("empty results?")
	}

	logrus.WithField("race_name", raceName).WithField("filename", outputFilename).Info("generating overall results...")

	adocWriter := bufio.NewWriter(outputFile)
	adocWriter.WriteString(fmt.Sprintf("= %s - Classement Général\n\n", raceName))
	// table header
	adocWriter.WriteString("[cols=\"2,5,5,5,10,10,5\"]\n")
	adocWriter.WriteString("|===\n")
	adocWriter.WriteString("|# |Dossard |Equipe |Catégorie |Coureurs |Club |Temps Total\n\n")

	// table rows
	for i, r := range results {
		adocWriter.WriteString(fmt.Sprintf("|%d |%d |%s |%s |%s |%s |%s \n",
			i+1,
			r.BibNumber,
			r.Name,
			getCategory(r.AgeCategory, r.Gender),
			getMemberNames(r.Members),
			getMemberClubs(r.Members),
			r.TotalTime.Round(time.Second).String()))
	}
	// close table
	adocWriter.WriteString("|===\n")
	err = adocWriter.Flush()
	if err != nil {
		return errors.Wrap(err, "unable to generate overall results")
	}
	return nil
}

func generateResultsPerCategory(raceName string, results []TeamResult, outputFilename string) error {
	outputFile, err := os.Create(outputFilename)
	if err != nil {
		return err
	}
	defer outputFile.Close()

	if len(results) == 0 {
		return fmt.Errorf("empty results?")
	}

	logrus.WithField("race_name", raceName).WithField("filename", outputFilename).Info("generating results per category...")

	adocWriter := bufio.NewWriter(outputFile)
	adocWriter.WriteString(fmt.Sprintf("= %s - Classement Par Catégorie\n\n", raceName))

	for _, c := range []string{MiniPoussin, Poussin, Pupille, Benjamin, Minime, Cadet, Junior, Senior, Master} {
		for _, g := range []string{"F", "H", "M"} {
			// retain 1st match
			for _, r := range results {
				if r.AgeCategory == c && r.Gender == g {
					// section title
					adocWriter.WriteString(fmt.Sprintf("== %s\n\n", getCategory(r.AgeCategory, r.Gender)))

					// table header
					adocWriter.WriteString("[cols=\"2,5,5,10,10,5\"]\n")
					adocWriter.WriteString("|===\n")
					adocWriter.WriteString("|# |Dossard |Equipe |Coureurs |Club |Temps Total\n\n")
					adocWriter.WriteString(fmt.Sprintf("|%d |%d |%s |%s |%s |%s \n",
						1,
						r.BibNumber,
						r.Name,
						getMemberNames(r.Members),
						getMemberClubs(r.Members),
						r.TotalTime.Round(time.Second).String()))
					adocWriter.WriteString("|===\n\n")
				}
				continue
			}
		}
	}
	err = adocWriter.Flush()
	if err != nil {
		return errors.Wrap(err, "unable to generate results per category")
	}
	return nil
}
