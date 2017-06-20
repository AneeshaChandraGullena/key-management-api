//Package logging .
// © Copyright 2016 IBM Corp. Licensed Materials – Property of IBM.
package logging

import (
	"os"
	"sync"

	"github.com/go-kit/kit/log"
)

var globalLogger log.Logger

type serializedLogger struct {
	mtx sync.Mutex
	log.Logger
}

func init() {
	globalLogger = log.NewLogfmtLogger(os.Stderr)
	globalLogger = &serializedLogger{Logger: globalLogger}
	globalLogger = log.NewContext(globalLogger).With("ts", log.DefaultTimestampUTC)
}

func (logger *serializedLogger) Log(keyvals ...interface{}) error {
	logger.mtx.Lock()
	defer logger.mtx.Unlock()
	return logger.Logger.Log(keyvals...)
}

// GlobalLogger returns a global logger if one does not exists, else return globalLogger
func GlobalLogger() log.Logger {
	return globalLogger
}
