package cmd_test

import (
	"testing"
	"time"

	"github.com/vatriathlon/stopwatch2/cmd"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAgeCategory(t *testing.T) {
	testcases := []struct {
		category    string
		dateOfBirth string
		expected    string
	}{
		{
			category:    "mini poussin",
			dateOfBirth: "2015-02-03",
			expected:    cmd.MiniPoussin,
		},
		{
			category:    "poussin",
			dateOfBirth: "2013-02-03",
			expected:    cmd.Poussin,
		},
		{
			category:    "pupille",
			dateOfBirth: "2012-02-03",
			expected:    cmd.Pupille,
		},
		{
			category:    "benjamin",
			dateOfBirth: "2010-02-03",
			expected:    cmd.Benjamin,
		},
		{
			category:    "cadet",
			dateOfBirth: "2005-02-03",
			expected:    cmd.Cadet,
		},
		{
			category:    "junior",
			dateOfBirth: "2003-02-03",
			expected:    cmd.Junior,
		},
		{
			category:    "senior",
			dateOfBirth: "1984-02-03",
			expected:    cmd.Senior,
		},
		{
			category:    "Master",
			dateOfBirth: "1975-02-03",
			expected:    cmd.Master,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.category, func(t *testing.T) {
			// given
			pattern := "2006-01-02"
			d, err := time.Parse(pattern, tc.dateOfBirth)
			require.NoError(t, err)
			// when
			actual := cmd.GetAgeCategory(d)
			// then
			assert.Equal(t, tc.expected, actual)
		})
	}
}
