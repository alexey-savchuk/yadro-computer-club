package system

import (
	"fmt"
	"strconv"
	"strings"
	"time"
	"unicode"
)

type Event struct {
	Time time.Time
	ID   int
	Body []any
}

const (
	InputEventEnter = iota + 1
	InputEventTakeTable
	InputEventWait
	InputEventLeave
)

const (
	OutputEventLeave = iota + 11
	OutputEventTakeTable
	OutputEventError
)

func ParseTime(s string) (time.Time, error) {
	s = strings.TrimSpace(s)

	parts := strings.Split(s, ":")
	if len(parts) != 2 {
		return time.Time{}, fmt.Errorf("invalid time format %s", s)
	}

	hours := parts[0]
	minutes := parts[1]

	if len(hours) != 2 || len(minutes) != 2 {
		return time.Time{}, fmt.Errorf("invalid time format %s", s)
	}

	for _, r := range hours + minutes {
		if !unicode.IsDigit(r) {
			return time.Time{}, fmt.Errorf("invalid time format %s", s)
		}
	}

	hoursInt, err := strconv.Atoi(hours)
	if err != nil {
		panic(fmt.Sprintf(
			"unexpected error while converting hours %s: %s",
			hours, err,
		))
	}

	minutesInt, err := strconv.Atoi(minutes)
	if err != nil {
		panic(fmt.Sprintf(
			"unexpected error while converting minutes %s: %s",
			minutes, err,
		))
	}

	if hoursInt > 23 {
		return time.Time{}, fmt.Errorf("invalid time format %s", s)
	}

	if minutesInt > 59 {
		return time.Time{}, fmt.Errorf("invalid time format %s", s)
	}

	time := time.Date(0, 0, 0, hoursInt, minutesInt, 0, 0, time.UTC)
	return time, nil
}

func ParsePositiveInt(s string) (int, error) {
	s = strings.TrimSpace(s)

	if len(s) == 0 {
		return 0, fmt.Errorf("empty number string")
	}

	isFirst := true
	for _, char := range s {
		if !unicode.IsDigit(char) {
			return 0, fmt.Errorf("invalid number %q format, %c is not a number", s, char)
		}
		if isFirst && char == '0' {
			return 0, fmt.Errorf("invalid number %q format, starts with zero", s)
		}
		isFirst = false
	}
	number, err := strconv.Atoi(s)
	if err != nil {
		panic("unexpected error while parsing number " + s)
	}
	return number, nil
}

func parseID(s string) (int, error) {
	s = strings.TrimSpace(s)

	id, err := ParsePositiveInt(s)
	if err != nil {
		return 0, fmt.Errorf("invalid id %q format: %w", s, err)
	}
	return id, nil
}

const nameSymbols = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_-"

func parseName(s string) (string, error) {
	s = strings.TrimSpace(s)

	if len(s) == 0 {
		return "", fmt.Errorf("empty name string")
	}

	for _, char := range s {
		if !strings.ContainsRune(nameSymbols, char) {
			return "", fmt.Errorf(
				"invalid name %q format, %c is not an allowed symbol",
				s, char,
			)
		}
	}
	return s, nil
}

func parseTableNum(s string) (int, error) {
	s = strings.TrimSpace(s)

	tableNum, err := ParsePositiveInt(s)
	if err != nil {
		return 0, fmt.Errorf("invalid table number %q format: %w", s, err)
	}
	return tableNum, nil
}

func ParseEvent(s string) (Event, error) {
	s = strings.TrimSpace(s)

	parts := strings.Split(s, " ")
	if len(parts) < 2 {
		return Event{}, fmt.Errorf(
			"invalid event %q format, expected at least <time> and <id>", s)
	}

	time, err := ParseTime(parts[0])
	if err != nil {
		return Event{}, fmt.Errorf("invalid event %q time format: %w", s, err)
	}

	id, err := parseID(parts[1])
	if err != nil {
		return Event{}, fmt.Errorf("invalid event %q id format: %w", s, err)
	}

	body := make([]any, 0)

	switch id {
	case InputEventTakeTable, OutputEventTakeTable:
		if len(parts) != 4 {
			return Event{}, fmt.Errorf(
				"invalid event %q format, event #%d requires <name> <table number> body",
				s, id,
			)
		}
		name, err := parseName(parts[2])
		if err != nil {
			return Event{}, fmt.Errorf("invalid event %q name format: %w", s, err)
		}
		tableNum, err := parseTableNum(parts[3])
		if err != nil {
			return Event{}, fmt.Errorf("invalid event %q table number format: %w", s, err)
		}
		body = append(body, name)
		body = append(body, tableNum)
	case OutputEventError:
		if len(parts) != 3 {
			return Event{}, fmt.Errorf(
				"invalid event %q format: event #%d requires <error> body",
				s, id,
			)
		}
		if len(parts[2]) == 0 {
			return Event{}, fmt.Errorf("empty event %q <error>", s)
		}
		body = append(body, parts[2])
	case InputEventEnter, InputEventWait, InputEventLeave, OutputEventLeave:
		if len(parts) != 3 {
			return Event{}, fmt.Errorf(
				"invalid event %q format, event #%d requires <name> body",
				s, id,
			)
		}
		name, err := parseName(parts[2])
		if err != nil {
			return Event{}, fmt.Errorf("invalid event %q name format: %w", s, err)
		}
		body = append(body, name)
	default:
		return Event{}, fmt.Errorf("unknown event %q, unknown event id=%d", s, id)
	}

	return Event{
		Time: time,
		ID:   id,
		Body: body,
	}, nil
}

func (e Event) String() string {

	builder := strings.Builder{}

	builder.WriteString(e.Time.Format("15:04"))
	builder.WriteString(" ")
	builder.WriteString(strconv.Itoa(e.ID))

	for _, v := range e.Body {
		builder.WriteString(" ")
		builder.WriteString(fmt.Sprint(v))
	}

	return builder.String()
}
