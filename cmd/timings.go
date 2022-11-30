package cmd

import (
	"fmt"
	"time"
)

// type RaceResults struct {
// 	StartTime Time          `yaml:"start"`
// 	Teams     []TeamScratch `yaml:"teams"`
// }

type Timings []Timing

type Timing map[string]interface{}

func (t Timing) Time() (time.Time, error) {
	for k := range t {
		return time.Parse(TimeFormat, k)
	}
	return time.Time{}, fmt.Errorf("unexpected timing: %v", t)
}

func (t Timing) IsStart() bool {
	for _, v := range t {
		if v, ok := v.(string); ok {
			return v == "start"
		}
	}
	return false
}

func (t Timing) BibNumber() (int, error) {

	for _, v := range t {
		if v, ok := v.(int); ok {
			return v, nil
		}
	}
	return -1, fmt.Errorf("unexpected timing: %v", t)
}

// type TeamScratch struct {
// 	BibNumber   int  `yaml:"bibNumber"`
// 	ScratchTime Time `yaml:"scratch"`
// }

// type Time time.Time

// func (d *Time) UnmarshalYAML(value *yaml.Node) error {
// 	var v string
// 	if err := value.Decode(&v); err != nil {
// 		return err
// 	}
// 	t, err := time.Parse("15:04:05", string(v))
// 	if err != nil {
// 		return err
// 	}
// 	*d = Time(t)
// 	return nil
// }

// func (d Time) MarshalYAML() (interface{}, error) {
// 	return time.Time(d).Format("15:04:05"), nil
// }

// func (t Time) Sub(s Time) time.Duration {
// 	return time.Time(t).Sub(time.Time(s))
// }
