package domain

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"strings"
)

type ASTParser struct {
	fset       *token.FileSet
	pkg        *ast.Package
	dirPath    string
	structName string
}

func NewASTParser(dirPath string) (*ASTParser, error) {
	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, dirPath, nil, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	for _, pkg := range pkgs {
		return &ASTParser{
			fset:    fset,
			pkg:     pkg,
			dirPath: dirPath,
		}, nil
	}

	return nil, fmt.Errorf("no packages found in %s", dirPath)
}

func (p *ASTParser) FindStruct(structName string) (*ast.TypeSpec, error) {
	p.structName = structName
	for _, file := range p.pkg.Files {
		for _, decl := range file.Decls {
			if genDecl, ok := decl.(*ast.GenDecl); ok {
				for _, spec := range genDecl.Specs {
					if typeSpec, ok := spec.(*ast.TypeSpec); ok {
						if typeSpec.Name.Name == structName {
							if _, isStruct := typeSpec.Type.(*ast.StructType); isStruct {
								return typeSpec, nil
							}
						}
					}
				}
			}
		}
	}
	return nil, fmt.Errorf("struct %s not found", structName)
}

func (p *ASTParser) AnalyzeStruct(structName string) (*ClassInfo, error) {
	typeSpec, err := p.FindStruct(structName)
	if err != nil {
		return nil, err
	}

	info := &ClassInfo{
		ClassName:  typeSpec.Name.Name,
		Fields:     []FieldInfo{},
		Methods:    []MethodInfo{},
		Interfaces: []string{},
	}

	if err := p.analyzeFields(typeSpec, info); err != nil {
		return nil, err
	}

	p.analyzeMethods(info)
	p.buildHierarchy(typeSpec, info)

	return info, nil
}

func (p *ASTParser) analyzeFields(typeSpec *ast.TypeSpec, info *ClassInfo) error {
	structType, ok := typeSpec.Type.(*ast.StructType)
	if !ok || structType.Fields == nil {
		return nil
	}

	for _, field := range structType.Fields.List {
		fieldType := p.typeToString(field.Type)

		if field.Names == nil && fieldType != "" {
			// Это встроенная структура (anonim field)
			info.Superclass = fieldType
			continue
		}

		// simple field
		for _, name := range field.Names {
			fieldInfo := FieldInfo{
				Name: name.Name,
				Type: fieldType,
			}
			if field.Tag != nil {
				fieldInfo.Tag = strings.Trim(field.Tag.Value, "`")
			}
			info.Fields = append(info.Fields, fieldInfo)
		}
	}
	return nil
}

func (p *ASTParser) analyzeMethods(info *ClassInfo) {
	for _, file := range p.pkg.Files {
		for _, decl := range file.Decls {
			if funcDecl, ok := decl.(*ast.FuncDecl); ok && funcDecl.Recv != nil {

				if p.isMethodForStruct(funcDecl, info.ClassName) {
					methodInfo := p.parseMethod(funcDecl)
					info.Methods = append(info.Methods, methodInfo)
				}
			}
		}
	}
}

func (p *ASTParser) isMethodForStruct(funcDecl *ast.FuncDecl, structName string) bool {
	if funcDecl.Recv == nil || len(funcDecl.Recv.List) == 0 {
		return false
	}

	recvType := funcDecl.Recv.List[0].Type

	if starExpr, ok := recvType.(*ast.StarExpr); ok {
		recvType = starExpr.X
	}

	if ident, ok := recvType.(*ast.Ident); ok {
		return ident.Name == structName
	}

	return false
}

func (p *ASTParser) parseMethod(funcDecl *ast.FuncDecl) MethodInfo {
	methodInfo := MethodInfo{
		Name:       funcDecl.Name.Name,
		Params:     []string{},
		ReturnType: "",
	}

	if funcDecl.Type.Params != nil {
		for _, param := range funcDecl.Type.Params.List {
			paramType := p.typeToString(param.Type)
			methodInfo.Params = append(methodInfo.Params, paramType)
		}
	}

	if funcDecl.Type.Results != nil {
		if len(funcDecl.Type.Results.List) > 0 {
			methodInfo.ReturnType = p.typeToString(funcDecl.Type.Results.List[0].Type)
		}
	}

	return methodInfo
}

func (p *ASTParser) buildHierarchy(typeSpec *ast.TypeSpec, info *ClassInfo) {
	hierarchy := make(map[string]interface{})

	p.buildHierarchyRecursive(typeSpec, hierarchy)

	info.Hierarchy = hierarchy
}

func (p *ASTParser) buildHierarchyRecursive(typeSpec *ast.TypeSpec, hierarchy map[string]interface{}) {
	structType, ok := typeSpec.Type.(*ast.StructType)
	if !ok || structType.Fields == nil {
		return
	}

	for _, field := range structType.Fields.List {
		if field.Names == nil && field.Type != nil {

			embeddedType := p.typeToString(field.Type)

			embeddedSpec, err := p.FindStruct(embeddedType)
			if err == nil {

				childHierarchy := make(map[string]interface{})
				hierarchy[embeddedType] = childHierarchy

				p.buildHierarchyRecursive(embeddedSpec, childHierarchy)
			} else {

				hierarchy[embeddedType] = make(map[string]interface{})
			}
		}
	}
}

func (p *ASTParser) typeToString(expr ast.Expr) string {
	switch t := expr.(type) {
	case *ast.Ident:
		return t.Name
	case *ast.StarExpr:
		return "*" + p.typeToString(t.X)
	case *ast.ArrayType:
		if t.Len != nil {
			return fmt.Sprintf("[%s]%s", p.exprToString(t.Len), p.typeToString(t.Elt))
		}
		return "[]" + p.typeToString(t.Elt)
	case *ast.SelectorExpr:
		return p.typeToString(t.X) + "." + t.Sel.Name
	case *ast.MapType:
		return fmt.Sprintf("map[%s]%s", p.typeToString(t.Key), p.typeToString(t.Value))
	case *ast.ChanType:
		dir := ""
		switch t.Dir {
		case ast.SEND:
			dir = "chan<- "
		case ast.RECV:
			dir = "<-chan "
		default:
			dir = "chan "
		}
		return dir + p.typeToString(t.Value)
	case *ast.FuncType:
		return "func" + p.funcTypeToString(t)
	case *ast.InterfaceType:
		return "interface{}"
	case *ast.StructType:
		return "struct{}"
	case *ast.Ellipsis:
		return "..." + p.typeToString(t.Elt)
	default:
		return fmt.Sprintf("%T", expr)
	}
}

func (p *ASTParser) exprToString(expr ast.Expr) string {
	switch t := expr.(type) {
	case *ast.BasicLit:
		return t.Value
	case *ast.Ident:
		return t.Name
	default:
		return p.typeToString(expr)
	}
}

func (p *ASTParser) funcTypeToString(funcType *ast.FuncType) string {
	var params, results []string

	if funcType.Params != nil {
		for _, field := range funcType.Params.List {
			paramType := p.typeToString(field.Type)
			params = append(params, paramType)
		}
	}

	if funcType.Results != nil {
		for _, field := range funcType.Results.List {
			resultType := p.typeToString(field.Type)
			results = append(results, resultType)
		}
	}

	result := "(" + strings.Join(params, ", ") + ")"

	if len(results) > 0 {
		if len(results) == 1 {
			result += " " + results[0]
		} else {
			result += " (" + strings.Join(results, ", ") + ")"
		}
	}

	return result
}
