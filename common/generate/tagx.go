package generate

import (
	"github.com/fatih/structtag"
	"github.com/siskinc/go-easy/common/str"
	"go/ast"
	"strings"
)

func GetBsonName(fieldName string, field *ast.Field) (string, error) {
	tagValue := ""
	if nil != field.Tag {
		tagValue = strings.Trim(field.Tag.Value, "`")
	}
	tags, err := structtag.Parse(tagValue)
	if nil != err {
		return "", err
	}
	bsonTag, _ := tags.Get("bson")
	bsonName := ""
	if nil != bsonTag {
		bsonName = bsonTag.Name
	} else {
		bsonName = str.ToSnakeCase(fieldName)
	}
	return bsonName, nil
}
