package jsonutil

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMarshalFile(t *testing.T) {
	var (
		expected = Foo{
			String: "test",
			Bar: Bar{
				Int: 42,
				Baz: []Baz{
					Baz{Bool: true, Float64: 42.42},
				},
			},
		}
		actual Foo
	)

	tmpFile, err := ioutil.TempFile("", "jsonutil-marshal_")
	assert.Nil(t, err)
	defer os.Remove(tmpFile.Name())

	err = MarshalFile(tmpFile.Name(), expected)
	assert.Nil(t, err)

	err = UnmarshalFile(tmpFile.Name(), &actual)
	assert.Nil(t, err)

	assert.Equal(t, expected, actual)
}
