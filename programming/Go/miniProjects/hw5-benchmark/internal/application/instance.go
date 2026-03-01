package application

import (
	"math/rand"
	"time"
	"unicode"

	"gitlab.education.tbank.ru/backend-academy-go-2025/homeworks/hw5-benchmark/internal/domain"
)

func Create(classInfo *domain.ClassInfo, cfg *domain.Config) {
	instance := RunAndCheck(createInstanceFromClassInfo, classInfo, "Error creating instance: %v\n")
	classInfo.Instance = instance
}

func randomString(rnd *rand.Rand) string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	b := make([]byte, 10)
	for i := range b {
		b[i] = letters[rnd.Intn(len(letters))]
	}
	return string(b)
}

func createInstanceFromClassInfo(info *domain.ClassInfo) (map[string]interface{}, error) {
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))

	instanceData := make(map[string]interface{})

	for _, field := range info.Fields {
		if len(field.Name) > 0 && unicode.IsUpper(rune(field.Name[0])) {
			value := createValueForTypeString(field.Type, rnd, 0)
			instanceData[field.Name] = value
		}
	}

	return instanceData, nil
}
func createValueForTypeString(typeStr string, rnd *rand.Rand, depth int) interface{} {
	if depth > 3 {
		return nil
	}

	switch typeStr {
	case "string":
		return randomString(rnd)
	case "int", "int8", "int16", "int32", "int64":
		return rnd.Int63()
	case "uint", "uint8", "uint16", "uint32", "uint64":
		return rnd.Uint64()
	case "float32", "float64":
		return rnd.Float64()
	case "bool":
		return rnd.Int()%2 == 0
	case "[]string":
		length := rnd.Intn(3) + 1
		result := make([]string, length)
		for i := 0; i < length; i++ {
			result[i] = randomString(rnd)
		}
		return result
	case "[]int":
		length := rnd.Intn(3) + 1
		result := make([]int, length)
		for i := 0; i < length; i++ {
			result[i] = int(rnd.Int63())
		}
		return result
	case "map[string]int":
		length := rnd.Intn(3) + 1
		result := make(map[string]int)
		for range length {
			result[randomString(rnd)] = int(rnd.Int63())
		}
		return result
	default:
		if len(typeStr) > 0 && typeStr[0] == '*' {
			return nil
		}
		return make(map[string]interface{})
	}
}
