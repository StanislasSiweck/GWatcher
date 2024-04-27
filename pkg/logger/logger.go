package logger

import (
	"github.com/lmittmann/tint"
	"log/slog"
	"os"
)

const (
	LevelFatal = slog.Level(12)
)

type levelInfo struct {
	color string
	name  string
}

var LevelStruct = map[slog.Leveler]levelInfo{
	LevelFatal:      {color: "\033[0;31m", name: "FATAL"}, // Red
	slog.LevelInfo:  {color: "\033[0;32m", name: "INFO"},  // Green
	slog.LevelWarn:  {color: "\033[0;33m", name: "WARN"},  // Yellow
	slog.LevelError: {color: "\033[0;31m", name: "ERROR"}, // Red
}

const colorNone = "\033[0m"

func (l *levelInfo) String() string {
	return l.color + l.name + colorNone
}

var attribute = func(groups []string, a slog.Attr) slog.Attr {
	if a.Key == slog.LevelKey {
		level := a.Value.Any().(slog.Level)

		levelLabel := level.String()
		levelStruct, exists := LevelStruct[level]
		if exists {
			levelLabel = levelStruct.String()
		}

		a.Value = slog.StringValue(levelLabel)
	}
	return a
}

func New() {

	level := getLogLevel(slog.LevelInfo)

	if os.Getenv("ENV") != "production" {
		slog.SetDefault(slog.New(
			tint.NewHandler(os.Stdout, &tint.Options{
				Level:       level,
				TimeFormat:  "2006-01-02 15:04:05",
				ReplaceAttr: attribute,
			}),
		))
	} else {
		slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level:       level,
			ReplaceAttr: attribute,
		})))
	}
}

func getLogLevel(defaultLevel slog.Level) (level slog.Level) {
	switch os.Getenv("LOG_LEVEL") {
	case "DEBUG":
		level = slog.LevelDebug
	case "INFO":
		level = slog.LevelInfo
	case "WARN":
		level = slog.LevelWarn
	case "ERROR":
		level = slog.LevelError
	default:
		level = defaultLevel
	}
	return
}

func Fatal(msg string, args ...any) {
	slog.Log(nil, LevelFatal, msg, args...)
	os.Exit(1)
}
