package log

import (
	log "github.com/labstack/gommon/log"
)

type Logger struct {
	*log.Logger
}

var (
	global        = log.New("log")
	defaultHeader = `{"time":"${time_rfc3339}","level":"${level}","prefix":"${prefix}",` +
		`"file":"${long_file}","line":"${line}"}`
)

// Init function
func init() {
	log.SetHeader(defaultHeader)
	log.SetLevel(log.DEBUG)

	global.SetHeader(defaultHeader)
	global.SetLevel(log.DEBUG)
}

func Debugf(format string, values ...interface{}) {
	global.Debugf(format, values...)
}