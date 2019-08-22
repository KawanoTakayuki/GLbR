package glbr

import (
	"context"
	"io"

	"cloud.google.com/go/logging"
)

var (
	loggerKey            = "loggerClient"       // logger key
	requestKey           = "request"            // request key
	severityKey          = "severity"           // severity key
	traceIDKey           = "trace-id"           // traceid key
	logIDKey             = "log-id"             // logid key
	iowriteKey           = "io-write"           // iowrite key
	monitoredResourceKey = "monitored-resource" // monitoredresource key
)

// logger setter
func setLogger(c context.Context, logger *logging.Logger) context.Context {
	return context.WithValue(c, &loggerKey, logger)
}

// logger getter
func getLogger(c context.Context) (*logging.Logger, bool) {
	logger, ok := c.Value(&loggerKey).(*logging.Logger)
	return logger, ok
}

// parent severity setter
func setSeverity(c context.Context, severity *logging.Severity) context.Context {
	return context.WithValue(c, &severityKey, severity)
}

// parent severity getter
func getSeverity(c context.Context) (*logging.Severity, bool) {
	severity, ok := c.Value(&severityKey).(*logging.Severity)
	return severity, ok
}

// traceid setter
func setTraceID(c context.Context, traceID *string) context.Context {
	return context.WithValue(c, &traceIDKey, traceID)
}

// traceid getter
func getTraceID(c context.Context) (*string, bool) {
	traceID, ok := c.Value(&traceIDKey).(*string)
	return traceID, ok
}

// io.Writer setter
func setIOWriter(c context.Context, w io.Writer) context.Context {
	return context.WithValue(c, &iowriteKey, w)
}

// io.Writer getter
func getIOWriter(c context.Context) (io.Writer, bool) {
	w, ok := c.Value(&iowriteKey).(io.Writer)
	return w, ok
}
