package system

import (
	"fmt"
	"strings"
	"time"
)

type Preamble struct {
	NumTables int
	OpenTime  time.Time
	CloseTime time.Time
	Cost      int
}

func ParsePreamble(lines []string) (Preamble, error) {
	if len(lines) != 3 {
		return Preamble{}, fmt.Errorf("invalid number of lines: %d", len(lines))
	}

	numTables, err := ParsePositiveInt(lines[0])
	if err != nil {
		return Preamble{}, fmt.Errorf("failed to parse number of tables: %w", err)
	}

	parts := strings.Split(strings.TrimSpace(lines[1]), " ")
	if len(parts) != 2 {
		return Preamble{}, fmt.Errorf("invalid time format")
	}

	openTime, err := ParseTime(parts[0])
	if err != nil {
		return Preamble{}, fmt.Errorf("failed to parse open time: %w", err)
	}

	closeTime, err := ParseTime(parts[1])
	if err != nil {
		return Preamble{}, fmt.Errorf("failed to parse close time: %w", err)
	}

	cost, err := ParsePositiveInt(lines[2])
	if err != nil {
		return Preamble{}, fmt.Errorf("failed to parse cost: %w", err)
	}

	return Preamble{
		NumTables: numTables,
		OpenTime:  openTime,
		CloseTime: closeTime,
		Cost:      cost,
	}, nil
}
