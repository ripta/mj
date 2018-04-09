package main

import (
	"fmt"
	"strings"
)

type Processor struct {
	Input             Struct
	KeyValueSeparator string
	KeyPathSeparator  string
}

func (p *Processor) Output() Struct {
	return p.Input
}

func (p *Processor) Process(arg string) error {
	segments := strings.SplitN(arg, p.KeyValueSeparator, 2)
	if len(segments) != 2 {
		return fmt.Errorf("missing separator (%q) in %q", p.KeyValueSeparator, arg)
	}

	keyPath := strings.Split(segments[0], p.KeyPathSeparator)
	return p.Input.Set(keyPath, segments[1])
}
