package main

import (
	"fmt"
	"os"

	"gitlab.education.tbank.ru/backend-academy-go-2025/homeworks/hw5-benchmark/internal/application"
	"gitlab.education.tbank.ru/backend-academy-go-2025/homeworks/hw5-benchmark/internal/domain"
	"gitlab.education.tbank.ru/backend-academy-go-2025/homeworks/hw5-benchmark/internal/infrastructure"
)

func main() {
	cfg := infrastructure.ParseCLI()
	exit := application.Validate(cfg)
	if exit >= 0 {
		os.Exit(exit)
	}

	classInfo := application.Inspect(cfg)
	application.Create(classInfo, cfg)

	formatter := application.RunAndCheck(domain.FormatterFactory, cfg.Format, "Error output format: %v\n")
	err := formatter.Format(classInfo, os.Stdout)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error formatting: %v\n", err)
		os.Exit(1)
	}
}
