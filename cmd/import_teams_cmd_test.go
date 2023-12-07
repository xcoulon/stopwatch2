package cmd_test

import (
	"bytes"
	"context"
	"errors"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/alecthomas/assert"
	"github.com/charmbracelet/log"
	"github.com/google/go-cmp/cmp"
	"github.com/sanity-io/litter"
	"github.com/stretchr/testify/require"
	"github.com/vatriathlon/stopwatch2/cmd"
	"gopkg.in/yaml.v3"
)

func TestImportCSV(t *testing.T) {

	// given
	source, err := os.Open("teams.csv")
	require.NoError(t, err)
	outputDir, err := os.MkdirTemp(os.TempDir(), "bikerun2023-")
	require.NoError(t, err)
	outputFilename := filepath.Join(outputDir, "teams.yaml")
	require.NoError(t, err)
	logger := cmd.NewLogger(false)

	// when
	err = cmd.ImportCSV(logger, source.Name(), outputFilename)
	// then

	require.NoError(t, err)
	expected := []cmd.Team{
		cmd.NewTeam("Team 1", "H", "Master", 1,
			cmd.NewTeamMember("Firstname1.1", "Lastname1.1", parseDate(t, "1977-01-01"), "Master", "H", ""),
			cmd.NewTeamMember("Firstname1.2", "Lastname1.2", parseDate(t, "1977-01-02"), "Master", "H", ""),
		),
		cmd.NewTeam("Team 2", "F", "Master", 2,
			cmd.NewTeamMember("Firstname2.1", "Lastname2.1", parseDate(t, "1977-01-01"), "Master", "F", "LILLE TRIATHLON"),
			cmd.NewTeamMember("Firstname2.2", "Lastname2.2", parseDate(t, "1977-01-02"), "Master", "F", ""),
		),
		cmd.NewTeam("Team 3", "M", "Master", 3,
			cmd.NewTeamMember("Firstname3.1", "Lastname3.1", parseDate(t, "1977-01-01"), "Master", "F", "VILLENEUVE D'ASCQ TRIATHLON"),
			cmd.NewTeamMember("Firstname3.2", "Lastname3.2", parseDate(t, "1977-01-02"), "Master", "H", ""),
		),
		cmd.NewTeam("Team 101", "H", "Minime", 101,
			cmd.NewTeamMember("Firstname101.1", "Lastname101.1", parseDate(t, "2007-01-01"), "Minime", "H", "VILLENEUVE D'ASCQ TRIATHLON"),
			cmd.NewTeamMember("Firstname101.2", "Lastname101.2", parseDate(t, "2007-01-02"), "Minime", "H", "VILLENEUVE D'ASCQ TRIATHLON"),
		),
		cmd.NewTeam("Team 201", "H", "Poussin", 201,
			cmd.NewTeamMember("Firstname201.1", "Lastname201.1", parseDate(t, "2014-01-01"), "Poussin", "H", "LILLE TRIATHLON"),
			cmd.NewTeamMember("Firstname201.2", "Lastname201.2", parseDate(t, "2014-01-02"), "Poussin", "H", "LILLE TRIATHLON"),
		),
	}

	assert.Condition(t, containsTeams(logger, expected, outputFilename))
}

func parseDate(t *testing.T, d string) cmd.ISO8601Date {
	r, err := time.Parse("2006-01-02", d)
	require.NoError(t, err)
	return cmd.ISO8601Date(r)
}

func containsTeams(logger *slog.Logger, expected []cmd.Team, filename string) assert.Comparison {
	return func() (success bool) {

		content, err := os.ReadFile(filename)
		if err != nil {
			logger.Error("failed to read file", "filename", filename, "error", err)
			return false
		}
		actual := []cmd.Team{}
		decoder := yaml.NewDecoder(bytes.NewReader(content))
		// decode 1 team at a time
		for {
			team := cmd.Team{}
			if err := decoder.Decode(&team); errors.Is(err, io.EOF) {
				break
			} else if err != nil {
				logger.Error("failed to decode YAML contents", "filename", filename, "error", err)
				return false
			}
			actual = append(actual, team)
		}

		if diff := cmp.Diff(expected, actual); diff != "" {
			if logger.Enabled(context.Background(), slog.LevelDebug) {
				log.Debugf("actual teams:\n%s", litter.Sdump(actual))
				log.Debugf("expected teams:\n%s", litter.Sdump(expected))
			}
			diff = strings.ReplaceAll(diff, "\u00a0", "")
			diff = strings.ReplaceAll(diff, "\t", "  ")
			logger.Info("contents are not equal", "diff", diff)
			return false
		}
		return true
	}
}
