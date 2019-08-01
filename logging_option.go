package glbr

import (
	"time"

	"cloud.google.com/go/logging"
	"google.golang.org/genproto/googleapis/api/monitoredres"
)

const (
	// GAEApplication .
	GAEApplication = "gae_app"
	//CloudFunction .
	CloudFunction = "cloud_function"
)

// Option option interface
type Option interface {
	loggerOption() logging.LoggerOption
}

// Option log service option
func (s Service) Option(opts ...Option) Service {
	s.option = make([]logging.LoggerOption, 0)
	for _, opt := range opts {
		if opt != nil {
			s.option = append(s.option, opt.loggerOption())
		}
	}
	return s
}

// Label ログエントリに付加する共通ラベル
func Label(label map[string]string) Option { return labelOption(label) }

type labelOption map[string]string

func (o labelOption) loggerOption() logging.LoggerOption {
	return logging.CommonLabels(o)
}

// MonitoredResource ログエントリに付加するリソースラベル
// https://cloud.google.com/monitoring/api/resources のResourceTypeのLabelsを自動で補完します。
// Default: resourceType = project, resourceLabel = {"project_id": $PROJECT_ID}
func MonitoredResource(resourceType string, resourceLabel map[string]string) Option {
	return monitoredResourceOption{&monitoredres.MonitoredResource{
		Type:   resourceType,
		Labels: resourceLabel,
	}}
}

type monitoredResourceOption struct {
	mr *monitoredres.MonitoredResource
}

func (o monitoredResourceOption) loggerOption() logging.LoggerOption {
	return logging.CommonResource(o.mr)
}

// ConcurrentWrite ログエントリの同時書き込み数　Default: 1
func ConcurrentWrite(limit int) Option { return concurrentOption(limit) }

type concurrentOption int

func (o concurrentOption) loggerOption() logging.LoggerOption {
	return logging.ConcurrentWriteLimit(int(o))
}

// WriteDelay ログエントリの遅延書き込み時間　Default: 1s
func WriteDelay(threshold int) Option { return writeDelayOption(threshold) }

type writeDelayOption time.Duration

func (o writeDelayOption) loggerOption() logging.LoggerOption {
	return logging.DelayThreshold(time.Duration(o))
}

// EntryCount バッファ可能なログエントリの最大数　Default: 1000
func EntryCount(threshold int) Option { return entryCountThresholdOption(threshold) }

type entryCountThresholdOption int

func (o entryCountThresholdOption) loggerOption() logging.LoggerOption {
	return logging.EntryCountThreshold(int(o))
}

// EntryByteThreshold バッファ可能なログエントリの最大サイズ　Default: 1MiB
func EntryByteThreshold(threshold int) Option { return entryByteOption(threshold) }

type entryByteOption int

func (o entryByteOption) loggerOption() logging.LoggerOption {
	return logging.EntryByteThreshold(int(o))
}

// EntryByteLimit 送信するログエントリの最大サイズ　Default: 0(無制限)
func EntryByteLimit(limit int) Option { return entryByteLimitOption(limit) }

type entryByteLimitOption int

func (o entryByteLimitOption) loggerOption() logging.LoggerOption {
	return logging.EntryByteLimit(int(o))
}

// BufferedByte ログバッファの最大サイズ　Default: 1GiB
func BufferedByte(limit int) Option { return bufferedByteLimitOption(limit) }

type bufferedByteLimitOption int

func (o bufferedByteLimitOption) loggerOption() logging.LoggerOption {
	return logging.BufferedByteLimit(int(o))
}
