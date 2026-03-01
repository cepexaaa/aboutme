package application

import (
	"math/rand"
	"testing"
	"time"
	"unicode"

	"gitlab.education.tbank.ru/backend-academy-go-2025/homeworks/hw5-benchmark/internal/domain"
)

func TestCreateInstanceFromClassInfo(t *testing.T) {
	tests := []struct {
		name      string
		classInfo *domain.ClassInfo
		wantErr   bool
	}{
		{
			name: "Simple struct with public fields",
			classInfo: &domain.ClassInfo{
				ClassName: "SimpleStruct",
				Fields: []domain.FieldInfo{
					{Name: "Name", Type: "string"},
					{Name: "Age", Type: "int"},
					{Name: "Score", Type: "float64"},
					{Name: "Active", Type: "bool"},
				},
			},
			wantErr: false,
		},
		{
			name: "Struct with private field",
			classInfo: &domain.ClassInfo{
				ClassName: "MixedStruct",
				Fields: []domain.FieldInfo{
					{Name: "PublicField", Type: "string"},
					{Name: "privateField", Type: "string"},
				},
			},
			wantErr: false,
		},
		{
			name: "Struct with slice and map",
			classInfo: &domain.ClassInfo{
				ClassName: "ComplexStruct",
				Fields: []domain.FieldInfo{
					{Name: "Names", Type: "[]string"},
					{Name: "Scores", Type: "[]int"},
					{Name: "Data", Type: "map[string]int"},
				},
			},
			wantErr: false,
		},
		{
			name: "Empty struct",
			classInfo: &domain.ClassInfo{
				ClassName: "EmptyStruct",
				Fields:    []domain.FieldInfo{},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			instance, err := createInstanceFromClassInfo(tt.classInfo)

			if tt.wantErr {
				if err == nil {
					t.Error("Expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("Failed to create instance: %v", err)
			}

			if instance == nil {
				t.Fatal("Instance should not be nil")
			}

			for _, field := range tt.classInfo.Fields {
				if len(field.Name) > 0 {
					firstChar := rune(field.Name[0])
					isPublic := unicode.IsUpper(firstChar)

					_, exists := instance[field.Name]

					if isPublic && !exists {
						t.Errorf("Public field %s should be created", field.Name)
					}

					if !isPublic && exists {
						t.Errorf("Private field %s should not be created", field.Name)
					}
				}
			}

			for fieldName, value := range instance {

				var fieldType string
				for _, field := range tt.classInfo.Fields {
					if field.Name == fieldName {
						fieldType = field.Type
						break
					}
				}

				if fieldType == "" {
					t.Errorf("Field %s not found in class info", fieldName)
					continue
				}

				switch fieldType {
				case "string":
					if _, ok := value.(string); !ok {
						t.Errorf("Field %s should be string, got %T", fieldName, value)
					}
				case "int", "int8", "int16", "int32", "int64":
					if _, ok := value.(int64); !ok {
						t.Errorf("Field %s should be int64, got %T", fieldName, value)
					}
				case "float32", "float64":
					if _, ok := value.(float64); !ok {
						t.Errorf("Field %s should be float64, got %T", fieldName, value)
					}
				case "bool":
					if _, ok := value.(bool); !ok {
						t.Errorf("Field %s should be bool, got %T", fieldName, value)
					}
				case "[]string":
					if _, ok := value.([]string); !ok {
						t.Errorf("Field %s should be []string, got %T", fieldName, value)
					}
				case "[]int":
					if _, ok := value.([]int); !ok {
						t.Errorf("Field %s should be []int, got %T", fieldName, value)
					}
				case "map[string]int":
					if _, ok := value.(map[string]int); !ok {
						t.Errorf("Field %s should be map[string]int, got %T", fieldName, value)
					}
				default:

					if value != nil {
						if _, ok := value.(map[string]interface{}); !ok {

							if fieldType[0] != '*' {
								t.Errorf("Field %s with type %s should be map or nil, got %T",
									fieldName, fieldType, value)
							}
						}
					}
				}
			}
		})
	}
}

func TestCreateValueForTypeString(t *testing.T) {
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))

	tests := []struct {
		name      string
		typeStr   string
		depth     int
		checkFunc func(interface{}) bool
	}{
		{
			name:    "string type",
			typeStr: "string",
			depth:   0,
			checkFunc: func(v interface{}) bool {
				str, ok := v.(string)
				return ok && len(str) == 10
			},
		},
		{
			name:    "int type",
			typeStr: "int",
			depth:   0,
			checkFunc: func(v interface{}) bool {
				_, ok := v.(int64)
				return ok
			},
		},
		{
			name:    "float type",
			typeStr: "float64",
			depth:   0,
			checkFunc: func(v interface{}) bool {
				_, ok := v.(float64)
				return ok
			},
		},
		{
			name:    "bool type",
			typeStr: "bool",
			depth:   0,
			checkFunc: func(v interface{}) bool {
				_, ok := v.(bool)
				return ok
			},
		},
		{
			name:    "string slice",
			typeStr: "[]string",
			depth:   0,
			checkFunc: func(v interface{}) bool {
				slice, ok := v.([]string)
				if !ok {
					return false
				}
				return len(slice) >= 1 && len(slice) <= 3
			},
		},
		{
			name:    "int slice",
			typeStr: "[]int",
			depth:   0,
			checkFunc: func(v interface{}) bool {
				slice, ok := v.([]int)
				if !ok {
					return false
				}
				return len(slice) >= 1 && len(slice) <= 3
			},
		},
		{
			name:    "string-int map",
			typeStr: "map[string]int",
			depth:   0,
			checkFunc: func(v interface{}) bool {
				m, ok := v.(map[string]int)
				if !ok {
					return false
				}
				return len(m) >= 1 && len(m) <= 3
			},
		},
		{
			name:    "depth limit",
			typeStr: "string",
			depth:   4,
			checkFunc: func(v interface{}) bool {
				return v == nil
			},
		},
		{
			name:    "pointer type",
			typeStr: "*string",
			depth:   0,
			checkFunc: func(v interface{}) bool {
				return v == nil
			},
		},
		{
			name:    "complex type",
			typeStr: "SomeStruct",
			depth:   0,
			checkFunc: func(v interface{}) bool {
				_, ok := v.(map[string]interface{})
				return ok
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := createValueForTypeString(tt.typeStr, rnd, tt.depth)

			if !tt.checkFunc(result) {
				t.Errorf("createValueForTypeString(%q, depth=%d) = %v, doesn't pass check",
					tt.typeStr, tt.depth, result)
			}
		})
	}
}

