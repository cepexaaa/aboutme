package domain

type Config struct {
	Size            Size             `json:"size"`
	Seed            float64          `json:"seed"`
	IterationCount  int              `json:"iteration_count"`
	OutputPath      string           `json:"output_path"`
	Threads         int              `json:"threads"`
	Functions       []FunctionConfig `json:"functions"`
	AffineParams    []AffineParams   `json:"affine_params"`
	GammaCorrection bool             `json:"gamma_correction"`
	Gamma           float64          `json:"gamma"`
	SymmetryLevel   int              `json:"symmetry_level"`
}
