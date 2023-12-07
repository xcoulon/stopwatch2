/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log/slog"
	"math"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

func NewGenerateReportCmd() *cobra.Command {
	generateCmd := &cobra.Command{
		Use:   "generate-report <race_name> <teams.yaml> <timings.yaml> <output-dir>",
		Short: "Generate a race report",
		Args:  cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {
			logger := NewLogger(debug)
			_, _, err := GenerateReport(logger, args[0], args[1], args[2], args[3])
			return err
		},
	}
	return generateCmd
}

func GenerateReport(logger *slog.Logger, raceName, teamFilename, timingFilename, outputDir string) (string, string, error) {
	results, err := NewOverallResults(teamFilename, timingFilename)
	if err != nil {
		return "", "", err
	}
	base := filepath.Base(timingFilename)[:len(filepath.Base(timingFilename))-len(filepath.Ext(timingFilename))]
	overallFilename := filepath.Join(outputDir, base+"-overall.adoc")
	if err := checkOutputFile(overallFilename); err != nil {
		return "", "", err
	}
	if err := GenerateOverallResultsReport(logger, raceName, results, overallFilename); err != nil {
		return "", "", err
	}
	perCategoryFilename := filepath.Join(outputDir, base+"-per-category.adoc")
	if err := checkOutputFile(perCategoryFilename); err != nil {
		return "", "", err
	}

	if err := GenerateResultsPerCategoryReport(logger, raceName, results, perCategoryFilename); err != nil {
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
			return nil, fmt.Errorf("unable to decode teams: %w", err)
		}
		teams[team.BibNumber] = team
	}
	return teams, nil
}

func NewOverallResults(teamFilename, timingFilename string) ([]TeamResult, error) {
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
		return nil, fmt.Errorf("unable to decode timings: %w", err)
	}
	// start time: assume it's the first entry
	if !timings[0].IsStart() {
		return nil, fmt.Errorf("missing start time")
	}
	startTime, err := timings[0].Time()
	if err != nil {
		return nil, fmt.Errorf("invalid start time: %w", err)
	}
	results := []TeamResult{}
	for i, t := range timings[1:] {
		bib, err := t.BibNumber()
		if err != nil {
			return nil, fmt.Errorf("invalid bib number: %w", err)
		}
		arrivalTime, err := t.Time()
		if err != nil {
			return nil, fmt.Errorf("invalid arrival time: %w", err)
		}
		team, found := teams[bib]
		if !found {
			return nil, fmt.Errorf("no team with number %d", bib)
		}
		results = append(results, TeamResult{
			Rank:      i + 1,
			Team:      team,
			TotalTime: arrivalTime.Sub(startTime),
		})
	}
	return results, nil
}

func getCategory(ageCategory, gender string) string {
	return fmt.Sprintf("%s/%s", ageCategory, gender)
}

func getMemberName(member TeamMember) string {
	return fmt.Sprintf("%s %s %s", member.FirstName, member.LastName, member.Club)
}

func getMemberNames(members []TeamMember) string {
	return fmt.Sprintf("%s %s \n %s %s", members[0].FirstName, members[0].LastName, members[1].FirstName, members[1].LastName)
}

func getMemberClubs(members []TeamMember) string {
	if members[0].Club == members[1].Club {
		return members[0].Club
	}
	if members[0].Club == "" {
		return members[1].Club
	}
	if members[1].Club == "" {
		return members[0].Club
	}
	return strings.TrimSpace(fmt.Sprintf("%s + %s", members[0].Club, members[1].Club))
}

