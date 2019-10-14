package dbgenerate

import (
	"fmt"
	"testing"
)

func TestCreateStruct(t *testing.T) {
	err := GenCode("dd:2222ddD_@tcp(47.94.168.30:3316)/ddc")
	fmt.Println(err)
}
