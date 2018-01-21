package gogiven

import (
	"go/ast"
	"go/token"
)

type IdentVisitor struct {
	ast.Visitor
	fset          *token.FileSet
	fileOffsetPos int
}

func (visitor *IdentVisitor) Visit(node ast.Node) ast.Visitor {
	if visitor.fileOffsetPos != -1 {
		return visitor
	}
	switch node.(type) {
	case *ast.Ident:
		ident := node.(*ast.Ident)
		if ident.Name == "Given" || ident.Name == "When" {
			visitor.fileOffsetPos = visitor.fset.Position(node.Pos()).Offset
		}
	}
	return visitor
}
