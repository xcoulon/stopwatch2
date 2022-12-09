package cmd

import (
	"math"
	"time"

	"github.com/sirupsen/logrus"
)

const (
	// MiniPoussin 2015 à 2016
	MiniPoussin = "Mini-poussin"
	// Poussin 		2013 à 2014
	Poussin = "Poussin"
	// Pupille 		2011 à 2012
	Pupille = "Pupille"
	// Benjamin 	2009 à 2010
	Benjamin = "Benjamin"
	// Minime 		2007 à 2008
	Minime = "Minime"
	// Cadet 		2005 à 2006
	Cadet = "Cadet"
	// Junior 		2003 à 2004
	Junior = "Junior"
	// Senior 	 	1983 à 2002
	Senior = "Senior"
	// Master 		1955 à 1982
	Master = "Master"
)

// GetAgeCategory gets the age category associated with the given date of birth
func GetAgeCategory(dateOfBirth time.Time) string {
	yearOfBirth := dateOfBirth.Year()
	logrus.WithField("year_of_birth", yearOfBirth).Debug("computing age category")
	switch {
	case yearOfBirth == 2015 || yearOfBirth == 2016:
		return MiniPoussin
	case yearOfBirth == 2013 || yearOfBirth == 2014:
		return Poussin
	case yearOfBirth == 2011 || yearOfBirth == 2012:
		return Pupille
	case yearOfBirth == 2009 || yearOfBirth == 2010:
		return Benjamin
	case yearOfBirth == 2007 || yearOfBirth == 2008:
		return Minime
	case yearOfBirth == 2005 || yearOfBirth == 2006:
		return Cadet
	case yearOfBirth == 2003 || yearOfBirth == 2004:
		return Junior
	case yearOfBirth >= 1983 && yearOfBirth <= 2002:
		return Senior
	default:
		return Master
	}

}

var ageCategories = map[string]int{
	MiniPoussin: 1,
	Poussin:     2,
	Pupille:     3,
	Benjamin:    4,
	Minime:      5,
	Cadet:       6,
	Junior:      7,
	Master:      8,
	Senior:      9, // so that max(master, senior) -> senior
}

// GetTeamAgeCategory computes the age category for the team
func GetTeamAgeCategory(ageCategory1, ageCategory2 string) string {
	cat1 := ageCategories[ageCategory1]
	cat2 := ageCategories[ageCategory2]
	// assign to senior if 1 veteran + 1 under senior
	if (ageCategory1 == Master && cat2 <= ageCategories[Junior]) || (ageCategory2 == Master && cat1 <= ageCategories[Junior]) {
		return Senior
	}
	teamAgeCategoryValue := math.Max(float64(cat1), float64(cat2))
	logrus.WithField("team_age_category_value", teamAgeCategoryValue).Debugf("computing team age category...")
	//

	for k, v := range ageCategories {
		if float64(v) == teamAgeCategoryValue {
			return k
		}
	}
	return ""

}
