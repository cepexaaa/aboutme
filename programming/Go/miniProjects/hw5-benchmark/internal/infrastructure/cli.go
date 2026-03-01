package infrastructure

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/pflag"
	"gitlab.education.tbank.ru/backend-academy-go-2025/homeworks/hw5-benchmark/internal/domain"
)

func ParseCLI() (cfg *domain.Config) {
	cfg = &domain.Config{}

	pflag.StringVarP(&cfg.PkgPath, "package", "p", "", "Go package path")
	pflag.StringVarP(&cfg.StructName, "struct", "s", "", "Struct name")
	pflag.StringVarP(&cfg.Format, "format", "f", "TEXT", "Output format (TEXT or JSON)")
	pflag.BoolVarP(&cfg.Help, "help", "h", false, "Show help")

	pflag.Parse()

	cfg.Format = strings.ToUpper(cfg.Format)

	return
}

func Usage() {
	fmt.Fprintf(os.Stderr, "Usage: inspector [options]\n")
	fmt.Fprintf(os.Stderr, "\nOptions:\n")
	flag.PrintDefaults()
	fmt.Fprintf(os.Stderr, "\nExamples:\n")
	fmt.Fprintf(os.Stderr, "  inspector -package ./testdata -struct Person -format TEXT\n")
	fmt.Fprintf(os.Stderr, "  inspector -p . -s User -f JSON\n")
}
