package main

import (
	"fmt"
	"strings"
)

type Processor struct {
	input Struct

	kvSeparator string
}

func (p *Processor) Output() Struct {
	return p.input
}

func (p *Processor) Process(arg string) error {
	segments := strings.SplitN(arg, p.kvSeparator, 2)
	if len(segments) != 2 {
		return fmt.Errorf("missing separator (%q) in %q", p.kvSeparator, arg)
	}

	return p.input.Set(segments[0], segments[1])
}
