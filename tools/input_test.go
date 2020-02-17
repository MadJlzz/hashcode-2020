package tools

import "testing"

func TestReadInput(t *testing.T) {
	data := ReadInput("../test/a_example.in")
	if len(data) != 2 {
		t.Errorf("map should contain two entry because there is two line in the file. (test/a_example.in)")
	}
	for _, list := range data {
		if len(list) <= 0 {
			t.Errorf("list values from keys should not be empty.")
		}
	}
}
