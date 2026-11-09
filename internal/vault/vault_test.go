package vault_test

import (
	"testing"

	"github.com/jonathantrevino/resumes-cli/internal/vault"
)

func TestSlug(t *testing.T) {
	testCases := []struct {
		input   string
		output  string
		wantErr bool
	}{
		{
			input:  "   Stripe, Inc.   ",
			output: "stripe-inc",
		},
		{
			input:  "Backend Engineer II",
			output: "backend-engineer-ii",
		},
		{
			input:  "stripe",
			output: "stripe",
		},
		{
			input:  "Sr. Software Engineer (Remote)",
			output: "sr-software-engineer-remote",
		},
		{
			input:  "AT&T",
			output: "at-t",
		},
		{
			input:  "Engineer   ---   Platform",
			output: "engineer-platform",
		},
		{
			input:   "!!!",
			output:  "", // expect error
			wantErr: true,
		},
		{
			input:   "",
			output:  "", // expect error
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			result, err := vault.Slug(tc.input)
			if tc.wantErr {
				if err == nil {
					t.Errorf("Slug(%q) = %q; expected an error", tc.input, result)
				}
				return
			}
			if err != nil {
				t.Fatalf("Slug(%q) unexpected error: %v", tc.input, err)
			}
			if result != tc.output {
				t.Errorf("Slug(%q) = %q; expected %q", tc.input, result, tc.output)
			}
		})
	}
}
