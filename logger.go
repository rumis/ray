package ray

import (
	"fmt"
	"io"
	"time"
)

func init() {
	defaultLogger = StdLogger
}

// LoggerHandle logger handle defined
type LoggerHandle func(opt *Options, err error) error

// default logger
var defaultLogger LoggerHandle

// SetGlobalLogger set global logger
func SetGlobalLogger(logger LoggerHandle) {
	defaultLogger = logger
}

// StdLogger default logger
func StdLogger(opt *Options, err error) error {
	logInfo := "ray trace \n"
	if opt.URL != "" {
		logInfo += "url:" + opt.URL + "\n"
	}
	if opt.Method != "" {
		logInfo += "method:" + opt.Method + "\n"
	}
	if opt.RetryTimes != 0 {
		logInfo += "retry:" + fmt.Sprintf("%d", opt.RetryTimes) + "\n"
	}
	if opt.Timeout != 0 {
		logInfo += "timeout:" + fmt.Sprintf("%d", opt.Timeout) + "\n"
	}
	if opt.ContentType != "" {
		logInfo += "content-type:" + opt.ContentType + "\n"
	}
	if opt.Query != nil {
		logInfo += "query:" + fmt.Sprintf("%+v", opt.Query) + "\n"
	}
	if opt.Header != nil {
		logInfo += "header:" + fmt.Sprintf("%+v", opt.Header) + "\n"
	}
	if opt.Body != nil {
		opt.Body.Seek(0, io.SeekStart)
		b, err := io.ReadAll(opt.Body)
		if err != nil {
			return err
		}
		logInfo += "body:" + string(b) + "\n"
	}
	logInfo += "time:" + time.Now().Format(time.DateTime)

	fmt.Println(logInfo)

	return nil
}
