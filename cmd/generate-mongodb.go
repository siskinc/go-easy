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
	typeNameList            []string
	clientNameForCollection string
)

type UniqueIndexInfo struct {
	Parameters               string
	FieldNameList            []string
	FieldBsonNameList        []string
	FieldVariableNameList    []string
	FieldTypeAndVariableName map[string][]string
}

func (u *UniqueIndexInfo) GenParameters() {
	for typeName, variableNameList := range u.FieldTypeAndVariableName {
		for _, variableName := range variableNameList {
			if "" == u.Parameters {
				u.Parameters = variableName
			} else {
				u.Parameters = fmt.Sprintf("%s, %s", u.Parameters, variableName)
			}
		}
		u.Parameters = fmt.Sprintf("%s %s", u.Parameters, typeName)
	}
}

func generateMongoDB(goPackage, typeName string, structFieldListMap map[string][]*ast.Field, structDocumentListMap map[string][]string) {
	uniqueIndexMap := make(map[string]UniqueIndexInfo)
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
		"client":  clientNameForCollection,
		"id":      "ID",
	}

	parseStr := generate.MongoDBParseBasic
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
					field, ok := fieldMap[fieldNameFromCommand]
					if !ok {
						logrus.Fatalf(" not found soft delete field %s", fieldNameFromCommand)
					}
					bsonName, err := generate.GetBsonName(fieldNameFromCommand, field)
					if nil != err {
						logrus.Fatalf(" Parse tag is err: %s, structName: %s, filedName: %s, command: %s", err,
							typeName, fieldNameFromCommand, command)
					}
					if 3 <= len(commandList) {
						if fieldName, ok := parseInfo[command]; !ok {
							if generate.SoftDelete == command {
								parseStr += generate.MongoDBParseSoftDelete
							}
						} else {
							logrus.Fatalf(" soft delete have been double declared, %s, %s, %s",
								command, fieldName, commandList[2])
						}
						parseInfo[command] = commandList[2]
						parseInfo[command+"_bson_name"] = bsonName
					}
				case generate.UniqueIndex:
					uniqueFieldList := commandList[2:]

					uniqueInfo := UniqueIndexInfo{
						FieldNameList:            make([]string, len(uniqueFieldList)),
						FieldBsonNameList:        make([]string, len(uniqueFieldList)),
						FieldVariableNameList:    make([]string, len(uniqueFieldList)),
						FieldTypeAndVariableName: make(map[string][]string),
					}
					uniqueIndexName := ""

					// check unique field is exist
					// get bson value of unique field
					for index, uniqueField := range uniqueFieldList {
						if field, ok := fieldMap[uniqueField]; !ok {
							logrus.Fatalf("generate unique field is failed, because the"+
								"field %s not found", uniqueField)
						} else {
							bsonName, err := generate.GetBsonName(uniqueField, field)
							if nil != err {
								logrus.Fatalf("can't find %s filed's bson value", uniqueField)
							}
							uniqueInfo.FieldNameList[index] = uniqueField
							uniqueInfo.FieldBsonNameList[index] = bsonName
							variableName := str.LatinCharFirst(uniqueField)
							variableType := generate.GetFieldType(fileContent, field)
							uniqueInfo.FieldTypeAndVariableName[variableType] =
								append(uniqueInfo.FieldTypeAndVariableName[variableType], variableName)
							uniqueInfo.FieldVariableNameList[index] = variableName
						}

						if index == 0 {
							uniqueIndexName = uniqueField
						} else {
							uniqueIndexName = fmt.Sprintf("%sAnd%s", uniqueIndexName, uniqueField)
						}
					}
					uniqueInfo.GenParameters()
					uniqueIndexMap[uniqueIndexName] = uniqueInfo
				}
			}
		}
	}

	if 0 < len(uniqueIndexMap) {
		parseInfo["unique_index_map"] = uniqueIndexMap
		parseStr += generate.MongoDBUniqueIndexParse
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
		contentByte, err := ioutil.ReadFile(gofile)
		if nil != err {
			logrus.Fatalf("Read File %s is err", gofile)
		}
		fileContent = string(contentByte)
		for _, typeName := range typeNameList {
			generateMongoDB(gopkg, typeName, structFieldListMap, structDocumentListMap)
		}
	},
}

func init() {
	generateMongoDBCmd.Flags().StringVar(&clientNameForCollection, "client", "", "client of mongo-go-driver")
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
