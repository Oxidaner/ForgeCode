package telemetry

import (
	"bytes"
	"strings"
	"testing"
)

func TestRedactorRedactsStructuredFieldsAndText(t *testing.T) {
	redactor := NewRedactor()

	fields := redactor.Redact([]Field{
		String("api_key", "sk-live-secret"),
		String("note", "token=abc123 should disappear"),
	}).([]Field)

	if fields[0].Value != Redacted {
		t.Fatalf("sensitive field was not redacted: %#v", fields[0].Value)
	}
	if strings.Contains(fields[1].Value.(string), "abc123") {
		t.Fatalf("free text token leaked: %q", fields[1].Value)
	}
}

func TestJSONLoggerRedactsOutput(t *testing.T) {
	var buf bytes.Buffer
	logger := NewJSONLogger(&buf, NewRedactor())

	logger.Info("calling provider token=abc123", String("authorization", "Bearer secret-token"))

	output := buf.String()
	for _, secret := range []string{"abc123", "secret-token"} {
		if strings.Contains(output, secret) {
			t.Fatalf("logger leaked secret %q in %s", secret, output)
		}
	}
	if !strings.Contains(output, Redacted) {
		t.Fatalf("logger output did not include redaction marker: %s", output)
	}
}
