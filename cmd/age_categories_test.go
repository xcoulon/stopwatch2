package cmd_test

import (
	"time"

	"github.com/vatriathlon/stopwatch2/cmd"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = DescribeTable("age categories",

	func(dateOfBirth string, expected string) {
		// given
		pattern := "2006-01-02"
		d, err := time.Parse(pattern, dateOfBirth)
		Expect(err).NotTo(HaveOccurred())
		// when
		result := cmd.GetAgeCategory(d)
		// then
		Expect(result).To(Equal(expected))
	},
	Entry("mini poussin", "2015-02-03", cmd.MiniPoussin),
	Entry("poussin", "2013-02-03", cmd.Poussin),
	Entry("pupille", "2012-02-03", cmd.Pupille),
	Entry("benjamin", "2010-02-03", cmd.Benjamin),
	Entry("cadet", "2005-02-03", cmd.Cadet),
	Entry("junior", "2003-02-03", cmd.Junior),
	Entry("senior", "1984-02-03", cmd.Senior),
	Entry("junior", "1975-02-03", cmd.Master),
)

var _ = DescribeTable("team age categories",
	func(category1, category2 string, expected string) {
		result := cmd.GetTeamAgeCategory(category1, category2)
		// then
		Expect(result).To(Equal(expected))
	},
	Entry("mini poussin/mini poussin", cmd.MiniPoussin, cmd.MiniPoussin, cmd.MiniPoussin),
	Entry("mini poussin/poussin", cmd.MiniPoussin, cmd.Poussin, cmd.Poussin),
	Entry("poussin/poussin", cmd.Poussin, cmd.Poussin, cmd.Poussin),
	Entry("poussin/pupille", cmd.Poussin, cmd.Pupille, cmd.Pupille),
	Entry("benjamin/minime", cmd.Benjamin, cmd.Minime, cmd.Minime),
	Entry("senior/senior", cmd.Senior, cmd.Senior, cmd.Senior),
	Entry("master/senior", cmd.Master, cmd.Senior, cmd.Senior),
	Entry("master/master", cmd.Master, cmd.Master, cmd.Master),
)
