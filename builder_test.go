package main

import (
	"reflect"
	"testing"
)

type builderTest struct {
	name     string
	keyPath  []string
	value    interface{}
	expected Struct
	err      error
}

var builderTests = []builderTest{
	// [foo] => "bar"
	{
		name:    "kv-string",
		keyPath: []string{"foo"},
		value:   "bar",
		expected: Struct{
			"foo": "bar",
		},
	},
	// [foo] => 123
	{
		name:    "kv-int",
		keyPath: []string{"foo"},
		value:   123,
		expected: Struct{
			"foo": 123,
		},
	},
	// [foo] => true
	{
		name:    "kv-bool",
		keyPath: []string{"foo"},
		value:   true,
		expected: Struct{
			"foo": true,
		},
	},
	// [foo, bar, baz] => "quux"
	{
		name:    "complex-path",
		keyPath: []string{"foo", "bar", "baz"},
		value:   "quux",
		expected: Struct{
			"foo": Struct{
				"bar": Struct{
					"baz": "quux",
				},
			},
		},
	},
	// [foo] => [bar]
	// Slice values
	{
		name:    "kv-array",
		keyPath: []string{"foo"},
		value:   []string{"bar"},
		expected: Struct{
			"foo": []string{"bar"},
		},
	},
	// [foo] => [bar: [baz: "quux"]]
	// Map values
	{
		name:    "kv-int",
		keyPath: []string{"foo"},
		value: map[string]interface{}{
			"bar": map[string]string{
				"baz": "quux",
			},
		},
		expected: Struct{
			"foo": map[string]interface{}{
				"bar": map[string]string{
					"baz": "quux",
				},
			},
		},
	},
}

func TestBuilder(t *testing.T) {
	for _, test := range builderTests {
		test := test // range capture
		t.Run(test.name, func(t *testing.T) {
			obj := Struct{}
			err := obj.Set(test.keyPath, test.value)
			if err != test.err {
				t.Fatalf("expected build to error with %v, but got %v", test.err, err)
			}
			if !reflect.DeepEqual(obj, test.expected) {
				t.Fatalf("expected build to create structure %+v, but got %+v", test.expected, obj)
			}
		})
	}
}
