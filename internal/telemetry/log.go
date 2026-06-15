package telemetry

import (
	"encoding/json"
	"io"
	"sync"
	"time"
)

type Logger interface {
	Info(msg string, kv ...Field)
	Error(msg string, err error, kv ...Field)
	With(kv ...Field) Logger
}

type JSONLogger struct {
	mu       sync.Mutex
	out      io.Writer
	redactor Redactor
	fields   []Field
}

func NewJSONLogger(out io.Writer, redactor Redactor) *JSONLogger {
	if redactor == nil {
		redactor = NewRedactor()
	}
	return &JSONLogger{out: out, redactor: redactor}
}

func (l *JSONLogger) Info(msg string, kv ...Field) {
	l.write("info", msg, nil, kv...)
}

func (l *JSONLogger) Error(msg string, err error, kv ...Field) {
	l.write("error", msg, err, kv...)
}

func (l *JSONLogger) With(kv ...Field) Logger {
	next := &JSONLogger{
		out:      l.out,
		redactor: l.redactor,
		fields:   make([]Field, 0, len(l.fields)+len(kv)),
	}
	next.fields = append(next.fields, l.fields...)
	next.fields = append(next.fields, kv...)
	return next
}

func (l *JSONLogger) write(level, msg string, err error, kv ...Field) {
	entry := map[string]any{
		"level": level,
		"msg":   l.redactor.RedactText(msg),
		"time":  time.Now().UTC().Format(time.RFC3339Nano),
	}
	if err != nil {
		entry["error"] = l.redactor.RedactText(err.Error())
	}

	fields := append([]Field{}, l.fields...)
	fields = append(fields, kv...)
	for _, field := range fields {
		redacted := l.redactor.Redact(field).(Field)
		entry[redacted.Key] = redacted.Value
	}

	encoded, encodeErr := json.Marshal(entry)
	if encodeErr != nil {
		return
	}

	l.mu.Lock()
	defer l.mu.Unlock()
	_, _ = l.out.Write(append(encoded, '\n'))
}