func GenerateOverallResultsReport(logger *slog.Logger, raceName string, results []TeamResult, outputFilename string) error {
	outputFile, err := os.Create(outputFilename)
	if err != nil {
		return err
	}
	defer outputFile.Close()

	if len(results) == 0 {
		return fmt.Errorf("empty results?")
	}

	logger.Info("generating overall results...", "race_name", raceName, "filename", outputFilename)

	adocWriter := bufio.NewWriter(outputFile)
	adocWriter.WriteString(fmt.Sprintf("== %s - Classement Général\n\n", raceName)) //nolint:errcheck
	// table header
	adocWriter.WriteString("[cols=\"1,5,5,10,10,5\"]\n")                             //nolint:errcheck
	adocWriter.WriteString("|===\n")                                                 //nolint:errcheck
	adocWriter.WriteString("|# |Equipe |Catégorie |Coureur 1 |Coureur 2 |Temps\n\n") //nolint:errcheck // |Club

	// table rows
	for i, r := range results {
		adocWriter.WriteString(fmt.Sprintf("|%d |%s |%s |%s |%s |%s \n", //nolint:errcheck
			i+1,
			fmt.Sprintf("%d %s", r.BibNumber, r.Name),
			getCategory(r.AgeCategory, r.Gender),
			getMemberName(r.Members[0]),
			getMemberName(r.Members[1]),
			r.TotalTime.Round(time.Second).String()))
	}
	// close table
	adocWriter.WriteString("|===\n") //nolint:errcheck
	err = adocWriter.Flush()
	if err != nil {
		return errors.Wrap(err, "unable to generate overall results")
	}
	return nil
}

func rankPerCategory(results []TeamResult) (map[string][]TeamResult, error) {
	if len(results) == 0 {
		return nil, fmt.Errorf("empty results?")
	}
	resultsPerCategory := map[string][]TeamResult{}

	for _, c := range []string{MiniPoussin, Poussin, Pupille, Benjamin, Minime, Cadet, Junior, Senior, Master} {
		for _, g := range []string{"F", "H", "M"} {
			// retain 3 first match
			for _, r := range results {
				if r.AgeCategory == c && r.Gender == g {
					cat := getCategory(r.AgeCategory, r.Gender)
					if resultsPerCategory[cat] == nil {
						resultsPerCategory[cat] = []TeamResult{}
					}
					resultsPerCategory[cat] = append(resultsPerCategory[cat], r)
				}
			}
		}
	}
	return resultsPerCategory, nil
}

func GenerateResultsPerCategoryReport(logger *slog.Logger, raceName string, results []TeamResult, outputFilename string) error {
	logger.Info("generating results per category...", "race_name", raceName, "filename", outputFilename)
	winners, err := rankPerCategory(results)
	if err != nil {
		return err
	}
	outputFile, err := os.Create(outputFilename)
	if err != nil {
		return err
	}
	defer outputFile.Close()

	adocWriter := bufio.NewWriter(outputFile)
	adocWriter.WriteString(fmt.Sprintf("== %s - Classement Par Catégorie\n\n", raceName)) //nolint:errcheck

	for _, c := range []string{MiniPoussin, Poussin, Pupille, Benjamin, Minime, Cadet, Junior, Senior, Master} {
	gender_loop:
		for _, g := range []string{"F", "H", "M"} {
			if rs, found := winners[getCategory(c, g)]; found {
				// section title
				adocWriter.WriteString(fmt.Sprintf("=== %s\n\n", getCategory(rs[0].AgeCategory, rs[0].Gender))) //nolint:errcheck

				// table header
				adocWriter.WriteString("[cols=\"1,10,10,10,5\"]\n")                   //nolint:errcheck
				adocWriter.WriteString("|===\n")                                      //nolint:errcheck
				adocWriter.WriteString("|# |Equipe |Coureur 1 |Coureur 2 |Temps\n\n") //nolint:errcheck  // |Club
				l := int(math.Min(3, float64(len(rs))))
				for i := 0; i < l; i++ {
					r := rs[i]
					adocWriter.WriteString(fmt.Sprintf("|%d |%s |%s |%s |%s\n", //nolint:errcheck // |%s
						r.Rank,
						fmt.Sprintf("%d %s", r.BibNumber, r.Name),
						getMemberName(r.Members[0]),
						getMemberName(r.Members[1]),
						r.TotalTime.Round(time.Second).String()))
				}
				adocWriter.WriteString("|===\n\n") //nolint:errcheck
				continue gender_loop
			}
		}
	}
	err = adocWriter.Flush()
	if err != nil {
		return errors.Wrap(err, "unable to generate results per category")
	}
	return nil
}
