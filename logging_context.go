package glbr

import (
	"context"

	"cloud.google.com/go/logging"
)

var (
	loggerKey            = "loggerClient"       // logger key
	requestKey           = "request"            // request key
	severityKey          = "severity"           // severity key
	traceIDKey           = "trace-id"           // traceid key
	logIDKey             = "log-id"             // logid key
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

// logid setter
func setLogID(c context.Context, logID string) context.Context {
	return context.WithValue(c, &logIDKey, logID)
}

// logid getter
func getLogID(c context.Context) (string, bool) {
	logID, ok := c.Value(&logIDKey).(string)
	return logID, ok
}
