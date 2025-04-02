package main

import (
	"fmt"
	"slices"
	"time"
)

type dateValue struct {
	time.Time
}

func (d *dateValue) String() string {
	return d.Time.Format(time.DateOnly)
}

func (d *dateValue) Set(value string) error {
	t, err := time.Parse(time.DateOnly, value)
	if err != nil {
		return err
	}
	d.Time = t
	return nil
}

type inputFormat string

func (i *inputFormat) String() string {
	return string(*i)
}

func (i *inputFormat) Set(value string) error {
	valid := []string{"bitcoin", "checking", "card"}
	if slices.Contains(valid, value) {
		*i = inputFormat(value)
		return nil
	}
	return fmt.Errorf("invalid input format: %q, must be one of: %q", value, valid)
}

type outputFormat string

func (o *outputFormat) String() string {
	return string(*o)
}

func (o *outputFormat) Set(value string) error {
	valid := []string{"ynab", "lunchmoney", "coinledger", "cointracker", "koinly"}
	if slices.Contains(valid, value) {
		*o = outputFormat(value)
		return nil
	}
	return fmt.Errorf("invalid output format: %q, must be one of: %q", value, valid)
}
