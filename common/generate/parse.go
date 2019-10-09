package generate

//var DontEdit = "// Code generated by go-easy generate-mongodb DO NOT EDIT."

var ParseBasic = `
// Code generated by go-easy generate-mongodb DO NOT EDIT.
package {{.package}}

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

var _ = time.Now()
var defaultContext = context.Background()

func (model *{{.struct}}) Save() (interface{}, error) {
	{{ if .create_at }}
	model.{{.create_at}} = time.Now().UTC()
	{{ end }}
	result, err := {{.client}}.InsertOne(defaultContext, model)
	if nil != err {
		return nil, err
	}
	return result.InsertedID, err
}

func (model *{{.struct}}) Delete(filter interface{}) (int64, error) {
	result, err := {{.client}}.DeleteMany(defaultContext, filter)
	if nil != err {
		return 0, err
	}
	return result.DeletedCount, err
}

func (model *{{.struct}}) DeleteByID() (int64, error) {
	filter := bson.M{
		"_id": model.{{.id}},
	}
	result, err := {{.client}}.DeleteOne(defaultContext, filter)
	if nil != err {
		return 0, err
	}
	return result.DeletedCount, err
}

func (model *{{.struct}}) FindByID() (err error) {
	filter := bson.M{
		"_id": model.{{.id}},
		{{ if .soft_delete }}
		"{{.soft_delete_bson_name}}": false,
		{{ end }}
	}

	result := {{.client}}.FindOne(defaultContext, filter)
	err = result.Decode(model)

	return
}

func (model *{{.struct}}) Find(filter interface{}) (err error) {
	{{ if .soft_delete }}
	mQuery := filter.(map[string]interface{})
	mQuery["{{.soft_delete_bson_name}}"] = false
	{{ end }}
	result := {{.client}}.FindOne(defaultContext, filter)
	err = result.Decode(model)
	return
}

func (model *{{.struct}}) FindAll(filter interface{}) (modelList []{{.struct}}, err error) {
	{{ if .soft_delete }}
	mQuery := filter.(map[string]interface{})
	mQuery["{{.soft_delete_bson_name}}"] = false
	{{ end }}
	var cursor *mongo.Cursor
	cursor, err = {{.client}}.Find(defaultContext, filter)
	for nil != cursor && cursor.Next(defaultContext) {
		temp := {{.struct}}{}
		err = cursor.Decode(&temp)
		if nil != err {
			return 
		}
		modelList = append(modelList, temp)
	}
	return
}

func (model *{{.struct}}) FindPage(filter interface{}, iPageSize, iPageIndex int64, SortedStrs ...string) (modelList []{{.struct}}, count int64, err error) {
	{{ if .soft_delete }}
	mQuery := filter.(map[string]interface{})
	mQuery["{{.soft_delete_bson_name}}"] = false
	{{ end }}
	count, err = {{.client}}.CountDocuments(defaultContext, filter)
	if nil != err {
		return
	}
	opt := &options.FindOptions{}
	skip := iPageSize * (iPageIndex - 1)
	opt = opt.SetLimit(iPageSize).SetSkip(skip)
	for _, sortStr := range SortedStrs {
		opt = opt.SetSort(sortStr)
	}
	var cursor *mongo.Cursor
	cursor, err = {{.client}}.Find(defaultContext, filter, opt)
	for nil != cursor && cursor.Next(defaultContext) {
		temp := {{.struct}}{}
		err = cursor.Decode(&temp)
		if nil != err {
			return 
		}
		modelList = append(modelList, temp)
	}
	return
}

func (model *{{.struct}}) UpdateByID() (err error) {
	objectID := model.{{.id}}
	filter := bson.M{
		"_id": objectID,
		{{ if .soft_delete }}
		"{{.soft_delete_bson_name}}": false,
		{{ end }}
	}
	{{ if .update_at }}
	model.{{.update_at}} = time.Now().UTC()
	{{ end }}
	_, err = {{.client}}.UpdateOne(defaultContext, filter, model)
	return
}

func (model *{{.struct}}) Update(filter interface{}) (err error) {
	{{ if .soft_delete }}
	mQuery := filter.(map[string]interface{})
	mQuery["{{.soft_delete_bson_name}}"] = false
	{{ end }}
	{{ if .update_at }}
	model.{{.update_at}} = time.Now().UTC()
	{{ end }}
	_, err = {{.client}}.UpdateMany(defaultContext, filter, model)
	return
}

`

var ParseSoftDelete = `
func (model *{{.struct}}) SoftDeleteByID() (err error) {
	filter := bson.M{
		"_id": model.{{.id}},
		"{{.soft_delete_bson_name}}": false,
	}

	_, err = {{.client}}.CountDocuments(defaultContext, filter)
	if nil != err {
		return
	}

	err = model.Find(filter)
	if nil != err {
		return
	}

	model.{{.soft_delete}} = true

	{{ if .soft_delete_at}}
	model.{{.soft_delete_at}} = time.Now().UTC()
	{{ end }}
	err = model.UpdateByID()
	return
}

`

var UniqueIndexParse = `
{{ range $unique_index_name, $unique_index_info := .unique_index_map }}
func (model *{{$.struct}}) FindBy{{$unique_index_name}}({{ $unique_index_info.Parameters }}) (err error) {
	filter := bson.M{
		"_id": model.{{$.id}},
		{{ if $.soft_delete }}
		"{{ $.soft_delete_bson_name }}": false,
		{{ end }}
		{{ range $i, $bson_name := $unique_index_info.FieldBsonNameList }}
		"{{ $bson_name }}": {{index $unique_index_info.FieldVariableNameList $i}},
		{{ end }}
	}

	result := {{$.client}}.FindOne(defaultContext, filter)
	err = result.Decode(model)

	return
}

{{ end }}
`