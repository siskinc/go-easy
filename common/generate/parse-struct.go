package generate

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"strings"
)

func GetNodeStructName(node ast.Node) string {
	ts, ok := node.(*ast.TypeSpec)
	if !ok || ts.Type == nil {
		return ""
	}

	// get struct name
	structName := ts.Name.Name

	return structName
}

func collectStructDocument(x *ast.GenDecl) (structName string, documents []string) {
	//fmt.Println("collectStructDocument")
	if x.Tok != token.TYPE {
		return
	}
	sts, ok := x.Specs[0].(*ast.TypeSpec)
	if !ok {
		return
	}
	structName = sts.Name.Name
	documents = []string{}
	for _, doc := range x.Doc.List {
		text := doc.Text
		if strings.HasPrefix(text, "//") {
			text = text[2:]
			text = strings.Trim(text, "")
			documents = append(documents, text)
		} else if strings.HasPrefix(text, "/*") && strings.HasSuffix(text, "*/") {
			text = text[2 : len(text)-2]
			textList := strings.Split(text, "\n")
			for _, t := range textList {
				temp := strings.Trim(t, "")
				documents = append(documents, temp)
			}
		}
	}
	return
}

func collectStructFields(ts *ast.StructType) []*ast.Field {
	//fmt.Println("collectStructFields")
	//x := ts.Type.(*ast.StructType)
	x := ts
	fmt.Println("len(x.Fields.List)", x.Fields.List)
	var filedList []*ast.Field
	for _, field := range x.Fields.List {
		if 0 == len(field.Names) {
			continue
		}
		filedList = append(filedList, field)
	}
	return filedList
}

type StructInfo struct {
	FieldListMap    map[string][]*ast.Field
	DocumentListMap map[string][]string
}

func ParseStruct(filename string, src []byte) (structInfo *StructInfo, err error) {
	structInfo = &StructInfo{}
	structInfo.FieldListMap = make(map[string][]*ast.Field)
	structInfo.DocumentListMap = make(map[string][]string)
	if src == nil {
		src, err = ioutil.ReadFile(filename)
		if err != nil {
			panic(err)
		}
	}
	file, err := parser.ParseFile(token.NewFileSet(), filename, src, parser.ParseComments)
	if err != nil {
		return
	}
	collectStructInfo := func(node ast.Node) bool {
		switch x := node.(type) {
		case *ast.TypeSpec:
			structName := x.Name.Name
			s, ok := x.Type.(*ast.StructType)
			if ok {
				filedList := collectStructFields(s)
				structInfo.FieldListMap[structName] = filedList
			}
			fmt.Printf("aaa : %+v\n", structInfo.FieldListMap[structName])
		case *ast.GenDecl:
			structName, documents := collectStructDocument(x)
			structInfo.DocumentListMap[structName] = documents
		default:
		}
		return true
	}

	ast.Inspect(file, collectStructInfo)

	return
}
