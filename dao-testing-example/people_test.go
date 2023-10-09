package dte

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_People(t *testing.T) {
	db, dbErr := NewDB("file::memory:")
	require.Nil(t, dbErr)

	t.Cleanup(func() {
		db.Close()
	})

	t.Run("can insert, retrieve, and update people", func(t *testing.T) {
		expected := Person{
			Id:    1,
			Name:  "John Doe",
			Email: "john.doe@example.com",
			Photo: "https://example.com/john_doe.jpg",
		}

		err := db.AddPerson(&expected)
		require.Nil(t, err)

		found, err := db.GetPerson(1)
		assert.Nil(t, err)
		assert.Equal(t, &expected, found)

		expected.Photo = "https://example.com/john_doe_2.jpg"
		err = db.AddPerson(&expected)
		require.Nil(t, err)
		found, err = db.GetPerson(1)
		assert.Nil(t, err)
		assert.Equal(t, &expected, found)
	})
}
