package domain

import (
	"os"
	"path/filepath"
	"testing"
)

type TestStruct struct {
	Name  string `json:"name" validate:"required"`
	Age   int    `json:"age"`
	Email string `json:"email,omitempty"`
}

func (t TestStruct) GetName() string {
	return t.Name
}

func (t *TestStruct) SetEmail(email string) {
	t.Email = email
}

func TestASTParser_FindStruct(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.go")

	content := `package testdata

type Person struct {
	Name string
	Age  int
}
`

	if err := os.WriteFile(testFile, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	parser, err := NewASTParser(tmpDir)
	if err != nil {
		t.Fatalf("Failed to create parser: %v", err)
	}

	_, err = parser.FindStruct("Person")
	if err != nil {
		t.Errorf("Expected to find Person struct, got error: %v", err)
	}

	_, err = parser.FindStruct("NonExistent")
	if err == nil {
		t.Error("Expected error for non-existent struct")
	}
}

func TestASTParser_AnalyzeStruct(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.go")

	content := `package testdata

type Person struct {
	Name string ` + "`json:\"name\"`" + `
	Age  int    ` + "`json:\"age\"`" + `
}

func (p Person) GetName() string {
	return p.Name
}

func (p *Person) SetAge(age int) {
	p.Age = age
}
`

	if err := os.WriteFile(testFile, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	parser, err := NewASTParser(tmpDir)
	if err != nil {
		t.Fatalf("Failed to create parser: %v", err)
	}

	info, err := parser.AnalyzeStruct("Person")
	if err != nil {
		t.Fatalf("Failed to analyze struct: %v", err)
	}

	if info.ClassName != "Person" {
		t.Errorf("Expected class name 'Person', got %s", info.ClassName)
	}

	if len(info.Fields) != 2 {
		t.Errorf("Expected 2 fields, got %d", len(info.Fields))
	}

	if len(info.Methods) != 2 {
		t.Errorf("Expected 2 methods, got %d", len(info.Methods))
	}

	expectedFields := map[string]string{
		"Name": "string",
		"Age":  "int",
	}

	for _, field := range info.Fields {
		expectedType, ok := expectedFields[field.Name]
		if !ok {
			t.Errorf("Unexpected field: %s", field.Name)
			continue
		}
		if field.Type != expectedType {
			t.Errorf("Field %s: expected type %s, got %s", field.Name, expectedType, field.Type)
		}
	}
}

func TestTypeToString(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.go")

	content := `package testdata

type TestTypes struct {
	Name   string
	Age    int
	Scores []int
	Data   map[string]interface{}
	Ptr    *TestTypes
}
`

	if err := os.WriteFile(testFile, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	parser, err := NewASTParser(tmpDir)
	if err != nil {
		t.Fatalf("Failed to create parser: %v", err)
	}

	_, err = parser.FindStruct("TestTypes")
	if err != nil {
		t.Fatalf("Failed to find struct: %v", err)
	}
}
