package jsonutil

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type Foo struct {
	String       string `json:"string"`
	Skip         string `json:"-"`
	OmitEmpty    string `json:"omit_empty,omitempty"`
	OmitNotEmpty string `json:"omit_notempty,omitempty"`
	Bar          Bar    `json:"bar"`
}

type Bar struct {
	Int int   `json:"int"`
	Baz []Baz `json:"baz"`
}

type Baz struct {
	Bool    bool    `json:"bool"`
	Float64 float64 `json:"float64"`
}

func TestFilterMap(t *testing.T) {
	source := map[string]interface{}{
		"string": "abc",
		"int":    123,
		"map": map[string]interface{}{
			"1": "a",
			"2": "b",
			"3": "c",
		},
	}

	actual := FilterMap(source, []string{})
	assert.Equal(t, source, actual)

	expected := map[string]interface{}{
		"string": "abc",
	}

	actual = FilterMap(source, []string{"string"})
	assert.Equal(t, expected, actual)

	expected = map[string]interface{}{
		"string": "abc",
		"map": map[string]interface{}{
			"3": "c",
		},
	}

	actual = FilterMap(source, []string{"string", "map.3"})
	assert.Equal(t, expected, actual)
}

func TestFilterSlice(t *testing.T) {
	foo := []Foo{{
		String:       "abc",
		OmitNotEmpty: "def",
		Bar: Bar{
			Int: 123,
			Baz: []Baz{
				{Bool: true, Float64: 45.6},
				{Bool: false, Float64: 78.9},
			},
		},
	}}

	expected := []map[string]interface{}{{
		"string":        foo[0].String,
		"omit_notempty": foo[0].OmitNotEmpty,
		"bar": map[string]interface{}{
			"int": foo[0].Bar.Int,
			"baz": []map[string]interface{}{
				{"bool": foo[0].Bar.Baz[0].Bool, "float64": foo[0].Bar.Baz[0].Float64},
				{"bool": foo[0].Bar.Baz[1].Bool, "float64": foo[0].Bar.Baz[1].Float64},
			},
		},
	}}

	actual := FilterSlice(foo, []string{})
	assert.Equal(t, expected, actual)
}

func TestFilterStruct(t *testing.T) {
	foo := Foo{
		String:       "abc",
		OmitNotEmpty: "def",
		Bar: Bar{
			Int: 123,
			Baz: []Baz{
				{Bool: true, Float64: 45.6},
				{Bool: false, Float64: 78.9},
			},
		},
	}

	// Test 1st level
	expected := map[string]interface{}{
		"string":        foo.String,
		"omit_notempty": foo.OmitNotEmpty,
		"bar": map[string]interface{}{
			"int": foo.Bar.Int,
			"baz": []map[string]interface{}{
				{"bool": foo.Bar.Baz[0].Bool, "float64": foo.Bar.Baz[0].Float64},
				{"bool": foo.Bar.Baz[1].Bool, "float64": foo.Bar.Baz[1].Float64},
			},
		},
	}

	actual := FilterStruct(foo, []string{})
	assert.Equal(t, expected, actual)

	delete(expected, "omit_notempty")

	actual = FilterStruct(foo, []string{"string", "bar"})
	assert.Equal(t, expected, actual)

	// Test 2nd level
	expected = map[string]interface{}{
		"string": foo.String,
		"bar": map[string]interface{}{
			"int": foo.Bar.Int,
		},
	}

	actual = FilterStruct(foo, []string{"string", "bar.int"})
	assert.Equal(t, expected, actual)

	// Test 3rd level
	expected = map[string]interface{}{
		"string": foo.String,
		"bar": map[string]interface{}{
			"baz": []map[string]interface{}{
				{"bool": foo.Bar.Baz[0].Bool},
				{"bool": foo.Bar.Baz[1].Bool},
			},
		},
	}

	actual = FilterStruct(foo, []string{"string", "bar.baz.bool"})
	assert.Equal(t, expected, actual)
}
