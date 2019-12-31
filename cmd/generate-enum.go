package cmd

import (
	"bytes"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/siskinc/go-easy/common/generate"
	"github.com/siskinc/go-easy/common/str"
	"github.com/spf13/cobra"
	"go/format"
	"html/template"
	"io/ioutil"
	"path/filepath"
	"strings"
)

func generateErrorCode(goPackage, typeName string, enumMap map[string]generate.EnumerationInfo) {
	enumInfo, ok := enumMap[typeName]
	if !ok {
		logrus.Fatalf("not found type name %s", typeName)
	}
	parseInfo := map[string]interface{}{
		"package":    goPackage,
		"error_type": typeName,
		"error_info": enumInfo.EnumDoc,
	}
	tmpl, err := template.New("").Parse(generate.ErrorCodeFormat)
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
	baseName := fmt.Sprintf("%s_generate_enum.go", str.ToSnakeCase(typeName))
	outputName := filepath.Join(".", strings.ToLower(baseName))
	err = ioutil.WriteFile(outputName, src, 0644)
	if err != nil {
		logrus.Fatalf(" write to file is err: %s", err)
	}
}

var generateErrorCodeCmd = &cobra.Command{
	Use:   "enum",
	Short: "generate code of error code and error information, generate file <dir>/<type>_generate_error_code.go .",
	Run: func(cmd *cobra.Command, args []string) {
		gofile := generate.GetGoFile()
		gopkg := generate.GetGoPackage()
		enumMap, err := generate.ParseEnum(gofile, nil)
		if nil != err {
			logrus.Fatalln("ParseEnum is err:", err)
		}
		contentByte, err := ioutil.ReadFile(gofile)
		if nil != err {
			logrus.Fatalf("Read File %s is err", gofile)
		}
		fileContent = string(contentByte)
		for _, typeName := range typeNameList {
			generateErrorCode(gopkg, typeName, enumMap)
		}
	},
}

func init() {
	generateErrorCodeCmd.Flags().StringArrayVar(&typeNameList, "type", []string{}, "type of need generate code")
	err := generateErrorCodeCmd.MarkFlagRequired("type")
	if nil != err {
		panic(err)
	}
}
