package mj

import (
	"github.com/rogpeppe/go-internal/testscript"
	"testing"
)

func TestMain(m *testing.M) {
	testscript.RunMain(m, map[string]func() int{
		"mj": Run,
	})
}

func TestScripts(t *testing.T) {
	testscript.Run(t, testscript.Params{
		Dir: "testdata/scripts",
	})
}
