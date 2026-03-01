package domain

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
)

func TestTextFormatter(t *testing.T) {
	info := &ClassInfo{
		ClassName: "Person",
		Fields: []FieldInfo{
			{Name: "Name", Type: "string", Tag: `json:"name"`},
			{Name: "Age", Type: "int", Tag: `json:"age"`},
		},
		Methods: []MethodInfo{
			{Name: "GetName", ReturnType: "string"},
			{Name: "SetAge", Params: []string{"int"}},
		},
	}

	formatter := &TextFormatter{}
	var buf bytes.Buffer

	err := formatter.Format(info, &buf)
	if err != nil {
		t.Fatalf("Format failed: %v", err)
	}

	output := buf.String()

	expected := []string{
		"Person",
		"Fields:",
		"  - Name (string) [json:\"name\"]",
		"  - Age (int) [json:\"age\"]",
		"Methods:",
		"  - GetName() : string",
		"  - SetAge(int)",
	}

	for _, exp := range expected {
		if !strings.Contains(output, exp) {
			t.Errorf("Output doesn't contain expected text: %s", exp)
			t.Error(output)
		}
	}
}

func TestJSONFormatter(t *testing.T) {
	info := &ClassInfo{
		ClassName: "Person",
		Fields: []FieldInfo{
			{Name: "Name", Type: "string"},
			{Name: "Age", Type: "int"},
		},
		Methods: []MethodInfo{
			{Name: "GetName", ReturnType: "string"},
		},
	}

	formatter := &JSONFormatter{}
	var buf bytes.Buffer

	err := formatter.Format(info, &buf)
	if err != nil {
		t.Fatalf("Format failed: %v", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &result); err != nil {
		t.Fatalf("Invalid JSON output: %v\nOutput: %s", err, buf.String())
	}

	if result["class"] != "Person" {
		t.Errorf("Expected class 'Person', got %v", result["class"])
	}

	fields, ok := result["fields"].([]interface{})
	if !ok || len(fields) != 2 {
		t.Errorf("Expected 2 fields, got %v", fields)
	}

	methods, ok := result["methods"].([]interface{})
	if !ok || len(methods) != 1 {
		t.Errorf("Expected 1 method, got %v", methods)
	}
}

func TestFormatterFactory(t *testing.T) {
	tests := []struct {
		name    string
		format  string
		wantErr bool
	}{
		{"Valid TEXT", "TEXT", false},
		{"Valid JSON", "JSON", false},
		{"Valid lowercase text", "text", false},
		{"Valid lowercase json", "json", false},
		{"Invalid format", "XML", true},
		{"Empty format", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			formatter, err := FormatterFactory(tt.format)

			if tt.wantErr {
				if err == nil {
					t.Error("Expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			if formatter == nil {
				t.Error("Formatter should not be nil")
			}
		})
	}
}
