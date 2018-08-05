package jsonutil

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNullStringMarshalEmpty(t *testing.T) {
	actual, err := json.Marshal(NullString(""))
	assert.Nil(t, err)
	assert.Equal(t, []byte(`null`), actual)
}

func TestNullStringMarshalString(t *testing.T) {
	actual, err := json.Marshal(NullString("test"))
	assert.Nil(t, err)
	assert.Equal(t, []byte(`"test"`), actual)
}

func TestNullStringUnmarshalEmpty(t *testing.T) {
	var actual NullString

	err := json.Unmarshal([]byte("null"), &actual)
	assert.Nil(t, err)
	assert.Equal(t, NullString(""), actual)
}
