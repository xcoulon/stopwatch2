package cmd

import (
	"fmt"
	"time"
)

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
