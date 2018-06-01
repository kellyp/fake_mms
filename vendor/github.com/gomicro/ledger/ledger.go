// Package ledger provides a threadsafe, minimalist logger on top of native Go
// logging.  Adding the ability to write to more than standard out and honor
// log level thresholds.
package ledger

import (
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"
)

// Ledger represents a logger
type Ledger struct {
	writer    io.Writer
	threshold Level
	mu        sync.Mutex
}

// New returns a new Ledger configured with the specified writer and log level
// threshold
func New(w io.Writer, level Level) *Ledger {
	return &Ledger{
		writer:    w,
		threshold: level,
	}
}

// Debug writes a log entry with the debug log level
func (l *Ledger) Debug(args ...interface{}) {
	l.write(DebugLevel, args...)
}

// Debugf writes a string formatted log entry with the debug log level
func (l *Ledger) Debugf(f string, args ...interface{}) {
	l.write(DebugLevel, fmt.Sprintf(f, args...))
}

// Info writes a log entry with the info log level
func (l *Ledger) Info(args ...interface{}) {
	l.write(InfoLevel, args...)
}

// Infof writes a string formatted log entry with the info log level
func (l *Ledger) Infof(f string, args ...interface{}) {
	l.write(InfoLevel, fmt.Sprintf(f, args...))
}

// Warn writes a log entry with the warn log level
func (l *Ledger) Warn(args ...interface{}) {
	l.write(WarnLevel, args...)
}

// Warnf writes a string formatted log entry with the warn log level
func (l *Ledger) Warnf(f string, args ...interface{}) {
	l.write(WarnLevel, fmt.Sprintf(f, args...))
}

// Error writes a log entry with the error log level
func (l *Ledger) Error(args ...interface{}) {
	l.write(ErrorLevel, args...)
}

// Errorf writes a string formatted log entry with the error log level
func (l *Ledger) Errorf(f string, args ...interface{}) {
	l.write(ErrorLevel, fmt.Sprintf(f, args...))
}

// Fatal writes a log entry with the fatal log level
func (l *Ledger) Fatal(args ...interface{}) {
	l.write(FatalLevel, args...)
}

// Fatalf writes a string formatted log entry with the fatal log level
func (l *Ledger) Fatalf(f string, args ...interface{}) {
	l.write(FatalLevel, fmt.Sprintf(f, args...))
}

func (l *Ledger) write(level Level, args ...interface{}) {
	if level <= l.threshold {
		out, args := fmt.Sprintf("[%v]", args[0]), args[1:]
		for _, arg := range args {
			out += fmt.Sprintf(" [%v]", arg)
		}

		l.mu.Lock()
		defer l.mu.Unlock()
		fmt.Fprintf(l.writer, "%s: %v\n", level, out)
	}
}

// Threshold sets the log level threshold for the defined logger
func (l *Ledger) Threshold(level Level) {
	l.threshold = level
}

// EndpointInfo wraps the given http handler and logs details of the endpoint at
// the Info threshold
func (l *Ledger) EndpointInfo(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		n := time.Now()
		next.ServeHTTP(w, r)
		l.Infof("method=%v path=%v duration=%v", r.Method, r.URL.String(), time.Since(n))
	})
}

// EndpointDebug wraps the given http handler and logs details of the endpoint
// at the Debug threshold
func (l *Ledger) EndpointDebug(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		n := time.Now()
		next.ServeHTTP(w, r)
		l.Debugf("method=%v path=%v duration=%v", r.Method, r.URL.String(), time.Since(n))
	})
}
