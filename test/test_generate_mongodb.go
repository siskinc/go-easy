package test

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var c = &mongo.Client{}
var client = c.Database("aaa").Collection("ddd")

//go:generate ./../go-easy generate mongodb --client client --type MongoDBTest
//@def soft_delete SoftDelete
type MongoDBTest struct {
	ID         primitive.ObjectID
	AA         string
	SoftDelete bool
}
