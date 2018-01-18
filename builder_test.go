package main

import "testing"

func TestBuildSimpleJSON(t *testing.T) {
	var obj Struct
	var err error

	obj = Struct{}

	err = obj.Set("hello", "world")
	if err != nil {
		t.Errorf("Error: %v", err)
		return
	}

	err = obj.Set("foo.bar", "baz")
	if err != nil {
		t.Errorf("Error: %v", err)
		return
	}

	actualJSON := obj.String()
	expectedJSON := `{"foo":{"bar":"baz"},"hello":"world"}`
	if actualJSON != expectedJSON {
		t.Errorf("Error: Expected %v, actual %v", expectedJSON, actualJSON)
	}
}
