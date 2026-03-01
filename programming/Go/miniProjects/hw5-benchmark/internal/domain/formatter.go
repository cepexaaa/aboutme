package domain

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

type Formatter interface {
	Format(*ClassInfo, io.Writer) error
}

type TextFormatter struct{}

func (f *TextFormatter) Format(info *ClassInfo, w io.Writer) error {
	var builder strings.Builder

	builder.WriteString(fmt.Sprintf("Class: %s\n", info.ClassName))

	if info.Superclass != "" {
		builder.WriteString(fmt.Sprintf("Embeds: %s\n", info.Superclass))
	}

	if len(info.Fields) > 0 {
		builder.WriteString("Fields:\n")
		for _, field := range info.Fields {
			tagInfo := ""
			if field.Tag != "" {
				tagInfo = fmt.Sprintf(" [%s]", field.Tag)
			}
			builder.WriteString(fmt.Sprintf("  - %s (%s)%s\n",
				field.Name, field.Type, tagInfo))
		}
	}

	if len(info.Methods) > 0 {
		builder.WriteString("Methods:\n")
		for _, method := range info.Methods {
			params := strings.Join(method.Params, ", ")
			returnType := ""
			if method.ReturnType != "" {
				returnType = fmt.Sprintf(" : %s", method.ReturnType)
			}
			builder.WriteString(fmt.Sprintf("  - %s(%s)%s\n", method.Name, params, returnType))
		}
	}

	if len(info.Interfaces) > 0 {
		builder.WriteString("Implements:\n")
		for _, iface := range info.Interfaces {
			builder.WriteString(fmt.Sprintf("  - %s\n", iface))
		}
	}

	builder.WriteString("Hierarchy:\n")
	writeHierarchyText(info.Hierarchy, &builder, 0)

	if info.Instance != nil {
		builder.WriteString("\nInstance:\n")
		for key, value := range info.Instance {
			builder.WriteString(fmt.Sprintf("  %s: %v\n", key, value))
		}
	}

	_, err := io.WriteString(w, builder.String())
	return err
}

func writeHierarchyText(hierarchy interface{}, buf *strings.Builder, level int) {
	if hierarchy == nil {
		return
	}

	if m, ok := hierarchy.(map[string]interface{}); ok {
		for key, value := range m {
			indent := strings.Repeat("    ", level)
			if level > 0 {
				indent += "└── "
			}
			buf.WriteString(fmt.Sprintf("%s%s\n", indent, key))

			if inner, ok := value.(map[string]interface{}); ok && len(inner) > 0 {
				writeHierarchyText(inner, buf, level+1)
			}
		}
	}
}

type JSONFormatter struct{}

func (f *JSONFormatter) Format(info *ClassInfo, w io.Writer) error {
	type JSONOutput struct {
		Class      string                 `json:"class"`
		Superclass string                 `json:"superclass,omitempty"`
		Interfaces []string               `json:"interfaces,omitempty"`
		Fields     []FieldInfo            `json:"fields"`
		Methods    []MethodInfo           `json:"methods"`
		Hierarchy  interface{}            `json:"hierarchy,omitempty"`
		Instance   map[string]interface{} `json:"instance,omitempty"`
	}

	output := JSONOutput{
		Class:      info.ClassName,
		Superclass: info.Superclass,
		Interfaces: info.Interfaces,
		Fields:     info.Fields,
		Methods:    info.Methods,
		Hierarchy:  info.Hierarchy,
		Instance:   info.Instance,
	}

	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	return encoder.Encode(output)
}

func FormatterFactory(format string) (Formatter, error) {
	switch strings.ToUpper(format) {
	case "TEXT":
		return &TextFormatter{}, nil
	case "JSON":
		return &JSONFormatter{}, nil
	default:
		return nil, fmt.Errorf("unsupported format: %s", format)
	}
}
