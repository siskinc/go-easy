package generate

import (
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"strings"
)

type EnumerationInfo struct {
	EnumDoc map[string]string
}

func ParseEnum(filename string, src []byte) (enumMap map[string]EnumerationInfo, err error) {
	enumMap = map[string]EnumerationInfo{}
	if src == nil {
		src, err = ioutil.ReadFile(filename)
		if err != nil {
			panic(err)
			return
		}
	}

	file, err := parser.ParseFile(token.NewFileSet(), filename, src, parser.ParseComments)
	if err != nil {
		return
	}

	collectEnumInfo := func(node ast.Node) bool {
		decl, ok := node.(*ast.GenDecl)
		if !ok || decl.Tok != token.CONST {
			return true
		}
		type_ := ""
		for _, spec := range decl.Specs {
			vspec := spec.(*ast.ValueSpec)
			if vspec.Type == nil && len(vspec.Values) > 0 {
				// 排除 v = 1 这种结构
				type_ = ""
				continue
			}
			//如果Type不为空，则确认typ
			//fmt.Println("Type", vspec.Type)
			//fmt.Println(vspec.Doc.Text())
			if vspec.Type != nil {
				ident, ok := vspec.Type.(*ast.Ident)
				if !ok {
					continue
				}
				type_ = ident.Name
			}
			info, ok := enumMap[type_]
			if !ok {
				info = EnumerationInfo{EnumDoc: map[string]string{}}
				enumMap[type_] = info
			}
			for _, n := range vspec.Names {
				info.EnumDoc[n.Name] = strings.Trim(vspec.Doc.Text(), " \n\t\r")
			}
		}
		return false
	}

	ast.Inspect(file, collectEnumInfo)

	return
}
