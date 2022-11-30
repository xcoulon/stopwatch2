package testsupport

import (
	"fmt"

	"github.com/vatriathlon/stopwatch2/cmd"

	"github.com/davecgh/go-spew/spew"
	"github.com/google/go-cmp/cmp"
	"github.com/onsi/gomega/types"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

func MatchTeamResults(expected []cmd.TeamResult) types.GomegaMatcher {
	return &TeamResultsMatcher{
		expected: expected,
	}
}

type TeamResultsMatcher struct {
	expected []cmd.TeamResult
	diffs    string
}

func (m *TeamResultsMatcher) Match(actual interface{}) (success bool, err error) {
	if _, ok := actual.([]cmd.TeamResult); !ok {
		return false, errors.Errorf("MatchTeams matcher expects a '[]cmd.TeamResult' (actual: %T)", actual)
	}
	if diff := cmp.Diff(m.expected, actual); diff != "" {
		if log.IsLevelEnabled(log.DebugLevel) {
			log.Debugf("actual team results:\n%s", spew.Sdump(actual))
			log.Debugf("expected team results:\n%s", spew.Sdump(m.expected))
		}
		m.diffs = diff
		return false, nil
	}
	return true, nil
}

func (m *TeamResultsMatcher) FailureMessage(_ interface{}) (message string) {
	return fmt.Sprintf("expected team results to match:\n%s", m.diffs)
}

func (m *TeamResultsMatcher) NegatedFailureMessage(_ interface{}) (message string) {
	return fmt.Sprintf("expected teams results not to match:\n%s", m.diffs)
}
