package main

import "encoding/json"
import "errors"
import "fmt"
import "strings"

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

	var data interface{} = s
	for keyIdx, key := range keyPath {
		if len(key) == 0 {
			return fmt.Errorf("key sub-path #%d in %s must not be empty", keyIdx, keyPath)
		}

		usedPath := strings.Join(keyPath[0:keyIdx], ".")
		newPath := strings.Join(keyPath[0:keyIdx+1], ".")

		// TODO(rpasay): a better, safer alternative?
		if nest, ok := data.(Struct); ok {
			if keyIdx == len(keyPath)-1 {
				// Check to ensure nest[key] wasn't already initialized
				if nest[key] == nil {
					nest[key] = value
				} else {
					return fmt.Errorf("key path %q was already assigned a %T value, and cannot be overwritten", newPath, nest[key])
				}
			} else if nest[key] == nil {
				nest[key] = Struct{}
			}

			data = nest[key]
		} else {
			return fmt.Errorf("key path %q was already assigned a %T value, and cannot be reassigned in %q", usedPath, data, newPath)
		}
	}

	return nil
}

func (s Struct) String() string {
	return string(s.Bytes())
}
