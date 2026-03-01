package domain

type FieldInfo struct {
	Name string `json:"name"`
	Type string `json:"type"`
	Tag  string `json:"tag,omitempty"`
}

type MethodInfo struct {
	Name       string   `json:"name"`
	Params     []string `json:"params"`
	ReturnType string   `json:"returnType,omitempty"`
}

type ClassInfo struct {
	ClassName  string                 `json:"class"`
	Superclass string                 `json:"superclass,omitempty"`
	Interfaces []string               `json:"interfaces,omitempty"`
	Fields     []FieldInfo            `json:"fields"`
	Methods    []MethodInfo           `json:"methods"`
	Hierarchy  interface{}            `json:"hierarchy,omitempty"`
	Instance   map[string]interface{} `json:"instance,omitempty"`
}

type MethodParam struct {
	Name string `json:"name,omitempty"`
	Type string `json:"type"`
}
