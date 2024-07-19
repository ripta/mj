package mj

import "fmt"

type ConflictResolutionMethod string

const (
	Die   ConflictResolutionMethod = "die"
	First                          = "first"
	Last                           = "last"
)

func (c *ConflictResolutionMethod) Set(v string) error {
	switch v2 := ConflictResolutionMethod(v); v2 {
	case Die, First, Last:
		*c = v2
	default:
		return fmt.Errorf("must be one of: die, first, last")
	}
	return nil
}

func (c *ConflictResolutionMethod) String() string {
	return string(*c)
}
