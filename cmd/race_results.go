package cmd

import (
	"time"

	"gopkg.in/yaml.v3"
)

type RaceResults struct {
	StartTime Time          `yaml:"start"`
	Teams     []TeamScratch `yaml:"teams"`
}

type TeamScratch struct {
	BibNumber   int  `yaml:"bibNumber"`
	ScratchTime Time `yaml:"scratch"`
}

type Time time.Time

func (d *Time) UnmarshalYAML(value *yaml.Node) error {
	var v string
	if err := value.Decode(&v); err != nil {
		return err
	}
	t, err := time.Parse("15:04:05", string(v))
	if err != nil {
		return err
	}
	*d = Time(t)
	return nil
}

func (d Time) MarshalYAML() (interface{}, error) {
	return time.Time(d).Format("15:04:05"), nil
}

func (t Time) Sub(s Time) time.Duration {
	return time.Time(t).Sub(time.Time(s))
}
