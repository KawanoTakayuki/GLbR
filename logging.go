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

// Service loggingService
type Service struct {
	ctx       context.Context
	client    *logging.Client
	option    []logging.LoggerOption
	logID     string
	projectID string
}

// NewLogging 新しいLoggingServiceを取得する
func NewLogging(c context.Context, projectID, logID string, opts ...option.ClientOption) (service Service, err error) {
	if logID == "" || 512 <= len(logID) {
		return Service{}, fmt.Errorf("logID empty or more than 512 char")
	}
	client, err := logging.NewClient(c, projectID, opts...)
	rand.Seed(time.Now().UnixNano())
	logctx := setLogID(c, logID)
	logctx = setTraceID(logctx, fmt.Sprintf("%d", rand.Uint64()))
	service = Service{
		ctx:       logctx,
		client:    client,
		option:    make([]logging.LoggerOption, 0),
		logID:     logID,
		projectID: projectID,
	}
	return
}

// Context log service context
func (s Service) Context() context.Context {
	return setLogger(s.ctx, s.client.Logger(s.logID, s.option...))
}

// Close serviceを閉じる
func (s Service) Close() (err error) {
	return s.client.Close()
}

// GroupingFunc グループ化される処理
// return code: httpStatusCode size: httpResponseSize
type GroupingFunc func(ctx context.Context) (code int, size int64)

// GroupingBy ログをリクエストでグループ化する
func (s Service) GroupingBy(r *http.Request, parentLogID string, f GroupingFunc) {
	if _, ok := getGroupKey(s.ctx); ok {
		// known group
		f(s.ctx)
		return
	}

	if r == nil {
		panic("empty to http.Request")
	}
	if parentLogID == "" {
		panic("empty to parentLogID")
	}
	if s.logID == parentLogID {
		panic("do not make parentLogID and the argument logID of 'NewLogging' functin identical")
	}

	severity := logging.Default
	traceID := fmt.Sprintf("%d", rand.Uint64())
	s.ctx = setSeverity(s.ctx, &severity)
	s.ctx = setTraceID(s.ctx, traceID)
	s.ctx = setGroupKey(s.ctx, "grouping")

	st := time.Now()
	code, size := f(s.Context())
	et := time.Now()

	if r.URL.String() == "" {
		r.URL.Path = "Empty_RequestUrl"
	}
	s.client.Logger(parentLogID, s.option...).Log(logging.Entry{
		HTTPRequest: &logging.HTTPRequest{
			Status:       code,
			ResponseSize: int64(size),
			Request:      r,
			Latency:      et.Sub(st),
		},
		Timestamp: et,
		Trace:     traceID,
		Severity:  severity,
		Resource:  getMonitoredResource(s.ctx),
	})
}
