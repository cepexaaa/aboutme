package application

import (
	"fmt"
	"strings"

	"fractalflame/internal/domain"
	"fractalflame/internal/infrastructure"
)

func LoadConfig() (*domain.Config, error) {
	cliConfig := infrastructure.ParseCLI()

	var err error
	var fileConfig *domain.Config
	if cliConfig.ConfigPath != "" {
		fileConfig, err = infrastructure.LoadConfigFromFile(cliConfig.ConfigPath)
		if err != nil {
			return nil, fmt.Errorf("failed to load config file: %w", err)
		}
	}

	config := mergeConfigs(cliConfig, fileConfig)

	return config, nil
}

// Что если есть конфигурации и из файла, и из терминала?
// Наверное человек хотел создать с конфигурацией из файла но со своими доработками (и ввёл их через cmd)
// Поэтому я обновляю значения из файла значениями из терминала, которые даны
func mergeConfigs(cliConfig *infrastructure.CLIConfig, fileConfig *domain.Config) *domain.Config {
	config := getDefaultConfig()
	if fileConfig != nil {
		applyFileConfig(config, fileConfig)
	}
	applyCLIConfig(config, cliConfig)
	return config
}

func getDefaultConfig() *domain.Config {
	return &domain.Config{
		Size: domain.Size{
			Width:  1920,
			Height: 1080,
		},
		Seed:            5.0,
		IterationCount:  2500,
		OutputPath:      "result.png",
		Threads:         1,
		GammaCorrection: false,
		Gamma:           2.2,
		SymmetryLevel:   1,
		Functions: []domain.FunctionConfig{
			{Name: "swirl", Weight: 1.0},
			{Name: "horseshoe", Weight: 0.8},
		},
		AffineParams: []domain.AffineParams{
			{A: 0.7, B: -0.3, C: 0.1, D: 0.3, E: 0.7, F: 0.1},
			{A: 0.3, B: 0.7, C: -0.2, D: -0.7, E: 0.3, F: 0.2},
		},
	}
}

func applyFileConfig(config *domain.Config, fileConfig *domain.Config) {
	if fileConfig.Size.Width > 0 && fileConfig.Size.Height > 0 {
		config.Size = fileConfig.Size
	}
	if fileConfig.Seed != 0 {
		config.Seed = fileConfig.Seed
	}
	if fileConfig.IterationCount > 0 {
		config.IterationCount = fileConfig.IterationCount
	}
	if fileConfig.OutputPath != "" {
		config.OutputPath = fileConfig.OutputPath
	}
	if fileConfig.Threads > 0 {
		config.Threads = fileConfig.Threads
	}
	config.GammaCorrection = fileConfig.GammaCorrection
	if fileConfig.Gamma > 0 {
		config.Gamma = fileConfig.Gamma
	}
	if fileConfig.SymmetryLevel >= 1 {
		config.SymmetryLevel = fileConfig.SymmetryLevel
	}
	if len(fileConfig.Functions) > 0 {
		config.Functions = fileConfig.Functions
	}
	if len(fileConfig.AffineParams) > 0 {
		config.AffineParams = fileConfig.AffineParams
	}
}

func applyCLIConfig(config *domain.Config, cliConfig *infrastructure.CLIConfig) {
	if cliConfig.Width > 0 {
		config.Size.Width = cliConfig.Width
	}
	if cliConfig.Height > 0 {
		config.Size.Height = cliConfig.Height
	}
	if cliConfig.Seed != 0 {
		config.Seed = cliConfig.Seed
	}
	if cliConfig.IterationCount > 0 {
		config.IterationCount = cliConfig.IterationCount
	}
	if cliConfig.OutputPath != "" {
		config.OutputPath = cliConfig.OutputPath
	}
	if cliConfig.Threads > 0 {
		config.Threads = cliConfig.Threads
	}
	config.GammaCorrection = cliConfig.GammaCorrection
	if cliConfig.Gamma > 0 {
		config.Gamma = cliConfig.Gamma
	}
	if cliConfig.SymmetryLevel >= 1 {
		config.SymmetryLevel = cliConfig.SymmetryLevel
	}

	if cliConfig.FunctionsStr != "" {
		funcs := parseFunctions(cliConfig.FunctionsStr)
		if len(funcs) > 0 {
			config.Functions = funcs
		}
	}

	if cliConfig.AffineParamsStr != "" {
		params := parseAffineParams(cliConfig.AffineParamsStr)
		if len(params) > 0 {
			config.AffineParams = params
		}
	}
}

func parseFunctions(funcStr string) []domain.FunctionConfig {
	var funcs []domain.FunctionConfig

	parts := strings.Split(funcStr, ",")
	for _, part := range parts {
		subparts := strings.Split(part, ":")
		if len(subparts) == 2 {
			name := strings.TrimSpace(subparts[0])
			weight := 1.0

			if w, err := parseFloat(subparts[1]); err == nil && w > 0 {
				weight = w
			}

			funcs = append(funcs, domain.FunctionConfig{
				Name:   name,
				Weight: weight,
			})
		}
	}

	return funcs
}

func parseAffineParams(paramsStr string) []domain.AffineParams {
	params := make([]domain.AffineParams, 0, 2)

	groups := strings.Split(paramsStr, "/")
	for _, group := range groups {
		values := strings.Split(group, ",")
		if len(values) == 6 {
			a, _ := parseFloat(values[0])
			b, _ := parseFloat(values[1])
			c, _ := parseFloat(values[2])
			d, _ := parseFloat(values[3])
			e, _ := parseFloat(values[4])
			f, _ := parseFloat(values[5])

			params = append(params, domain.AffineParams{
				A: a, B: b, C: c,
				D: d, E: e, F: f,
			})
		}
	}

	return params
}

func parseFloat(s string) (float64, error) {
	var f float64
	_, err := fmt.Sscanf(strings.TrimSpace(s), "%f", &f)
	return f, err
}
