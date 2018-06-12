package jsonutil

import (
	"os"
	"reflect"
	"testing"
)

func TestMarshalFile(t *testing.T) {
	var (
		testFile = "./test.json"
		foo      = Foo{
			String: "test",
			Bar: Bar{
				Int: 42,
				Baz: []Baz{
					Baz{Bool: true, Float64: 42.42},
				},
			},
		}
		foo2 Foo
	)

	if err := MarshalFile(testFile, foo); err != nil {
		t.Logf("error: %s", err)
		t.Fail()
	}

	if err := UnmarshalFile(testFile, &foo2); err != nil {
		t.Logf("error: %s", err)
		t.Fail()
	}

	if !reflect.DeepEqual(foo, foo2) {
		t.Logf("\nExpected %#v\nbut got  %#v", foo, foo2)
		t.Fail()
	}

	os.Remove(testFile)
}
