package telemetry

import "time"

type Field struct {
	Key   string
	Value any
}

type Tag struct {
	Key   string
	Value string
}

func String(key, value string) Field {
	return Field{Key: key, Value: value}
}

func Any(key string, value any) Field {
	return Field{Key: key, Value: value}
}

func Time(key string, value time.Time) Field {
	return Field{Key: key, Value: value}
}
