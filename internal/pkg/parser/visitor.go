package parser

import (
	"go/ast"
)

type visitor struct {
	allowedStructs map[string]bool
	structs        []structDecl
}

func newVisitor(allowedStructs []string) *visitor {
	v := &visitor{
		allowedStructs: map[string]bool{},
	}

	for _, s := range allowedStructs {
		v.allowedStructs[s] = true
	}
	return v
}

func (v *visitor) checkStructName(name string) bool {
	if len(v.allowedStructs) == 0 {
		return true
	}

	_, ok := v.allowedStructs[name]
	return ok
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
			name:   structName,
			fields: fields,
		})
	}
	return v
}
