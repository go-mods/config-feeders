// Package default_feeder provides a golobby/config feeder that fills a structure
// with the values declared in its `default` tags.
package default_feeder

import (
	"errors"
	"fmt"
	"reflect"
	"time"
	"unsafe"

	"github.com/golobby/cast"
)

// Default is a feeder to be used with golobby/config.
//
// It feeds using a default value declared as tags in the structure.
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
	return f.fillStruct(v.Elem())
}

// fillStruct sets the tagged fields of s to their default value, walking through
// nested structs and non-nil pointers to structs.
func (f Default) fillStruct(s reflect.Value) error {
	for i := 0; i < s.NumField(); i++ {
		if t, exist := s.Type().Field(i).Tag.Lookup("default"); exist {
			// handle time.Duration
			if s.Type().Field(i).Type.String() == "time.Duration" {
				d, err := time.ParseDuration(t)
				if err != nil {
					return fmt.Errorf("default: cannot set `%v` field; err: %w", s.Type().Field(i).Name, err)
				}
				// #nosec G103
				ptr := reflect.NewAt(s.Field(i).Type(), unsafe.Pointer(s.Field(i).UnsafeAddr())).Elem()
				ptr.Set(reflect.ValueOf(d))
			} else {
				v, err := cast.FromType(t, s.Type().Field(i).Type)
				if err != nil {
					return fmt.Errorf("default: cannot set `%v` field; err: %w", s.Type().Field(i).Name, err)
				}
				// #nosec G103
				ptr := reflect.NewAt(s.Field(i).Type(), unsafe.Pointer(s.Field(i).UnsafeAddr())).Elem()
				ptr.Set(reflect.ValueOf(v))
			}
		} else if s.Type().Field(i).Type.Kind() == reflect.Struct {
			if err := f.fillStruct(s.Field(i)); err != nil {
				return err
			}
		} else if s.Type().Field(i).Type.Kind() == reflect.Pointer {
			if !s.Field(i).IsZero() && s.Field(i).Elem().Type().Kind() == reflect.Struct {
				if err := f.fillStruct(s.Field(i).Elem()); err != nil {
					return err
				}
			}
		}
	}
	return nil
}
