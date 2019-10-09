package generate

import (
	"fmt"
	"testing"
)

func TestParseEnum(t *testing.T) {
	enumMap, err := ParseEnum("../../test/test_generate_enum.go", nil)
	if err != nil {
		panic(err)
	}
	fmt.Println("enumMap", enumMap)
}
