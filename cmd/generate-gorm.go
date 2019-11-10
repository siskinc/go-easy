package cmd

import (
	"bytes"
	"fmt"
	"github.com/fatih/structtag"
	"github.com/sirupsen/logrus"
	"github.com/siskinc/go-easy/common/generate"
	"github.com/siskinc/go-easy/common/str"
	"github.com/spf13/cobra"
	"go/ast"
	"go/format"
	"html/template"
	"io/ioutil"
	"path/filepath"
	"strings"
)

type GormFiledInfo struct {
	FiledName  string
	ColumnName string
}

func generateGorm(goPackage, typeName string, structFieldListMap map[string][]*ast.Field) {
	var fieldList []*ast.Field
	var ok bool
	if fieldList, ok = structFieldListMap[typeName]; !ok {
		logrus.Fatalln("not fount type: ", typeName)
	}
	parseInfo := map[string]interface{}{
		"package":            goPackage,
		"struct":             typeName,
		"client":             clientNameForCollection,
		"query_builder_name": fmt.Sprintf("%sQueryBuilder", typeName),
	}
	gormFieldInfoList := make([]GormFiledInfo, len(fieldList))

	for i, fieldObj := range fieldList {
		fmt.Printf("filedObj, %+v\n", fieldObj)
		//ast.Inspect(fieldObj, func(node ast.Node) bool {
		//	fmt.Printf("%+v %+v\n", node, reflect.TypeOf(node))
		//	switch x := node.(type) {
		//	case *ast.Ident:
		//		fmt.Println(x.Name, x.NamePos, x.Obj, x.String())
		//	}
		//	return true
		//})
		fieldName := fieldObj.Names[0].Name
		gormFieldInfoList[i].FiledName = fieldName
		gormFieldInfoList[i].ColumnName = str.LatinCharFirst(fieldName)
		if nil == fieldObj.Tag {
			continue
		}
		tags, err := structtag.Parse(strings.Trim(fieldObj.Tag.Value, "`"))
		if nil != err {
			logrus.Fatalf(" Parse tag is err: %s, typeName: %s, filedName: %s", err, typeName, fieldName)
		}
		gromTag, _ := tags.Get("grom")
		if nil == gromTag {
			continue
		}
		gormInfoList := strings.Split(gromTag.Name, ";")
		for _, info := range gormInfoList {
			info = strings.Trim(info, " ")
			if strings.HasPrefix(info, "column") {
				gormFieldInfoList[i].ColumnName = info[len("column"):]
			}
		}
	}
	parseInfo["fields"] = gormFieldInfoList

	tmpl, err := template.New("").Parse(generate.GormParseBasic)
	if nil != err {
		logrus.Fatalf("New %s struct template is err: %s", typeName, err)
	}
	buff := bytes.NewBufferString("")
	err = tmpl.Execute(buff, parseInfo)
	if nil != err {
		logrus.Fatalf("template Execute is err: %s, typeName is %s", err, typeName)
	}
	// 格式化
	src, err := format.Source(buff.Bytes())
	if nil != err {
		logrus.Fatalf("format Source is err: %s, typeName is %s", err, typeName)
	}
	baseName := fmt.Sprintf("%s_generate_gorm.go", str.ToSnakeCase(typeName))
	outputName := filepath.Join(".", strings.ToLower(baseName))
	err = ioutil.WriteFile(outputName, src, 0644)
	if err != nil {
		logrus.Fatalf(" write to file is err: %s", err)
	}
}

var generateGormCmd = &cobra.Command{
	Use:   "gorm",
	Short: "generate code of operate grom, generate file <dir>/<type>_generate_gorm.go .",
	Run: func(cmd *cobra.Command, args []string) {
		gofile := generate.GetGoFile()
		gopkg := generate.GetGoPackage()
		parser := generate.NewParser()
		parser.ParseFiles([]string{gofile})
		structInfo, err := generate.ParseStruct(gofile, nil)
		//fmt.Printf("%+v", structFieldListMap)
		if nil != err {
			logrus.Fatalln(" ParseStruct is err:", err)
		}
		structFieldListMap := structInfo.FieldListMap
		fmt.Printf("%+v\n", structFieldListMap)
		contentByte, err := ioutil.ReadFile(gofile)
		if nil != err {
			logrus.Fatalf("Read File %s is err", gofile)
		}
		fileContent = string(contentByte)
		for _, typeName := range typeNameList {
			generateGorm(gopkg, typeName, structFieldListMap)
		}
	},
}

func init() {
	generateGormCmd.Flags().StringVar(&clientNameForCollection, "client", "", "db of gorm")
	err := generateGormCmd.MarkFlagRequired("client")
	if nil != err {
		panic(err)
	}
	generateGormCmd.Flags().StringArrayVar(&typeNameList, "type", []string{}, "type of need generate code")
	err = generateGormCmd.MarkFlagRequired("type")
	if nil != err {
		panic(err)
	}
}
