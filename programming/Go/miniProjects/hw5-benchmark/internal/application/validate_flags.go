package application

import (
	"flag"
	"fmt"
	"os"

	"gitlab.education.tbank.ru/backend-academy-go-2025/homeworks/hw5-benchmark/internal/domain"
)

func Validate(cfg *domain.Config) (exit int) {
	if cfg.Help {
		flag.Usage()
		return 0
	}
	if cfg.PkgPath == "" || cfg.StructName == "" {
		fmt.Fprintf(os.Stderr, "Error: both package and struct parameters are required\n")
		flag.Usage()
		return 2
	}
	if _, err := os.Stat(cfg.PkgPath); os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "Error: package directory %s does not exist\n", cfg.PkgPath)
		return 2
	}
	if cfg.Format != "TEXT" && cfg.Format != "JSON" {
		fmt.Fprintf(os.Stderr, "Error: uncorrect output format '%s'\n", cfg.Format)
		return 2
	}
	return -1
}
