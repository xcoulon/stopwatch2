package cmd

import (
	"time"

	"gopkg.in/yaml.v3"
)

type Team struct {
	Name        string       `yaml:"name"`
	Gender      string       `yaml:"gender"`
	AgeCategory string       `yaml:"ageCategory"`
	BibNumber   int          `yaml:"bibNumber"`
	Members     []TeamMember `yaml:"members"`
}

// TeamMember a member of a team
type TeamMember struct {
	FirstName   string      `yaml:"firstName"`
	LastName    string      `yaml:"lastName"`
	DateOfBirth ISO8601Date `yaml:"dateOfBirth"`
	Category    string      `yaml:"category"`
	Gender      string      `yaml:"gender"`
	Club        string      `yaml:"club"`
}

type ISO8601Date time.Time

func (d *ISO8601Date) UnmarshalYAML(value *yaml.Node) error {
	var v string
	if err := value.Decode(&v); err != nil {
		return err
	}
	t, err := time.Parse("2006-01-02", string(v))
	if err != nil {
		return err
	}
	*d = ISO8601Date(t)
	return nil
}

func (d ISO8601Date) MarshalYAML() (interface{}, error) {
	return time.Time(d).Format("2006-01-02"), nil
}
