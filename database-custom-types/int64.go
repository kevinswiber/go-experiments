package dbtypes

import (
	"database/sql/driver"
	"reflect"
)

// The purpose of this set is to determine if "inherited" custom types
// need their own Valuer and Scanner implementations, or if the root one
// will suffice.

// RootLong is the base type that other "long" types will extend.
type RootLong int64

// ChildLong extends [RootLong].
type ChildLong RootLong

func (rl RootLong) Value() (driver.Value, error) {
	return int64(rl), nil
}

func (rl *RootLong) Scanner(value any) error {
	if value == nil {
		return nil
	}

	rv := reflect.TypeOf(value)
	println(rv.Name())

	return nil
}
