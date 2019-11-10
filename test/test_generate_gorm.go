package test

import "github.com/jinzhu/gorm"

var db *gorm.DB

//go:generate go-easy generate gorm --client db --type GromTest
type GromTest struct {
	gorm.Model
	AAA int64
}
