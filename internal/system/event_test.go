package system_test

import (
	"reflect"
	"testing"
	"time"

	system "github.com/alexey-savchuk/computer-club/internal/system"
)

var testCases = []struct {
	input   string
	want    system.Event
	isError bool
}{
	{
		"24:03 1 j0hn_d0e",
		system.Event{},
		true,
	},
	{
		"15:60 1 j0hn_d0e",
		system.Event{},
		true,
	},
	{
		"8:24 1 j0hn_d0e",
		system.Event{},
		true,
	},
	{
		"12:8 1 j0hn_d0e",
		system.Event{},
		true,
	},
	{
		"15:03 01 j0hn_d0e",
		system.Event{},
		true,
	},
	{
		"15:03 -1 j0hn_d0e",
		system.Event{},
		true,
	},
	{
		"15:03 1 j0hn_d0e$",
		system.Event{},
		true,
	},
	{
		"15:03 1 j0hn_d0e",
		system.Event{
			Time: time.Date(0, 0, 0, 15, 3, 0, 0, time.UTC),
			ID:   1,
			Body: []any{"j0hn_d0e"},
		},
		false,
	},
	{
		"15:03 1",
		system.Event{},
		true,
	},
	{
		"15:03 1 j0hn_d0e something",
		system.Event{},
		true,
	},
	{
		"15:03 2 j0hn_d0e 12",
		system.Event{
			Time: time.Date(0, 0, 0, 15, 3, 0, 0, time.UTC),
			ID:   2,
			Body: []any{"j0hn_d0e", 12},
		},
		false,
	},
	{
		"15:03 2",
		system.Event{},
		true,
	},
	{
		"15:03 2 j0hn_d0e",
		system.Event{},
		true,
	},
	{
		"15:03 2 j0hn_d0e 12 something",
		system.Event{},
		true,
	},
	{
		"15:03 13 error_happened",
		system.Event{
			Time: time.Date(0, 0, 0, 15, 3, 0, 0, time.UTC),
			ID:   13,
			Body: []any{"error_happened"},
		},
		false,
	},
	{
		"15:03 13",
		system.Event{},
		true,
	},
	{
		"15:03 13 error_happened something",
		system.Event{},
		true,
	},
}

func testParse(t *testing.T, input string, want system.Event, isError bool) {
	t.Logf("parse %q", input)

	got, err := system.ParseEvent(input)
	if isError && err == nil {
		t.Fatalf("expected error while parse %q", input)
	}
	if !isError && err != nil {
		t.Fatalf("unexpected error while parse %q: %s", input, err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("expected %q, got %q", want, got)
	}
}

func TestParse(t *testing.T) {
	for _, tc := range testCases {
		testParse(t, tc.input, tc.want, tc.isError)
	}
}
