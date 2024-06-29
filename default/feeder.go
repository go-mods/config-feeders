package default_feeder

import (
	"errors"
	"fmt"
	"github.com/golobby/cast"
	"reflect"
	"time"
	"unsafe"
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
func (f Default) Feed(structure interface{}) error {
	inputType := reflect.TypeOf(structure)
	if inputType != nil {
		if inputType.Kind() == reflect.Ptr {
			if inputType.Elem().Kind() == reflect.Struct {
				return f.fillStruct(reflect.ValueOf(structure).Elem())
			}
		}
	}
	return errors.New("env: invalid structure")
}

// fillStruct permet de définir les champs de la structure avec la valeur par défaut
func (f Default) fillStruct(s reflect.Value) error {
	for i := 0; i < s.NumField(); i++ {
		if t, exist := s.Type().Field(i).Tag.Lookup("default"); exist {
			// handle time.Duration
			if s.Type().Field(i).Type.String() == "time.Duration" {
				d, err := time.ParseDuration(t)
				if err != nil {
					return fmt.Errorf("default: cannot set `%v` field; err: %v", s.Type().Field(i).Name, err)
				}
				// #nosec G103
				ptr := reflect.NewAt(s.Field(i).Type(), unsafe.Pointer(s.Field(i).UnsafeAddr())).Elem()
				ptr.Set(reflect.ValueOf(d))
			} else {
				v, err := cast.FromType(t, s.Type().Field(i).Type)
				if err != nil {
					return fmt.Errorf("default: cannot set `%v` field; err: %v", s.Type().Field(i).Name, err)
				}
				// #nosec G103
				ptr := reflect.NewAt(s.Field(i).Type(), unsafe.Pointer(s.Field(i).UnsafeAddr())).Elem()
				ptr.Set(reflect.ValueOf(v))
			}
		} else if s.Type().Field(i).Type.Kind() == reflect.Struct {
			if err := f.fillStruct(s.Field(i)); err != nil {
				return err
			}
		} else if s.Type().Field(i).Type.Kind() == reflect.Ptr {
			if !s.Field(i).IsZero() && s.Field(i).Elem().Type().Kind() == reflect.Struct {
				if err := f.fillStruct(s.Field(i).Elem()); err != nil {
					return err
				}
			}
		}
	}
	return nil
}
