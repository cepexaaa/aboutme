package application

import (
	"fmt"
	"os"
	"path/filepath"

	"gitlab.education.tbank.ru/backend-academy-go-2025/homeworks/hw5-benchmark/internal/domain"
)

// Inspect[T any] - но типизация в даннойреализации не нужна, но по спецификации (по зпдпнию) - это нужная функция
func Inspect(cfg *domain.Config) *domain.ClassInfo {
	absPath := RunAndCheck(filepath.Abs, cfg.PkgPath, "Error: %v\n")
	parser := RunAndCheck(domain.NewASTParser, absPath, "Error parsing package: %v\n")
	classInfo := RunAndCheck(parser.AnalyzeStruct, cfg.StructName, "Error: %v\n")
	return classInfo
}

func RunAndCheck[R, T any](fn func(T) (R, error), args T, errMsg string) R {
	result, err := fn(args)
	if err != nil {
		fmt.Fprintf(os.Stderr, errMsg, err)
		os.Exit(1)
	}
	return result
}
