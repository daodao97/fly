package ggm

import (
	"github.com/fatih/color"
	"log"
	"os"
)

var (
	red    = color.New(color.BgRed).Sprint
	yellow = color.New(color.BgHiYellow).Sprint
	green  = color.New(color.BgGreen).Sprint
	gray   = color.New(color.BgCyan).Sprint
)

type Level int

func (l Level) String() string {
	switch l {
	case LevelDebug:
		return gray(" DEBUG ")
	case LevelInfo:
		return green(" INFO ")
	case LevelWarn:
		return yellow(" WARNING ")
	case LevelErr:
		return red(" ERROR ")
	}
	return gray(" DEBUG ")
}

const LevelDebug = Level(0)
const LevelInfo = Level(1)
const LevelWarn = Level(2)
const LevelErr = Level(3)

type Logger interface {
	Log(level Level, keyValues ...any) error
}

var (
	logger     = newStdOutLogger()
	limitLevel = LevelDebug
)

func SetLogger(customLogger Logger, customLimitLevel Level) {
	logger = customLogger
	limitLevel = customLimitLevel
}

func newStdOutLogger() Logger {
	return &stdOutLogger{
		logger: log.New(os.Stdout, "GGM LOG: ", log.Lmicroseconds),
	}
}

type stdOutLogger struct {
	logger *log.Logger
}

func (s stdOutLogger) Log(level Level, keyValues ...any) error {
	if level < limitLevel {
		return nil
	}
	args := []any{level.String()}
	args = append(args, keyValues...)
	s.logger.Println(args...)
	return nil
}
