package main

import (
	"fmt"
	"io/ioutil"
	"strings"
)

type Processor struct {
	Input             Struct
	KeyValueSeparator string
	KeyPathSeparator  string
	ReadFilePrefix    string
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
	value := segments[1]

	if p.ReadFilePrefix != "" && strings.HasPrefix(value, p.ReadFilePrefix) {
		fn := strings.TrimPrefix(value, p.ReadFilePrefix)
		v, err := ioutil.ReadFile(fn)
		if err != nil {
			return err
		}
		value = string(v)
	}

	return p.Input.Set(keyPath, value)
}
