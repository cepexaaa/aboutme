package domain

type Word struct {
	Name  string   `json:"name"`
	Hints []string `json:"hints"`
}

type ComplexityWords struct {
	Easy []Word `json:"easy"`
	Hard []Word `json:"hard"`
}

type Words struct {
	En ComplexityWords `json:"en"`
	Ru ComplexityWords `json:"ru"`
}
