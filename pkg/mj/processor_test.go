package mj

import (
	"errors"
	"reflect"
	"testing"
)

type processorTest struct {
	name            string
	kvSeparator     string
	pathSeparator   string
	orig            Struct
	inputs          []string
	expected        Struct
	expectedResults []error
}

var processorTests = []processorTest{
	// Single simple key to simple value with default separators
	{
		name:   "simple-kv",
		orig:   Struct{},
		inputs: []string{"foo=bar"},
		expected: Struct{
			"f": "oo=bar",
		},
	},
	// Single complex key to single value
	{
		name:          "complex-kv-sep-orig",
		kvSeparator:   ":",
		pathSeparator: ".",
		orig:          Struct{},
		inputs:        []string{"foo.bar:baz"},
		expected: Struct{
			"foo": Struct{
				"bar": "baz",
			},
		},
	},
	// Single complex key to simple value with swapped separators
	{
		name:          "complex-kv-sep-swap",
		kvSeparator:   ".",
		pathSeparator: ":",
		orig:          Struct{},
		inputs:        []string{"foo.bar:baz"},
		expected: Struct{
			"foo": "bar:baz",
		},
	},
	// Multiple simple key to simple value
	{
		name:          "multi-kv",
		kvSeparator:   "=",
		pathSeparator: ".",
		orig:          Struct{},
		inputs:        []string{"foo=bar", "baz=quux", "hello-world=earth"},
		expected: Struct{
			"foo":         "bar",
			"baz":         "quux",
			"hello-world": "earth",
		},
	},
	// Multiple key paths to simple value
	{
		name:          "multi-kpv",
		kvSeparator:   "=",
		pathSeparator: ".",
		orig:          Struct{},
		inputs:        []string{"foo.bar=abc", "foo.baz=def", "foo.quux=ghi"},
		expected: Struct{
			"foo": Struct{
				"bar":  "abc",
				"baz":  "def",
				"quux": "ghi",
			},
		},
	},
	// Multiple keys with duplicate error
	{
		name:          "multikey-duplicate",
		kvSeparator:   "=",
		pathSeparator: ".",
		orig:          Struct{},
		inputs:        []string{"foo=bar", "foo=quux"},
		expected: Struct{
			"foo": "bar",
		},
		expectedResults: []error{
			nil,
			ErrAlreadyExists,
		},
	},
	// Multiple keys with type overwrite error
	{
		name:          "multikey-type-mismatch",
		kvSeparator:   "=",
		pathSeparator: ".",
		orig:          Struct{},
		inputs:        []string{"foo.bar=baz", "foo=quux"},
		expected: Struct{
			"foo": Struct{
				"bar": "baz",
			},
		},
		expectedResults: []error{
			nil,
			ErrAlreadyExists,
		},
	},
	// Slice keys are autovivified and merged correctly
	{
		name:          "slice-values-only",
		kvSeparator:   "=",
		pathSeparator: ".",
		orig:          Struct{},
		inputs:        []string{"foo.bar[]=abc", "foo.bar[]=def", "foo.bar[]=ghi"},
		expected: Struct{
			"foo": Struct{
				"bar": []interface{}{"abc", "def", "ghi"},
			},
		},
	},
	// Slice keys on multiple levels
	{
		name:          "slice-values-multi",
		kvSeparator:   "=",
		pathSeparator: ".",
		orig:          Struct{},
		inputs:        []string{"foo.bar[]=abc", "foo.bar[]=def", "quux[]=ghi", "quux[]=jkl", "hello=world"},
		expected: Struct{
			"foo": Struct{
				"bar": []interface{}{"abc", "def"},
			},
			"quux":  []interface{}{"ghi", "jkl"},
			"hello": "world",
		},
	},
	{
		name:          "non-empty-original",
		kvSeparator:   "=",
		pathSeparator: ".",
		orig: Struct{
			"foo": "bar",
		},
		inputs: []string{"baz=quux"},
		expected: Struct{
			"baz": "quux",
			"foo": "bar",
		},
	},
	{
		name:          "error-slice",
		kvSeparator:   "=",
		pathSeparator: ".",
		orig:          Struct{},
		inputs:        []string{"foo[].bar=abc", "foo[].baz=def"},
		expected: Struct{
			"foo": []interface{}{
				Struct{},
				Struct{},
			},
		},
		expectedResults: []error{
			ErrUnsupported,
			ErrUnsupported,
		},
	},
}

func TestProcessor(t *testing.T) {
	for _, test := range processorTests {
		test := test // range capture
		t.Run(test.name, func(t *testing.T) {
			results := []error{}
			errorCount := 0
			p := &Processor{
				Input:             test.orig,
				KeyPathSeparator:  test.pathSeparator,
				KeyValueSeparator: test.kvSeparator,
			}
			for _, input := range test.inputs {
				err := p.Process(input)
				results = append(results, err)
				if err != nil {
					errorCount++
				}
			}
			if test.expectedResults == nil {
				if errorCount > 0 {
					resultStrings := []string{}
					for _, res := range results {
						if res == nil {
							resultStrings = append(resultStrings, "")
						} else {
							resultStrings = append(resultStrings, res.Error())
						}
					}
					t.Fatalf("expected no errors in process, but got %d errors: %#v", errorCount, resultStrings)
				}
			} else {
				if len(results) != len(test.expectedResults) {
					t.Fatalf("error count mismatch, expected %d but got %d\n\texpected: %+v\n\tbut got : %+v",
						len(test.expectedResults), len(results), test.expectedResults, results)
				}
				for i, res := range test.expectedResults {
					if (res == nil && results[i] != nil) ||
						(res != nil && results[i] == nil) ||
						(res != nil && results[i] != nil && !errors.Is(results[i], res)) {
						t.Fatalf("error mismatch on result #%d,\n\texpected: %v,\n\tbut got : %v", i, res, results[i])
					}
				}
			}
			if out := p.Output(); !reflect.DeepEqual(out, test.expected) {
				t.Fatalf("expected to build structure %#v, but got %#v", test.expected, out)
			}
		})
	}
}
