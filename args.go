package main

import (
	"fmt"
	"slices"
	"strings"
	"time"
)

type dateValue struct {
	time.Time
}

func (d *dateValue) String() string {
	return d.Format(time.DateOnly)
}

func (d *dateValue) Set(value string) (err error) {
	d.Time, err = time.Parse(time.DateOnly, value)
	return err
}

type enumFlag struct {
	Allowed []string
	Value   string
}

// newEnumFlag give a list of allowed flag parameters, where the second argument is the default.
func newEnumFlag(allowed []string, d string) *enumFlag {
	return &enumFlag{
		Allowed: allowed,
		Value:   d,
	}
}

func (e enumFlag) String() string {
	return e.Value
}

func (e *enumFlag) Set(p string) error {
	if !slices.Contains(e.Allowed, p) {
		return fmt.Errorf("must be one of: %s", strings.Join(e.Allowed, ", "))
	}
	e.Value = p
	return nil
}

func (e *enumFlag) Usage(t string) string {
	return fmt.Sprintf("%s (one of: %s)", t, strings.Join(e.Allowed, ", "))
}

func (e *enumFlag) Type() string {
	return "string"
}
