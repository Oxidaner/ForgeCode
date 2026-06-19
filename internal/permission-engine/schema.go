package permission

import (
	"encoding/json"
	"fmt"
	"strings"
	"unicode/utf8"

	toolruntime "forgecode/internal/tool-runtime"
)

type schemaSpec struct {
	Type                 string                    `json:"type"`
	Required             []string                  `json:"required"`
	Properties           map[string]schemaProperty `json:"properties"`
	AdditionalProperties *bool                     `json:"additionalProperties"`
}

type schemaProperty struct {
	Type      string `json:"type"`
	MinLength *int   `json:"minLength"`
	Minimum   *int   `json:"minimum"`
}

func validateInput(descriptor toolruntime.ToolDescriptor, raw json.RawMessage) (map[string]any, error) {
	if len(raw) == 0 {
		return nil, toolruntime.NewError(toolruntime.ValidationError, "tool input is required")
	}
	if !utf8.Valid(raw) {
		return nil, toolruntime.NewError(toolruntime.ValidationError, "tool input must be UTF-8")
	}
	if strings.Contains(string(raw), "\x00") {
		return nil, toolruntime.NewError(toolruntime.ValidationError, "tool input contains null byte")
	}
	if !json.Valid(raw) {
		return nil, toolruntime.NewError(toolruntime.ValidationError, "tool input must be valid JSON")
	}

	var input map[string]any
	if err := json.Unmarshal(raw, &input); err != nil {
		return nil, toolruntime.WrapError(toolruntime.ValidationError, "decode tool input", err)
	}

	var spec schemaSpec
	if err := json.Unmarshal(descriptor.InputSchema, &spec); err != nil {
		return nil, toolruntime.WrapError(toolruntime.ValidationError, "decode tool input schema", err)
	}
	if spec.Type != "" && spec.Type != "object" {
		return nil, toolruntime.NewError(toolruntime.ValidationError, "only object input schemas are supported")
	}
	for _, field := range spec.Required {
		if _, ok := input[field]; !ok {
			return nil, toolruntime.NewError(toolruntime.ValidationError, "missing required field: "+field)
		}
	}
	if spec.AdditionalProperties != nil && !*spec.AdditionalProperties {
		for key := range input {
			if _, ok := spec.Properties[key]; !ok {
				return nil, toolruntime.NewError(toolruntime.ValidationError, "unexpected field: "+key)
			}
		}
	}
	for key, value := range input {
		prop, ok := spec.Properties[key]
		if !ok {
			continue
		}
		if err := validateProperty(key, value, prop); err != nil {
			return nil, err
		}
	}
	if err := validateStrings(input); err != nil {
		return nil, err
	}
	return input, nil
}

func validateProperty(key string, value any, prop schemaProperty) error {
	switch prop.Type {
	case "string":
		str, ok := value.(string)
		if !ok {
			return toolruntime.NewError(toolruntime.ValidationError, key+" must be a string")
		}
		if prop.MinLength != nil && len(str) < *prop.MinLength {
			return toolruntime.NewError(toolruntime.ValidationError, key+" is too short")
		}
	case "integer":
		num, ok := value.(float64)
		if !ok || num != float64(int(num)) {
			return toolruntime.NewError(toolruntime.ValidationError, key+" must be an integer")
		}
		if prop.Minimum != nil && int(num) < *prop.Minimum {
			return toolruntime.NewError(toolruntime.ValidationError, fmt.Sprintf("%s must be >= %d", key, *prop.Minimum))
		}
	case "boolean":
		if _, ok := value.(bool); !ok {
			return toolruntime.NewError(toolruntime.ValidationError, key+" must be a boolean")
		}
	}
	return nil
}

func validateStrings(value any) error {
	switch typed := value.(type) {
	case string:
		if strings.Contains(typed, "\x00") || !utf8.ValidString(typed) {
			return toolruntime.NewError(toolruntime.ValidationError, "string input contains null byte or invalid UTF-8")
		}
	case map[string]any:
		for _, child := range typed {
			if err := validateStrings(child); err != nil {
				return err
			}
		}
	case []any:
		for _, child := range typed {
			if err := validateStrings(child); err != nil {
				return err
			}
		}
	}
	return nil
}
