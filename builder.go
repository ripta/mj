package main

import "encoding/json"
import "fmt"
import "strings"
import "github.com/pkg/errors"

// Uselessly generic interfaces deserve uselessly generic names, right?
// I hate myself (a little) for this.
type Struct map[string]interface{}

var (
	ErrAlreadyExists  = errors.New("already exists")
	ErrUnknownHandler = errors.New("unknown handler")
)

func (s Struct) Bytes() []byte {
	if bytes, err := json.Marshal(s); err == nil {
		return bytes
	}
	return []byte("{ }")
}

func (s Struct) Set(keyPath []string, value interface{}) error {
	if len(keyPath) == 0 {
		return errors.New("key path must not be empty")
	}
	keyName := keyPath[len(keyPath)-1]
	keyDir := keyPath[:len(keyPath)-1]

	var data interface{} = s
	for keyIdx, key := range keyDir {
		if len(key) == 0 {
			return fmt.Errorf("key sub-path #%d in %s must not be empty", keyIdx, keyPath)
		}

		usedPath := strings.Join(keyPath[0:keyIdx], ".")
		newPath := strings.Join(keyPath[0:keyIdx+1], ".")

		dataNext, err := s.setOn(data, key, Struct{})
		if err != nil && err != ErrAlreadyExists {
			return errors.Wrapf(err, "in key path %q while processing key path %q", usedPath, newPath)
		}
		data = dataNext
	}

	_, err := s.setOn(data, keyName, value)
	return errors.Wrapf(err, "while processing key path %v", keyPath)
}

func (s Struct) String() string {
	return string(s.Bytes())
}

func (s Struct) setOn(data interface{}, key string, value interface{}) (interface{}, error) {
	if nest, ok := data.(Struct); ok {
		if strings.HasSuffix(key, "[]") {
			key = strings.TrimSuffix(key, "[]")
			if nest[key] == nil {
				nest[key] = make([]interface{}, 0)
			}
			if slice, ok2 := nest[key].([]interface{}); ok2 {
				nest[key] = append(slice, value)
				return nest[key], nil
			}
		}
		if nest[key] != nil {
			// fmt.Errorf("already assigned a %T value and cannot be re-typed", nest[key])
			return nest[key], ErrAlreadyExists
		}
		nest[key] = value
		return nest[key], nil
	}
	// fmt.Errorf("no handler for type %T", data)
	return data, ErrUnknownHandler
}
