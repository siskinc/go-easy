package generate

import (
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

func collectStructFields(ts *ast.TypeSpec) []*ast.Field {
	//fmt.Println("collectStructFields")
	x := ts.Type.(*ast.StructType)
	//fmt.Println("len(x.Fields.List)", len(x.Fields.List))
	filedList := make([]*ast.Field, len(x.Fields.List))
	for index, field := range x.Fields.List {
		filedList[index] = field
	}
	return filedList
}

func ParseStruct(filename string, src []byte) (structFieldListMap map[string][]*ast.Field, structDocumentListMap map[string][]string, err error) {
	structFieldListMap = make(map[string][]*ast.Field)
	structDocumentListMap = make(map[string][]string)
	if src == nil {
		src, err = ioutil.ReadFile(filename)
		if err != nil {
			return
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
			filedList := collectStructFields(x)
			structFieldListMap[structName] = filedList
		case *ast.GenDecl:
			structName, documents := collectStructDocument(x)
			structDocumentListMap[structName] = documents
		default:
		}
		return true
	}

	ast.Inspect(file, collectStructInfo)

	return
}
