package telemetry

import (
	"fmt"
	"regexp"
	"strings"
)

const Redacted = "[REDACTED]"

type Redactor interface {
	Redact(v any) any
	RedactText(text string) string
}

type RedactionRule struct {
	KeyPatterns []string
}

type DefaultRedactor struct {
	keyFragments []string
	patterns     []*regexp.Regexp
}

func NewRedactor(extraKeys ...string) *DefaultRedactor {
	keys := []string{
		"api_key",
		"apikey",
		"authorization",
		"password",
		"secret",
		"token",
	}
	keys = append(keys, extraKeys...)

	return &DefaultRedactor{
		keyFragments: keys,
		patterns: []*regexp.Regexp{
			regexp.MustCompile(`(?i)\b(api[_-]?key|token|password|secret)\s*[:=]\s*([^\s,;"'}]+)`),
			regexp.MustCompile(`(?i)\b(authorization)\s*[:=]\s*(bearer\s+)?([^\s,;"'}]+)`),
		},
	}
}

func (r *DefaultRedactor) Redact(v any) any {
	switch typed := v.(type) {
	case nil:
		return nil
	case string:
		return r.RedactText(typed)
	case Field:
		return Field{Key: typed.Key, Value: r.redactByKey(typed.Key, typed.Value)}
	case []Field:
		out := make([]Field, len(typed))
		for i, f := range typed {
			out[i] = r.Redact(f).(Field)
		}
		return out
	case map[string]any:
		out := make(map[string]any, len(typed))
		for k, value := range typed {
			out[k] = r.redactByKey(k, value)
		}
		return out
	case []any:
		out := make([]any, len(typed))
		for i, value := range typed {
			out[i] = r.Redact(value)
		}
		return out
	default:
		return r.RedactText(fmt.Sprint(typed))
	}
}

func (r *DefaultRedactor) RedactText(text string) string {
	out := text
	for _, pattern := range r.patterns {
		out = pattern.ReplaceAllStringFunc(out, func(match string) string {
			pieces := strings.SplitN(match, "=", 2)
			if len(pieces) == 2 {
				return pieces[0] + "=" + Redacted
			}
			pieces = strings.SplitN(match, ":", 2)
			if len(pieces) == 2 {
				return pieces[0] + ": " + Redacted
			}
			return Redacted
		})
	}
	return out
}

func (r *DefaultRedactor) redactByKey(key string, value any) any {
	if r.isSensitiveKey(key) {
		return Redacted
	}
	return r.Redact(value)
}

func (r *DefaultRedactor) isSensitiveKey(key string) bool {
	normalized := strings.ToLower(strings.ReplaceAll(strings.ReplaceAll(key, "-", "_"), ".", "_"))
	for _, fragment := range r.keyFragments {
		if strings.Contains(normalized, strings.ToLower(fragment)) {
			return true
		}
	}
	return false
}
