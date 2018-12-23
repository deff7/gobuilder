package parser

import (
	"go/ast"
)

type visitor struct {
	lastIdent string
	structs   []structDecl
}

func newVisitor() *visitor {
	return &visitor{}
}

func (v *visitor) Visit(node ast.Node) ast.Visitor {
	if typeSpec, ok := node.(*ast.TypeSpec); ok {
		s, ok := typeSpec.Type.(*ast.StructType)
		if !ok {
			return v
		}

		fields := []field{}
		for _, list := range s.Fields.List {
			var typeName string
			// we consider only simple types (not anonymous structs)
			t, ok := list.Type.(*ast.Ident)
			if !ok {
				continue
			}
			typeName = t.Name

			for _, f := range list.Names {
				fields = append(fields, field{
					name:     f.Name,
					typeName: typeName,
				})
			}
		}

		v.structs = append(v.structs, structDecl{
			name:   typeSpec.Name.Name,
			fields: fields,
		})
	}
	return v
}
