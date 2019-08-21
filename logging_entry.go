package glbr

import (
	"context"
	"fmt"
	"io"
	"time"

	"cloud.google.com/go/logging"
)

func push(c context.Context, entry logging.Entry) {
	if logger, ok := getLogger(c); ok {
		if logger == nil {
			panic("logger is nil, call initilize function 'NewLogging'")
		}
		// logging.client.errc is closed in the logging.Close function,
		// it will panic if called after Close function.
		logger.Log(entry)
	} else {
		fmt.Println("logger not found")
	}
	if w, ok := getIOWriter(c); ok {
		pl, _ := entry.Payload.(string)
		tm := entry.Timestamp.Format("2006/01/02 03:04:05")
		io.WriteString(w, fmt.Sprintf("%s %s: %s\n", tm, entry.Severity, pl))
	}
}

// sendEntry ログを送信する
func sendEntry(c context.Context, severity logging.Severity, format string, value ...interface{}) {
	if maxSeverity, ok := getSeverity(c); ok {
		if *maxSeverity < severity {
			*maxSeverity = severity
		}
	}
	traceID, ok := getTraceID(c)
	if !ok {
		traceID = new(string)
		*traceID = newTraceID()
	}
	push(c, logging.Entry{
		Payload:   fmt.Sprintf(format, value...),
		Severity:  severity,
		Trace:     *traceID,
		Timestamp: time.Now(),
	})
}

// CustomSeverityf 0 < Debugf(100) < ... < Emergencyf(700)
func CustomSeverityf(c context.Context, severity int, format string, value ...interface{}) {
	sendEntry(c, logging.Severity(severity), format, value...)
}

// Debugf 0 < Debugf
func Debugf(c context.Context, format string, value ...interface{}) {
	sendEntry(c, logging.Debug, format, value...)
}

// Infof Debugf < Infof
func Infof(c context.Context, format string, value ...interface{}) {
	sendEntry(c, logging.Info, format, value...)
}

// Noticef Infof < Noticef
func Noticef(c context.Context, format string, value ...interface{}) {
	sendEntry(c, logging.Notice, format, value...)
}

// Warningf Noticef < Warningf
func Warningf(c context.Context, format string, value ...interface{}) {
	sendEntry(c, logging.Warning, format, value...)
}

// Errorf Warningf < Errorf
func Errorf(c context.Context, format string, value ...interface{}) {
	sendEntry(c, logging.Error, format, value...)
}

// Criticalf Errorf < Criticalf
func Criticalf(c context.Context, format string, value ...interface{}) {
	sendEntry(c, logging.Critical, format, value...)
}

// Alertf Criticalf < Alertf
func Alertf(c context.Context, format string, value ...interface{}) {
	sendEntry(c, logging.Alert, format, value...)
}

// Emergencyf Alertf < Emergencyf
func Emergencyf(c context.Context, format string, value ...interface{}) {
	sendEntry(c, logging.Emergency, format, value...)
}
