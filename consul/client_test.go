package consul

import (
	"strconv"
	"testing"
)

func TestLoadCustomValueForString(t *testing.T) {
	var strVar string
	input := "test-data"
	data := []byte(input)
	loadCustomValue(&strVar, data, ConfigTypeString)
	if strVar != "test-data" {
		t.Errorf("unable to parse %s to string", input)
	}
}

func TestLoadCustomValueForBool(t *testing.T) {
	var boolVar bool
	input := []string{
		"true",
		"false",
		"True",
		"False",
		"TRUE",
		"FALSE",
	}

	for i, in := range input {
		data := []byte(in)
		loadCustomValue(&boolVar, data, ConfigTypeBool)
		if i%2 == 0 && !boolVar {
			t.Errorf("unable to parse %s to boolean", in)
		}
	}

}

func TestLoadCustomValueForInt(t *testing.T) {
	var intVar int
	input := "235"
	data := []byte(input)
	loadCustomValue(&intVar, data, ConfigTypeInt)
	if in, _ := strconv.Atoi(input); in != intVar {
		t.Errorf("unable to parse %s to int", input)
	}
}
