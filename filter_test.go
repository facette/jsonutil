package jsonutil

import (
	"reflect"
	"testing"
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

	result := FilterSlice(foo, []string{})
	if !reflect.DeepEqual(result, expected) {
		t.Logf("\nExpected %#v\nbut got  %#v", expected, result)
		t.Fail()
	}
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

	result := FilterStruct(foo, []string{})
	if !reflect.DeepEqual(result, expected) {
		t.Logf("\nExpected %#v\nbut got  %#v", expected, result)
		t.Fail()
	}

	delete(expected, "omit_notempty")

	result = FilterStruct(foo, []string{"string", "bar"})
	if !reflect.DeepEqual(result, expected) {
		t.Logf("\nExpected %#v\nbut got  %#v", expected, result)
		t.Fail()
	}

	// Test 2nd level
	expected = map[string]interface{}{
		"string": foo.String,
		"bar": map[string]interface{}{
			"int": foo.Bar.Int,
		},
	}

	result = FilterStruct(foo, []string{"string", "bar.int"})
	if !reflect.DeepEqual(result, expected) {
		t.Logf("\nExpected %#v\nbut got  %#v", expected, result)
		t.Fail()
	}

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

	result = FilterStruct(foo, []string{"string", "bar.baz.bool"})
	if !reflect.DeepEqual(result, expected) {
		t.Logf("\nExpected %#v\nbut got  %#v", expected, result)
		t.Fail()
	}
}
