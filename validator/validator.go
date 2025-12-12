package validator

import (
	"regexp"
	"strings"
	"time"
)

var datetimePattern = "2006-01-02T15:04:05"

func isValidAtExpression(pattern string) bool {
	if _, err := time.Parse(datetimePattern, pattern); err != nil {
		return false
	}

	return true
}

const (
	asterisk     = "\\*"
	questionMark = "\\?"
	minutes      = "[0-5]?[0-9]"                                                      // [0]0-59
	hours        = "([0-1]?[0-9]|2[0-3])"                                             // [0]0-23
	dayOfMonth   = "([0-2]?[0-9]|3[0-1])"                                             // [0]1-31
	month        = "(0?[1-9]|1[0-2]|JAN|FEB|MAR|APR|MAY|JUN|JUL|AUG|SEP|OCT|NOV|DEC)" // [0]1-12 or JAN-DEC
	dayOfWeek    = "([1-7]|SUN|MON|TUE|WED|THU|FRI|SAT)"                              // 1-7 or SUN-SAT
	year         = "(19[7-9][0-9]|2[0-1][0-9]{2})"                                    // 1970-2199
)

func rangeOf(pattern string) string {
	return pattern + "(-" + pattern + ")?"
}

func incrementOf(pattern string) string {
	return pattern + "(\\/[0-9]*[1-9][0-9]*)?"
}

func listOf(pattern string) string {
	return pattern + "(," + pattern + ")*"
}

func either(patterns ...string) string {
	return "(" + strings.Join(patterns, "|") + ")"
}

func exactly(pattern string) string {
	return "^" + pattern + "$"
}

var (
	minutesRegex = regexp.MustCompile(
		exactly(listOf(incrementOf(either(asterisk, rangeOf(minutes))))),
	)
	hoursRegex = regexp.MustCompile(
		exactly(listOf(incrementOf(either(asterisk, rangeOf(hours))))),
	)
	dayOfMonthRegex = regexp.MustCompile(
		exactly(
			either(
				listOf(either(incrementOf(either(asterisk, rangeOf(dayOfMonth))), dayOfMonth+"W")),
				questionMark,
				"L",
			),
		),
	)
	monthRegex = regexp.MustCompile(
		exactly(listOf(incrementOf(either(asterisk, rangeOf(month))))),
	)
	dayOfWeekRegex = regexp.MustCompile(
		exactly(
			either(
				listOf(either(asterisk, rangeOf(dayOfWeek))),
				dayOfWeek+"#[1-5]",
				questionMark,
				dayOfWeek+"?L",
			),
		),
	)
	yearRegex = regexp.MustCompile(
		exactly(listOf(incrementOf(either(asterisk, rangeOf(year))))),
	)
)

func isValidCronExpression(pattern string) bool {
	elements := strings.Split(pattern, " ")
	if len(elements) != 6 {
		return false
	}

	minutes, hours, dayOfMonth, month, dayOfWeek, year := elements[0], elements[1], elements[2], elements[3], elements[4], elements[5]

	if (dayOfMonth == "?" && dayOfWeek == "?") || (dayOfMonth != "?" && dayOfWeek != "?") {
		return false
	}

	return minutesRegex.MatchString(minutes) &&
		hoursRegex.MatchString(hours) &&
		dayOfMonthRegex.MatchString(dayOfMonth) &&
		monthRegex.MatchString(month) &&
		dayOfWeekRegex.MatchString(dayOfWeek) &&
		yearRegex.MatchString(year)
}

var (
	valueRegex = regexp.MustCompile(exactly("[1-9][0-9]*"))
	unitRegex  = regexp.MustCompile(exactly(either("minutes?", "hours?", "days?")))
)

func isValidRateExpression(pattern string) bool {
	elements := strings.Split(pattern, " ")
	if len(elements) != 2 {
		return false
	}

	value, unit := elements[0], elements[1]

	return valueRegex.MatchString(value) && unitRegex.MatchString(unit)
}

var exprRegex = regexp.MustCompile(`^([^(]*)\((.*)\)$`)

func Validate(expr string) bool {
	matches := exprRegex.FindStringSubmatch(expr)
	if len(matches) != 3 {
		return false
	}

	_, typ, pattern := matches[0], matches[1], matches[2]

	switch typ {
	case "at":
		return isValidAtExpression(pattern)
	case "cron":
		return isValidCronExpression(pattern)
	case "rate":
		return isValidRateExpression(pattern)
	default:
		return false
	}
}
