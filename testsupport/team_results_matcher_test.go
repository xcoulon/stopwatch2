package testsupport_test

import (
	"fmt"

	"github.com/vatriathlon/stopwatch2/cmd"
	"github.com/vatriathlon/stopwatch2/testsupport"

	"github.com/google/go-cmp/cmp"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("document fragments matcher", func() {

	// given
	expected := []cmd.Team{
		{
			Name:        "Team 1",
			Gender:      "F",
			AgeCategory: "Master",
			BibNumber:   1,
			Members: []cmd.TeamMember{
				{
					FirstName: "Firstname1.1",
					LastName:  "Lastname1.1",
					Category:  "Master",
					Gender:    "F",
					Club:      "VILLENEUVE D'ASCQ TRIATHLON",
				},
				{
					FirstName: "Firstname1.2",
					LastName:  "Lastname1.2",
					Category:  "Master",
					Gender:    "F",
					Club:      "",
				},
			},
		},
	}
	matcher := testsupport.MatchTeams(expected)

	It("should match", func() {
		// given
		actual := []cmd.Team{
			{
				Name:        "Team 1",
				Gender:      "F",
				AgeCategory: "Master",
				BibNumber:   1,
				Members: []cmd.TeamMember{
					{
						FirstName: "Firstname1.1",
						LastName:  "Lastname1.1",
						Category:  "Master",
						Gender:    "F",
						Club:      "VILLENEUVE D'ASCQ TRIATHLON",
					},
					{
						FirstName: "Firstname1.2",
						LastName:  "Lastname1.2",
						Category:  "Master",
						Gender:    "F",
						Club:      "",
					},
				},
			},
		}
		// when
		result, err := matcher.Match(actual)
		// then
		Expect(err).ToNot(HaveOccurred())
		Expect(result).To(BeTrue())
	})

	It("should not match", func() {
		// given
		actual := []cmd.Team{
			{
				Name:        "Team 2",
				Gender:      "F",
				AgeCategory: "Master",
				BibNumber:   2,
			},
		}
		// when
		result, err := matcher.Match(actual)
		// then
		Expect(err).ToNot(HaveOccurred())
		Expect(result).To(BeFalse())
		diffs := cmp.Diff(expected, actual)
		Expect(matcher.FailureMessage(actual)).To(Equal(fmt.Sprintf("expected teams to match:\n%s", diffs)))
		Expect(matcher.NegatedFailureMessage(actual)).To(Equal(fmt.Sprintf("expected teams not to match:\n%s", diffs)))
	})

	It("should return error when invalid type is input", func() {
		// when
		result, err := matcher.Match(1)
		// then
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(Equal("MatchTeams matcher expects a '[]cmd.Team' (actual: int)"))
		Expect(result).To(BeFalse())
	})

})
