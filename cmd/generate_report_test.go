package cmd_test

import (
	"io/ioutil"
	"os"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/vatriathlon/stopwatch2/cmd"
)

var _ = Describe("generate reports", func() {

	It("should generate scratch report", func() {
		// given
		teamFilename := "../test/teams-xs.yaml"
		raceFilename := "../test/race-xs.yaml"
		outputFile, err := ioutil.TempFile(os.TempDir(), "output-*.adoc")
		Expect(err).NotTo(HaveOccurred())

		// when
		err = cmd.GenerateReport("Race XS", teamFilename, raceFilename, outputFile.Name())
		Expect(err).NotTo(HaveOccurred())
		result, err := os.ReadFile(outputFile.Name())
		Expect(err).NotTo(HaveOccurred())
		Expect(string(result)).To(Equal(`= Classement Race XS

== Classement scratch

[cols="2,5,5,5,8,8,4"]
|===
|# |Dossard |Equipe |Catégorie |Coureurs |Club |Temps Total

|1 |3 |Team 3 |M/F |Lamy - Poulain |VILLENEUVE D ASCQ TRIATHLON |50m0s 
|2 |2 |Team 2 |M/F |Aubert - Le Gall | |51m30s 
|3 |1 |Team 1 |M/M |Bonnin - Georges |TOBESPORT |52m55s 
|===
`))
	})

})
