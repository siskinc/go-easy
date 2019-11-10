package generate

import (
	"testing"
)

func TestParseStruct(t *testing.T) {
	ParseStruct("../../test/test_generate_gorm.go", nil)
}
