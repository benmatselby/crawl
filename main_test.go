package main

import (
	"errors"
	"testing"
)

func TestRunCanParseFlags(t *testing.T) {
	tt := []struct {
		name     string
		args     []string
		expected error
	}{
		{name: "no arguments passed in", args: []string{}, expected: errors.New("no URL specified")},
		{name: "bad url", args: []string{"flim flam"}, expected: errors.New("invalid URL")},
		{name: "valid url", args: []string{"https://bbc.co.uk"}, expected: nil},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			err := Run(tc.args)
			if tc.expected == nil && err != nil {
				t.Fatalf("did not expect error, got %v", err)
			}

			if tc.expected != nil {
				if err.Error() != tc.expected.Error() {
					t.Fatalf("expected %v, got %v", tc.expected, err)
				}
			}
		})
	}
}
