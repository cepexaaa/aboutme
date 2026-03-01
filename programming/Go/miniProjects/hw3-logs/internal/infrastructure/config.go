package infrastructure

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/pflag"
)

type OutputFormat string

const (
	JSONFormat     OutputFormat = "json"
	MarkdownFormat OutputFormat = "markdown"
	AsciiDocFormat OutputFormat = "adoc"
)

type Config struct {
	Path   []string
	Format OutputFormat
	Output string
	From   *time.Time
	To     *time.Time
}

func ParseFlags() (*Config, error) {
	var path string
	var format string
	var output string
	var fromStr string
	var toStr string

	pflag.StringVarP(&path, "path", "p", "", "Path to log-files NGINX")
	pflag.StringVarP(&format, "format", "f", "markdown", "Output foramt (json, markdown, adoc)")
	pflag.StringVarP(&output, "output", "o", "", "path to save file")
	pflag.StringVar(&fromStr, "from", "", "First date (ISO8601)")
	pflag.StringVar(&toStr, "to", "", "Last date (ISO8601)")
	pflag.Parse()

	if path == "" {
		return nil, errors.New("parameter --path required")
	}

	if output == "" {
		return nil, errors.New("parameter --output required")
	}

	if format == "" {
		return nil, errors.New("parameter --format required")
	}

	fmt.Println(path)

	files, err := validateInputPath(path)

	if err != nil {
		return nil, err
	}

	cfg := &Config{
		Path:   files,
		Output: output,
	}

	switch OutputFormat(format) {
	case JSONFormat, MarkdownFormat, AsciiDocFormat:
		cfg.Format = OutputFormat(format)
	default:
		return nil, fmt.Errorf("unsupported format: %s", format)
	}

	if err := validateOutputExtension(cfg.Format, output); err != nil {
		return nil, err
	}

	if _, err := os.Stat(output); err == nil {
		return nil, fmt.Errorf("file is already exist: %s", output)
	}

	dir := filepath.Dir(output)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return nil, fmt.Errorf("the directory does not exist: %s", dir)
	}

	if fromStr != "" {
		from, err := time.Parse(time.RFC3339, fromStr)
		if err != nil {
			return nil, fmt.Errorf("incorrect date format --from: %v", err)
		}
		cfg.From = &from
	}

	if toStr != "" {
		to, err := time.Parse(time.RFC3339, toStr)
		if err != nil {
			return nil, fmt.Errorf("incorrect date format --to: %v", err)
		}
		cfg.To = &to
	}

	if cfg.From != nil && cfg.To != nil && !cfg.From.Before(*cfg.To) {
		return nil, errors.New("date --from it should be earlier than --to")
	}

	return cfg, nil
}

func validateOutputExtension(format OutputFormat, output string) error {
	ext := filepath.Ext(output)

	switch format {
	case JSONFormat:
		if ext != ".json" {
			return errors.New("for the format json an extension is required .json")
		}
	case MarkdownFormat:
		if ext != ".md" {
			return errors.New("for the format markdown an extension is required .md")
		}
	case AsciiDocFormat:
		if ext != ".adoc" && ext != ".ad" {
			return errors.New("for the format adoc an extension is required .adoc or .ad")
		}
	}

	return nil
}

func validateInputPath(path string) ([]string, error) {
	var files []string
	if strings.ContainsAny(path, "*?[") {
		matches, err := filepath.Glob(path)
		if err != nil {
			return nil, fmt.Errorf("error in the search template: %v", err)
		}
		if len(matches) == 0 {
			return nil, fmt.Errorf("files were not found according to the template: %s", path)
		}
		files = make([]string, 0, len(matches))

		for _, file := range matches {
			if err := validateSingleFile(file); err != nil {
				return nil, err
			}
			files = append(files, file)
		}
	} else {
		if err := validateSingleFile(path); err != nil {
			return nil, err
		}
		files = append(files, path)
	}

	return files, nil
}

func validateSingleFile(path string) error {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return fmt.Errorf("the file does not exist: %s", path)
	}
	if err != nil {
		return fmt.Errorf("file access error: %v", err)
	}
	if info.IsDir() {
		return fmt.Errorf("the path is a directory: %s", path)
	}

	ext := strings.ToLower(filepath.Ext(path))
	if ext != ".txt" && ext != ".log" {
		return fmt.Errorf("unsupported file extension: %s. Log files are supported only .txt и .log", path)
	}

	return nil
}
