package helper

import (
	"testing"
)

func TestContains(t *testing.T) {
	var list = []string{"one", "two", "three", "four"}
	var key = "one"

	result := Contains(key, list)
	if result != true {
		t.Log("error: should be true but got ", result)
		t.Fail()
	}
}
