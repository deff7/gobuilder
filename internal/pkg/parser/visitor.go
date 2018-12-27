package parser

import (
	"go/ast"
)

type visitor struct {
	unexported     bool
	allowedStructs map[string]bool
	structs        []StructDecl
}

func newVisitor(allowedStructs []string, unexported bool) *visitor {
	v := &visitor{
		allowedStructs: map[string]bool{},
		unexported:     unexported,
	}

	for _, s := range allowedStructs {
		v.allowedStructs[s] = true
	}
	return v
}

func (v *visitor) checkStructName(name string) bool {
	if !v.unexported {
		if !ast.IsExported(name) {
			return false
		}
	}

	if len(v.allowedStructs) == 0 {
		return true
	}

	_, ok := v.allowedStructs[name]
	return ok
}

func (v *visitor) checkFieldName(name string) bool {
	if !v.unexported {
		if !ast.IsExported(name) {
			return false
		}
	}

	return true
}

func collectTypeName(expr ast.Expr) string {
	switch t := expr.(type) {
	case *ast.Ident:
		return t.Name
	case *ast.StarExpr:
		return "*" + collectTypeName(t.X)
	case *ast.ArrayType:
		if t.Len == nil {
			return "[]" + collectTypeName(t.Elt)
		}

		if v, ok := t.Len.(*ast.BasicLit); ok {
			return "[" + v.Value + "]" + collectTypeName(t.Elt)
		}
	case *ast.SelectorExpr:
		return t.X.(*ast.Ident).Name + "." + t.Sel.Name
	}
	return ""
}

func (v *visitor) Visit(node ast.Node) ast.Visitor {
	if typeSpec, ok := node.(*ast.TypeSpec); ok {
		s, ok := typeSpec.Type.(*ast.StructType)
		if !ok {
			return v
		}

		structName := typeSpec.Name.Name
		if !v.checkStructName(structName) {
			return v
		}

		fields := []Field{}
		for _, list := range s.Fields.List {
			typeName := collectTypeName(list.Type)
			if typeName == "" {
				continue
			}

			for _, f := range list.Names {
				if !v.checkFieldName(f.Name) {
					continue
				}
				fields = append(fields, Field{
					Name:     f.Name,
					TypeName: typeName,
				})
			}
		}

		v.structs = append(v.structs, StructDecl{
			Name:   structName,
			Fields: fields,
		})
	}
	return v
}
