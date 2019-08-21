package glbr

import (
	"context"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"time"

	"cloud.google.com/go/logging"
	"google.golang.org/api/option"
)

// Service loggingService
type Service struct {
	ctx    context.Context
	client *logging.Client
	option []logging.LoggerOption
	logID  string
}

// NewLogging 新しいLoggingServiceを取得する
func NewLogging(c context.Context, projectID, logID string, opts ...option.ClientOption) (service Service, err error) {
	if logID == "" || 512 <= len(logID) {
		return Service{}, fmt.Errorf("logID empty or more than 512 char")
	}
	client, err := logging.NewClient(c, projectID, opts...)
	logctx := setLogID(c, logID)
	service = Service{
		ctx:    logctx,
		client: client,
		option: make([]logging.LoggerOption, 0),
		logID:  logID,
	}
	return
}

// WithIOWriter write buffer after log output
func (s Service) WithIOWriter(w io.Writer) Service {
	s.ctx = setIOWriter(s.ctx, w)
	return s
}

// Context log service context
func (s Service) Context() context.Context {
	return setLogger(s.ctx, s.client.Logger(s.logID, s.option...))
}

// Close serviceを閉じる
func (s Service) Close() (err error) {
	return s.client.Close()
}

// NewTraceID 新しいTraceIDを返す
func newTraceID() string {
	rand.Seed(time.Now().UnixNano())
	return fmt.Sprintf("%d", rand.Uint64())
}

// http.ResponseWriter interface
type logResponse struct {
	body   []byte
	code   int
	origin http.ResponseWriter
}

func (lr *logResponse) Header() http.Header {
	return lr.origin.Header()
}
func (lr *logResponse) Write(body []byte) (int, error) {
	lr.body = body
	return lr.origin.Write(body)
}
func (lr *logResponse) WriteHeader(statusCode int) {
	lr.code = statusCode
	lr.origin.WriteHeader(statusCode)
}

// GroupingHandler グループ化される処理
type GroupingHandler func(http.Handler) http.Handler

// GroupedBy ログをリクエストでグループ化する
func (s Service) GroupedBy(parentLogID string) GroupingHandler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r == nil {
				panic("http.Request is nil")
			}
			if parentLogID == "" {
				panic("empty to parentLogID")
			}
			if s.logID == parentLogID {
				panic("do not make parentLogID and the argument logID of 'NewLogging' functin identical")
			}

			severity := logging.Default
			traceID := newTraceID()
			ctx := s.Context()
			ctx = setSeverity(ctx, &severity)
			ctx = setTraceID(ctx, &traceID)

			res := &logResponse{code: http.StatusOK, origin: w}
			st := time.Now()
			next.ServeHTTP(res, r.WithContext(ctx))
			et := time.Now()
			if r.URL.String() == "" {
				r.URL.Path = "Empty_RequestUrl"
			}
			s.client.Logger(parentLogID, s.option...).Log(logging.Entry{
				HTTPRequest: &logging.HTTPRequest{
					Status:       res.code,
					ResponseSize: int64(len(res.body)),
					Request:      r,
					Latency:      et.Sub(st),
				},
				Timestamp: et,
				Trace:     traceID,
				Severity:  severity,
			})
		})
	}
}
