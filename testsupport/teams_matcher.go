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

func MatchTeams(expected []cmd.Team) types.GomegaMatcher {
	return &TeamsMatcher{
		expected: expected,
	}
}

type TeamsMatcher struct {
	expected []cmd.Team
	diffs    string
}

func (m *TeamsMatcher) Match(actual interface{}) (success bool, err error) {
	if _, ok := actual.([]cmd.Team); !ok {
		return false, errors.Errorf("MatchTeams matcher expects a '[]cmd.Team' (actual: %T)", actual)
	}
	if diff := cmp.Diff(m.expected, actual); diff != "" {
		if log.IsLevelEnabled(log.DebugLevel) {
			log.Debugf("actual teams:\n%s", spew.Sdump(actual))
			log.Debugf("expected teams:\n%s", spew.Sdump(m.expected))
		}
		m.diffs = diff
		return false, nil
	}
	return true, nil
}

func (m *TeamsMatcher) FailureMessage(_ interface{}) (message string) {
	return fmt.Sprintf("expected teams to match:\n%s", m.diffs)
}

func (m *TeamsMatcher) NegatedFailureMessage(_ interface{}) (message string) {
	return fmt.Sprintf("expected teams not to match:\n%s", m.diffs)
}
