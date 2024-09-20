package utils

import (
	"fmt"
	"testing"
)

func TestPathExists(t *testing.T) {
	exists, err := PathExists("../config")
	fmt.Println(exists, err)
}