func TestRandomString(t *testing.T) {
	rnd := rand.New(rand.NewSource(42))

	results := make(map[string]bool)
	for i := 0; i < 100; i++ {
		str := randomString(rnd)

		if len(str) != 10 {
			t.Errorf("randomString() length = %d, want 10", len(str))
		}

		for _, ch := range str {
			if !((ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z')) {
				t.Errorf("randomString() contains invalid character: %c", ch)
			}
		}

		if results[str] {
			t.Errorf("randomString() produced duplicate: %s", str)
		}
		results[str] = true
	}
}

func TestDepthLimit(t *testing.T) {

	classInfo := &domain.ClassInfo{
		ClassName: "RecursiveStruct",
		Fields: []domain.FieldInfo{
			{Name: "Value", Type: "int"},
			{Name: "Child", Type: "RecursiveStruct"},
		},
	}

	instance, err := createInstanceFromClassInfo(classInfo)
	if err != nil {
		t.Fatalf("Failed to create instance: %v", err)
	}

	if val, ok := instance["Value"]; !ok {
		t.Error("Value field should be created")
	} else if _, ok := val.(int64); !ok {
		t.Error("Value should be int64")
	}

	if child, ok := instance["Child"]; !ok {
		t.Error("Child field should be created")
	} else if childMap, ok := child.(map[string]interface{}); !ok {
		t.Errorf("Child should be map[string]interface{} due to depth limit, got %T", child)
	} else if len(childMap) != 0 {
		t.Error("Child map should be empty due to depth limit")
	}
}

func TestCreateFunction(t *testing.T) {
	cfg := &domain.Config{
		Format: "TEXT",
	}

	classInfo := &domain.ClassInfo{
		ClassName: "TestStruct",
		Fields: []domain.FieldInfo{
			{Name: "Name", Type: "string"},
			{Name: "Count", Type: "int"},
		},
	}
	Create(classInfo, cfg)

	if classInfo.Instance == nil {
		t.Error("Create should set Instance field")
	}

	if len(classInfo.Instance) != 2 {
		t.Errorf("Instance should have 2 fields, got %d", len(classInfo.Instance))
	}

	if _, ok := classInfo.Instance["Name"]; !ok {
		t.Error("Instance should have Name field")
	}

	if _, ok := classInfo.Instance["Count"]; !ok {
		t.Error("Instance should have Count field")
	}
}
