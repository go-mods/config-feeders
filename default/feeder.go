// Package default_feeder provides a golobby/config feeder that fills a structure
// with the values declared in its `default` tags.
package default_feeder

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"
	"unsafe"

	"github.com/golobby/cast"
)

var durationType = reflect.TypeOf(time.Duration(0))

// Default is a feeder to be used with golobby/config.
//
// It feeds using a default value declared as tags in the structure.
//
// The supported field types are strings, booleans, integers and floats, named types
// based on them, time.Duration, and slices of these (comma-separated values).
//
// Example:
//
//	type Config struct {
//		Host     string `default:"localhost"`
//		Port     int    `default:"8080"`
//	}
type Default struct{}

// Feed feeds the structure with the default values declared as tags in the structure.
//
// The structure must be a non-nil pointer to a struct. A tagged field is always
// overwritten, even if it already holds a value: add this feeder before the others.
func (f Default) Feed(structure interface{}) error {
	v := reflect.ValueOf(structure)
	if v.Kind() != reflect.Pointer || v.IsNil() || v.Elem().Kind() != reflect.Struct {
		return errors.New("default: invalid structure")
	}
	return f.fillStruct(v.Elem(), map[uintptr]bool{v.Pointer(): true})
}

// fillStruct sets the tagged fields of s to their default value, walking through
// nested structs and non-nil pointers to structs. The pointers already followed are
// kept in seen, so that a cyclic structure does not recurse forever.
func (f Default) fillStruct(s reflect.Value, seen map[uintptr]bool) error {
	for i := 0; i < s.NumField(); i++ {
		field, value := s.Type().Field(i), s.Field(i)

		if t, exist := field.Tag.Lookup("default"); exist {
			v, err := parse(t, field.Type)
			if err != nil {
				return fmt.Errorf("default: cannot set `%v` field; err: %w", field.Name, err)
			}
			// #nosec G103
			ptr := reflect.NewAt(value.Type(), unsafe.Pointer(value.UnsafeAddr())).Elem()
			ptr.Set(v)
		} else if field.Type.Kind() == reflect.Struct {
			if err := f.fillStruct(value, seen); err != nil {
				return err
			}
		} else if field.Type.Kind() == reflect.Pointer {
			if !value.IsNil() && value.Elem().Kind() == reflect.Struct && !seen[value.Pointer()] {
				seen[value.Pointer()] = true
				if err := f.fillStruct(value.Elem(), seen); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

// parse converts a default value to the type t: a slice takes comma-separated values.
func parse(value string, t reflect.Type) (reflect.Value, error) {
	if t.Kind() != reflect.Slice {
		return parseValue(value, t)
	}

	items := strings.Split(value, ",")
	slice := reflect.MakeSlice(t, 0, len(items))
	for _, item := range items {
		v, err := parseValue(strings.Trim(item, " \n\r"), t.Elem())
		if err != nil {
			return reflect.Value{}, err
		}
		slice = reflect.Append(slice, v)
	}
	return slice, nil
}

// parseValue converts a single default value to the type t, which may be a named type.
func parseValue(value string, t reflect.Type) (reflect.Value, error) {
	if t == durationType {
		d, err := time.ParseDuration(value)
		if err != nil {
			return reflect.Value{}, err
		}
		return reflect.ValueOf(d), nil
	}

	v, err := cast.FromString(value, t.Kind().String())
	if err != nil {
		return reflect.Value{}, err
	}
	return reflect.ValueOf(v).Convert(t), nil
}
