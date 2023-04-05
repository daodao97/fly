package fly

import (
	"time"

	_logger "github.com/daodao97/gokit/logger"
)

var logger = _logger.Default()

func dbLog(prefix string, start time.Time, err *error, kv *[]interface{}) {
	tc := time.Since(start)
	_log := []interface{}{
		prefix,
		"ums:", tc.Milliseconds(),
	}
	_log = append(_log, *kv...)
	if *err != nil {
		_log = append(_log, "error:", *err)
		logger.Log(_logger.LevelError, "", _log...)
		return
	}
	logger.Log(_logger.LevelDebug, "", _log...)
}
