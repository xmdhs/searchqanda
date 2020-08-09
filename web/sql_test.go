package web

import (
	"fmt"
	"testing"
)

func TestSearch(t *testing.T) {
	r, err := search("启动器", "0")
	if err != nil {
		t.Error(err)
	}
	fmt.Println(r)
}
