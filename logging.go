package glbr

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"cloud.google.com/go/logging"
	"google.golang.org/api/option"
)

// http.ResponseWriter interface
type logResponse struct {
	body []byte
	code int
}

func (lr *logResponse) Header() http.Header {
	return make(http.Header)
}
func (lr *logResponse) Write(body []byte) (int, error) {
	lr.body = body
	return len(body), nil
}
func (lr *logResponse) WriteHeader(statusCode int) {
	lr.code = statusCode
}

var (
	client      *logging.Client // logClient
	severityKey = "severity"    // severity key
	traceIDKey  = "trace-id"    // traceid key
	logIDKey    = "log-id"      // logid key
)

// GroupingBy ログをリクエストでグループ化する
func GroupingBy(c context.Context, w http.ResponseWriter, r *http.Request, f func(c context.Context, w http.ResponseWriter, r *http.Request)) {
	s := time.Now()
	defaultSeverity := logging.Default
	logCtx := setTraceID(c, fmt.Sprintf("%d", rand.Uint64()))
	logCtx = setSeverity(logCtx, &defaultSeverity)

	res := &logResponse{code: http.StatusOK}
	f(logCtx, res, r)
	w.WriteHeader(res.code)
	if _, err := w.Write(res.body); err != nil {
		Criticalf(logCtx, err.Error())
	}

	maxSeverity, _ := getSeverity(logCtx)
	traceID, _ := getTraceID(logCtx)
	client.Logger("request_log").Log(logging.Entry{
		HTTPRequest: &logging.HTTPRequest{
			Status:       res.code,
			ResponseSize: int64(len(res.body)),
			Request:      r,
			Latency:      time.Now().Sub(s),
		},
		Trace:    traceID,
		Severity: *maxSeverity,
	})
}

// NewLogging 新しいLoggingContextを取得する
func NewLogging(c context.Context, parent string, opts ...option.ClientOption) (logctx context.Context, err error) {
	client, err = logging.NewClient(c, parent, opts...)
	rand.Seed(time.Now().Unix())
	logctx = setLogID(c, "message_log")
	logctx = setTraceID(logctx, fmt.Sprintf("%d", rand.Uint64()))
	return
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
func setTraceID(c context.Context, traceID string) context.Context {
	return context.WithValue(c, &traceIDKey, traceID)
}

// traceid getter
func getTraceID(c context.Context) (string, bool) {
	traceID, ok := c.Value(&traceIDKey).(string)
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

// sendEntry ログを送信する
func sendEntry(c context.Context, severity logging.Severity, format string, value ...interface{}) {
	if maxSeverity, ok := getSeverity(c); ok {
		if *maxSeverity < severity {
			*maxSeverity = severity
			setSeverity(c, maxSeverity)
		}
	}
	traceID, _ := getTraceID(c)
	logID, _ := getLogID(c)
	client.Logger(logID).Log(logging.Entry{
		Payload:  fmt.Sprintf(format, value...),
		Severity: severity,
		Trace:    traceID,
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
