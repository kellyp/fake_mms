package ledger

import (
	"strings"
)

type Level int8

const (
	FatalLevel Level = iota
	ErrorLevel
	WarnLevel
	InfoLevel
	DebugLevel
)

var (
	levelStrings = []string{
		"FATAL",
		"ERROR",
		"WARN",
		"INFO",
		"DEBUG",
	}
)

func (level Level) String() string {
	return levelStrings[level]
}

func ParseLevel(lvl string) Level {
	switch strings.ToUpper(lvl) {
	case levelStrings[0]:
		return FatalLevel
	case levelStrings[1]:
		return ErrorLevel
	case levelStrings[2]:
		return WarnLevel
	case levelStrings[3]:
		return InfoLevel
	default:
		return DebugLevel
	}
}
