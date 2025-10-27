package mj

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type ValueType int

const (
	TypeString ValueType = iota
	TypeInt
	TypeFloat
	TypeBool
	TypeNull
)

type Processor struct {
	Input             Struct
	KeyValueSeparator string
	KeyPathSeparator  string
	ReadFilePrefix    string
	TypeSeparator     string
}

func (p *Processor) Output() Struct {
	return p.Input
}

// parseTypeSuffix extracts type from key like "age:int" -> ("age", TypeInt)
func (p *Processor) parseTypeSuffix(key string) (string, ValueType) {
	if p.TypeSeparator == "" {
		return key, TypeString
	}

	// Find the last occurrence of TypeSeparator to handle keys with multiple separators
	lastIdx := strings.LastIndex(key, p.TypeSeparator)
	if lastIdx == -1 {
		return key, TypeString
	}

	suffix := key[lastIdx+len(p.TypeSeparator):]
	cleanKey := key[:lastIdx]

	switch suffix {
	case "int":
		return cleanKey, TypeInt
	case "float":
		return cleanKey, TypeFloat
	case "bool":
		return cleanKey, TypeBool
	case "null":
		return cleanKey, TypeNull
	case "string":
		return cleanKey, TypeString
	default:
		// Not a recognized type suffix, treat the whole thing as the key
		return key, TypeString
	}
}

// parseTypedValue converts string value to appropriate type
func parseTypedValue(value string, vtype ValueType) (interface{}, error) {
	switch vtype {
	case TypeInt:
		if value == "" {
			return 0, nil
		}
		v, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("cannot parse %q as int: %w", value, err)
		}
		return v, nil
	case TypeFloat:
		if value == "" {
			return 0.0, nil
		}
		v, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return nil, fmt.Errorf("cannot parse %q as float: %w", value, err)
		}
		return v, nil
	case TypeBool:
		if value == "" {
			return false, nil
		}
		v, err := strconv.ParseBool(value)
		if err != nil {
			return nil, fmt.Errorf("cannot parse %q as bool: %w", value, err)
		}
		return v, nil
	case TypeNull:
		if value != "" {
			return nil, fmt.Errorf("cannot parse non-empty %q as null", value)
		}
		return nil, nil
	case TypeString:
		return value, nil
	default:
		return value, nil
	}
}

func (p *Processor) Process(arg string) error {
	segments := strings.SplitN(arg, p.KeyValueSeparator, 2)
	if len(segments) != 2 {
		return fmt.Errorf("missing separator (%q) in %q", p.KeyValueSeparator, arg)
	}

	keyPathStr := segments[0]
	value := segments[1]

	keyPathStr, vtype := p.parseTypeSuffix(keyPathStr)
	keyPath := strings.Split(keyPathStr, p.KeyPathSeparator)

	if p.ReadFilePrefix != "" && strings.HasPrefix(value, p.ReadFilePrefix) {
		fn := strings.TrimPrefix(value, p.ReadFilePrefix)
		v, err := os.ReadFile(fn)
		if err != nil {
			return err
		}
		value = string(v)
	}

	typedValue, err := parseTypedValue(value, vtype)
	if err != nil {
		return err
	}

	return p.Input.Set(keyPath, typedValue)
}
