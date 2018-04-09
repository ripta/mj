package main

import "encoding/json"
import "fmt"
import "strings"
import "github.com/pkg/errors"

// Uselessly generic interfaces deserve uselessly generic names, right?
// I hate myself (a little) for this.
type Struct map[string]interface{}

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

		// TODO(rpasay): a better, safer alternative?
		if nest, ok := data.(Struct); ok {
			if nest[key] == nil {
				nest[key] = Struct{}
			}

			data = nest[key]
		} else {
			return fmt.Errorf("key path %q was already assigned a %T value, and cannot be reassigned in %q", usedPath, data, newPath)
		}
	}

	err := s.setOn(data, keyName, value)
	return errors.Wrapf(err, "while processing key path %v", keyPath)
}

func (s Struct) String() string {
	return string(s.Bytes())
}

func (s Struct) setOn(data interface{}, key string, value interface{}) error {
	if nest, ok := data.(Struct); ok {
		if nest[key] == nil {
			nest[key] = value
			return nil
		}
		return fmt.Errorf("already assigned a %T value and cannot be re-typed", nest[key])
	}
	return fmt.Errorf("no handler for type %T", data)
}
