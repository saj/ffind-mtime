package main

import (
	"fmt"
	"regexp"
	"strconv"
	"time"
)

var durationRE = regexp.MustCompile("^([0-9]+)(y|w|d|h|m|s)?$")

func ParseDuration(s string) (time.Duration, error) {
	matches := durationRE.FindStringSubmatch(s)
	if len(matches) != 3 {
		return 0, fmt.Errorf("invalid duration string: %q", s)
	}
	var (
		n, _ = strconv.Atoi(matches[1])
		d    = time.Duration(n) * time.Second
	)
	switch unit := matches[2]; unit {
	case "y":
		d *= 60 * 60 * 24 * 365
	case "w":
		d *= 60 * 60 * 24 * 7
	case "d":
		d *= 60 * 60 * 24
	case "h":
		d *= 60 * 60
	case "m":
		d *= 60
	case "s":
	case "":
	default:
		return 0, fmt.Errorf("invalid time unit in duration string: %q", unit)
	}
	return d, nil
}
