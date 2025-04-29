package mj

import (
	"github.com/rogpeppe/go-internal/testscript"
	"testing"
)

func TestMain(m *testing.M) {
	testscript.Main(m, map[string]func(){
		"mj": func() {
			_ = Run()
		},
	})
}

func TestScripts(t *testing.T) {
	testscript.Run(t, testscript.Params{
		Dir: "testdata/scripts",
	})
}
