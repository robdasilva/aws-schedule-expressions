package validator_test

import (
	"testing"

	"github.com/robdasilva/aws-schedule-expressions/validator"
)

func TestValidate(t *testing.T) {
	t.Parallel()

	cases := map[bool][]string{
		true: {
			"at(2023-10-15T10:00:00)",
			"cron(0 12 * * ? *)",
			"cron(0 18 ? * MON-FRI *)",
			"cron(15 10 ? * 6L 2022)",
			"rate(1 days)",
			"rate(5 minutes)",
		},
		false: {
			"at 2023-12-15T10:00:00",      // Invalid notation
			"at(2023-13-15T10:00:00)",     // Invalid month
			"cron(0 18 ? * MON-FRIDAY *)", // Invalid day of week format
			"cron(15 10 ? * 8L 2022)",     // Invalid day of week
			"cron(60 12 * * ? *)",         // Invalid minute
			"every(5 minutes)",            // Invalid type
			"rate(-5 minutes)",            // Negative rate
			"rate(10 hour)",               // Singular unit
			"rate(0 days)",                // Zero rate
		},
	}

	for expected, expressions := range cases {
		for _, expr := range expressions {
			t.Run(expr, func(t *testing.T) {
				t.Parallel()

				actual := validator.Validate(expr)

				if actual != expected {
					t.Errorf("Validate(%q) = %v; want %v", expr, actual, expected)
				}
			})
		}
	}
}
