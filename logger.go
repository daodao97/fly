package fly

import (
	"encoding/json"
	"log"
	"os"
	"time"

	"github.com/fatih/color"
)

var (
	red    = color.New(color.BgRed).Sprint
	yellow = color.New(color.BgHiYellow).Sprint
	green  = color.New(color.BgGreen).Sprint
	gray   = color.New(color.BgCyan).Sprint
	prefix = color.New(color.FgGreen).Sprint
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
	Log(level Level, keyValues ...interface{}) error
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
		logger: log.New(os.Stdout, prefix("FLY LOG")+" ", log.Lmicroseconds),
	}
}

type stdOutLogger struct {
	logger *log.Logger
}

func jsonEncode(v interface{}) string {
	bt, _ := json.Marshal(v)
	return string(bt)
}

func (s stdOutLogger) Log(level Level, keyValues ...interface{}) error {
	if level < limitLevel {
		return nil
	}
	args := []interface{}{level.String()}
	args = append(args, keyValues...)

	for i, v := range args {
		switch t := v.(type) {
		case []interface{}:
			args[i] = jsonEncode(t)
		}
	}

	s.logger.Println(args...)
	return nil
}

func dbLog(prefix string, start time.Time, err *error, kv *[]interface{}) {
	tc := time.Since(start)
	_log := []interface{}{
		prefix,
		"ums:", tc.Milliseconds(),
	}
	_log = append(_log, *kv...)
	if *err != nil {
		_log = append(_log, "error:", *err)
		_ = logger.Log(LevelErr, _log...)
		return
	}
	_ = logger.Log(LevelDebug, _log...)
}
