package main

import (
	"regexp"
	"runtime"
)

var reservedKeywords *regexp.Regexp

// TODO: Find a library that can make this more flexible with varying filesystem
func clean(s string) string {
	if reservedKeywords == nil {
		if runtime.GOOS == "windows" {
			reservedKeywords = regexp.MustCompile(`[\pC"*/:<>?\\|]+`)
		} else {
			reservedKeywords = regexp.MustCompile(`[/\x{0}]+`)
		}
	}
	return reservedKeywords.ReplaceAllString(s, "-")
}
