package ledger

import (
	"fmt"
	"net/http"
	"os"
	"time"
)

var (
	std = New(os.Stdout, DebugLevel)
)

// Debug writes a log entry with the debug log level
func Debug(args ...interface{}) {
	std.write(DebugLevel, args...)
}

// Debugf writes a string formatted log entry with the debug log level
func Debugf(f string, args ...interface{}) {
	std.write(DebugLevel, fmt.Sprintf(f, args...))
}

// Info writes a log entry with the info log level
func Info(args ...interface{}) {
	std.write(InfoLevel, args...)
}

// Infof writes a string formatted log entry with the info log level
func Infof(f string, args ...interface{}) {
	std.write(InfoLevel, fmt.Sprintf(f, args...))
}

// Warn writes a log entry with the warn log level
func Warn(args ...interface{}) {
	std.write(WarnLevel, args...)
}

// Warnf writes a string formatted log entry with the warn log level
func Warnf(f string, args ...interface{}) {
	std.write(WarnLevel, fmt.Sprintf(f, args...))
}

// Error writes a log entry with the error log level
func Error(args ...interface{}) {
	std.write(ErrorLevel, args...)
}

// Errorf writes a string formatted log entry with the error log level
func Errorf(f string, args ...interface{}) {
	std.write(ErrorLevel, fmt.Sprintf(f, args...))
}

// Fatal writes a log entry with the fatal log level
func Fatal(args ...interface{}) {
	std.write(FatalLevel, args...)
}

// Fatalf writes a string formatted log entry with the fatal log level
func Fatalf(f string, args ...interface{}) {
	std.write(FatalLevel, fmt.Sprintf(f, args...))
}

// Threshold sets the log level threshold for the exported logger
func Threshold(level Level) {
	std.threshold = level
}

// EndpointInfo wraps the given http handler and logs details of the endpoint at
// the Info threshold
func EndpointInfo(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		n := time.Now()
		next.ServeHTTP(w, r)
		std.Infof("method=%v path=%v duration=%v", r.Method, r.URL.String(), time.Since(n))
	})
}

// EndpointDebug wraps the given http handler and logs details of the endpoint
// at the Debug threshold
func EndpointDebug(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		n := time.Now()
		next.ServeHTTP(w, r)
		std.Debugf("method=%v path=%v duration=%v", r.Method, r.URL.String(), time.Since(n))
	})
}
