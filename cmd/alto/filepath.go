package main

import (
	"errors"
	"path/filepath"
	"regexp"
	"runtime"
)

type Filepath struct {
	val *string
}

func (f Filepath) String() string {
	if f.val != nil {
		return *f.val
	}
	return ""
}

func (f Filepath) Set(value string) error {
	if v, err := filepath.Abs(filepath.Clean(value)); err != nil {
		return err
	} else {
		if value == "" {
			return errors.New("this flag requires a filepath")
		}
		*f.val = v
	}
	return nil
}

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
