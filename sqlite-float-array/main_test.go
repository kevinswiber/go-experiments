package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_String(t *testing.T) {
	ll := LatLng{42.12345, -42.54321}
	str := fmt.Sprintf("%v", ll)
	assert.Equal(t, "(42.12345,-42.54321)", str)
}

func Test_Value(t *testing.T) {
	ll := LatLng{42.12345, -42.54321}
	str, err := ll.Value()
	assert.Nil(t, err)
	assert.Equal(t, "(42.12345,-42.54321)", str)
}

func Test_Scan(t *testing.T) {
	t.Run("only scans strings", func(t *testing.T) {
		ll := LatLng{}
		err := ll.Scan(42)
		assert.ErrorContains(t, err, "value must be a string, got: int")
	})

	t.Run("empty instance is empty", func(t *testing.T) {
		expected := LatLng{}
		source := LatLng{}

		err := source.Scan("")
		assert.Nil(t, err)
		assert.Equal(t, expected, source)
	})

	t.Run("parts must match expectation", func(t *testing.T) {
		ll := LatLng{}
		err := ll.Scan("(42.12345)")
		assert.ErrorContains(t, err, "expected (float64,float64) but received: (42.12345)")
	})

	t.Run("parts must be floats", func(t *testing.T) {
		ll := LatLng{}

		err := ll.Scan("(foo,bar)")
		assert.ErrorContains(t, err, "parsing \"foo\"")

		err = ll.Scan("(42,bar)")
		assert.ErrorContains(t, err, "parsing \"bar\"")
	})

	t.Run("scans correctly", func(t *testing.T) {
		ll := LatLng{}
		err := ll.Scan("(42.12345,-42.54321)")
		assert.Nil(t, err)
		assert.Equal(t, LatLng{42.12345, -42.54321}, ll)
	})
}
