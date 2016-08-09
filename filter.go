package jsonutil

import (
	"encoding/json"
	"reflect"
	"strings"
)

// FilterStruct filters a struct given a list of JSON field paths.
func FilterStruct(v interface{}, fields []string) map[string]interface{} {
	var current reflect.Value

	result := make(map[string]interface{})

	stack := []reflect.Value{reflect.ValueOf(v)}

	for len(stack) > 0 {
		current, stack = stack[0], stack[1:]

		n := current.NumField()
		for i := 0; i < n; i++ {
			ft := current.Type().Field(i)
			f := current.Field(i)

			// Handle nested structures
			if ft.Anonymous {
				stack = append(stack, f)
				continue
			}

			// Get field name
			fname := filterBaseField(ft.Tag.Get("json"))

			if _, ok := f.Interface().(json.Marshaler); !ok && f.Kind() == reflect.Struct {
				// Handle sub struct filtering
				if !filterMatch(fname, fields) {
					continue
				} else if smap := FilterStruct(f.Interface(), filterFields(fname, fields)); len(smap) > 0 {
					result[fname] = smap
				}
			} else if f.Kind() == reflect.Slice {
				// Handle slice filtering
				slice := []map[string]interface{}{}

				n := f.Len()
				for i := 0; i < n; i++ {
					if !filterMatch(fname, fields) {
						continue
					} else if smap := FilterStruct(f.Index(i).Interface(), filterFields(fname, fields)); len(smap) > 0 {
						slice = append(slice, smap)
					}
				}

				if len(slice) > 0 {
					result[fname] = slice
				}
			} else if !filterMatch(fname, fields) {
				// Skip unwanted fields
				continue
			} else {
				// Set item value
				result[fname] = f.Interface()
			}
		}
	}

	return result
}

func filterMatch(name string, fields []string) bool {
	if len(fields) == 0 {
		return true
	}

	for _, s := range fields {
		if name == filterBaseField(s) {
			return true
		}
	}

	return false
}

func filterFields(prefix string, fields []string) []string {
	result := []string{}
	for _, s := range fields {
		if strings.HasPrefix(s, prefix+".") {
			result = append(result, strings.TrimPrefix(s, prefix+"."))
		}
	}

	return result
}

func filterBaseField(name string) string {
	return strings.SplitN(strings.SplitN(name, ",", 2)[0], ".", 2)[0]
}
