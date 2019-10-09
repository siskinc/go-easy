package generate

import (
	"go/ast"
)

func GetFieldType(fileContent string, field *ast.Field) string {
	typeExpr := field.Type
	start := typeExpr.Pos() - 1
	end := typeExpr.End() - 1

	// grab it in source
	typeInSource := fileContent[start:end]

	return typeInSource
}
