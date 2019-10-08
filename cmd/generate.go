/*
Copyright © 2019 daryl <susecjh@gmail.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
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

var (
	clientName   string
	typeNameList []string
)

// generateCmd represents the generate command
var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "go generate tool",
	Long:  `A go generate tool for generate code.`,
	//Run: func(cmd *cobra.Command, args []string) {
	//	fmt.Println("generate called")
	//},
}

var generateMongoDBCmd = &cobra.Command{
	Use:   "mongodb",
	Short: "generate code of operate mongodb by mongo-go-driver, generate file <dir>/<type>_generate_mongodb.go .",
	Run: func(cmd *cobra.Command, args []string) {
		gofile := generate.GetGoFile()
		gopkg := generate.GetGoPackage()
		structFieldListMap, structDocumentListMap, err := generate.ParseStruct(gofile, nil)
		if nil != err {
			logrus.Fatalln(" ParseStruct is err:", err)
		}
		for _, typeName := range typeNameList {
			generateMongoDB(gopkg, typeName, structFieldListMap, structDocumentListMap)
		}
	},
}

func generateMongoDB(goPackage, typeName string, structFieldListMap map[string][]*ast.Field, structDocumentListMap map[string][]string) {
	var documentList []string
	var fieldList []*ast.Field
	var ok bool
	if documentList, ok = structDocumentListMap[typeName]; !ok {
		logrus.Fatalln(" not fount type: ", typeName)
	}
	if fieldList, ok = structFieldListMap[typeName]; !ok {
		logrus.Fatalln(" not fount type: ", typeName)
	}
	parseInfo := map[string]interface{}{
		"package": goPackage,
		"struct":  typeName,
		"client":  clientName,
		"id":      "ID",
	}

	parseStr := generate.ParseBasic
	fieldMap := make(map[string]*ast.Field)

	for _, fieldObj := range fieldList {
		fieldName := fieldObj.Names[0].Name
		fieldMap[fieldName] = fieldObj
		if nil != fieldObj.Tag {
			tags, err := structtag.Parse(strings.Trim(fieldObj.Tag.Value, "`"))
			if nil != err {
				logrus.Fatalf(" Parse tag is err: %s, typeName: %s, filedName: %s", err, typeName, fieldName)
			}
			bsonTag, _ := tags.Get("bson")
			if nil != bsonTag {
				if "_id" == bsonTag.Name {
					parseInfo["id"] = fieldName
				}
			}
		}
	}

	for _, document := range documentList {
		if strings.HasPrefix(document, "@def") {
			commandList := strings.Split(document, " ")
			if 2 <= len(commandList) {
				fieldNameFromCommand := commandList[2]
				command := commandList[1]
				switch command {
				case generate.SoftDelete, generate.SoftDeleteAt, generate.UpdateAt, generate.CreateAt:
					softDeleteField, ok := fieldMap[fieldNameFromCommand]
					if !ok {
						logrus.Fatalf(" not found soft delete field %s", fieldNameFromCommand)
					}
					tagValue := ""
					if nil != softDeleteField.Tag {
						tagValue = strings.Trim(softDeleteField.Tag.Value, "`")
					}
					tags, err := structtag.Parse(tagValue)
					if nil != err {
						logrus.Fatalf(" Parse tag is err: %s, structName: %s, filedName: %s, command: %s", err,
							typeName, fieldNameFromCommand, command)
					}
					bsonTag, _ := tags.Get("bson")
					bsonName := ""
					if nil != bsonTag {
						bsonName = bsonTag.Name
					} else {
						bsonName = str.ToSnakeCase(fieldNameFromCommand)
					}
					if 3 <= len(commandList) {
						if fieldName, ok := parseInfo[command]; !ok {
							if generate.SoftDelete == command {
								parseStr += generate.ParseSoftDelete
							}
						} else {
							logrus.Fatalf(" soft delete have been double declared, %s, %s, %s",
								command, fieldName, commandList[2])
						}
						parseInfo[command] = commandList[2]
						parseInfo[command+"_bson_name"] = bsonName
					}
				case generate.UniqueIndex:

				}
			}
		}
	}

	tmpl, err := template.New("").Parse(parseStr)
	if nil != err {
		logrus.Fatalf(" New %s struct template is err: %s", typeName, err)
	}
	buff := bytes.NewBufferString("")
	err = tmpl.Execute(buff, parseInfo)
	if nil != err {
		logrus.Fatalf(" template Execute is err: %s, typeName is %s", err, typeName)
	}
	// 格式化
	src, err := format.Source(buff.Bytes())
	if nil != err {
		logrus.Fatalf(" format Source is err: %s, typeName is %s", err, typeName)
	}
	baseName := fmt.Sprintf("%s_generate_mgorm.go", str.ToSnakeCase(typeName))
	outputName := filepath.Join(".", strings.ToLower(baseName))
	err = ioutil.WriteFile(outputName, src, 0644)
	if err != nil {
		logrus.Fatalf(" write to file is err: %s", err)
	}

}

func init() {
	generateCmd.AddCommand(generateMongoDBCmd)
	rootCmd.AddCommand(generateCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// generateCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// generateCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	generateMongoDBCmd.Flags().StringVar(&clientName, "client", "", "client of mongo-go-driver")
	err := generateMongoDBCmd.MarkFlagRequired("client")
	if nil != err {
		panic(err)
	}
	generateMongoDBCmd.Flags().StringArrayVar(&typeNameList, "type", []string{}, "type of need generate code")
	err = generateMongoDBCmd.MarkFlagRequired("type")
	if nil != err {
		panic(err)
	}
}
