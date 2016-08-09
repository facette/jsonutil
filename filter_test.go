package jsonutil

import (
	"reflect"
	"testing"
)

type Foo struct {
	String string `json:"string"`
	Bar    Bar    `json:"bar"`
}

type Bar struct {
	Int int   `json:"int"`
	Baz []Baz `json:"baz"`
}

type Baz struct {
	Bool    bool    `json:"bool"`
	Float64 float64 `json:"float64"`
}

func TestFilterStruct(t *testing.T) {
	foo := Foo{
		String: "abc",
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
		"string": foo.String,
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
