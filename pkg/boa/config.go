// Copyright 2020 oncilla
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package boa

import (
	"fmt"
	"net"
	"reflect"
	"strings"
	"time"

	"github.com/mitchellh/mapstructure"
	"github.com/spf13/pflag"
)

// ConfigRegistry is an abstraction of viper.Viper.
type ConfigRegistry interface {
	BindEnv(input ...string) error
	SetDefault(key string, value interface{})
}

// SetDefaults sets the default values based on the values contained in the
// provided config struct.
func SetDefaults(r ConfigRegistry, config interface{}) error {
	m := map[string]interface{}{}
	if err := mapstructure.Decode(config, &m); err != nil {
		return err
	}
	setDefaults(r, nil, m)
	return nil
}

func setDefaults(r ConfigRegistry, p path, config map[string]interface{}) {
	for key, value := range config {
		keyPath := p.Extend(key)
		if m, ok := value.(map[string]interface{}); ok {
			setDefaults(r, keyPath, m)
			continue
		}
		r.SetDefault(keyPath.String(), value)
	}
}

// BindEnv binds the environemt variables based on the config struct.
func BindEnv(r ConfigRegistry, config interface{}) error {
	m := map[string]interface{}{}
	if err := mapstructure.Decode(config, &m); err != nil {
		return err
	}
	return bindEnv(r, nil, m)
}

func bindEnv(r ConfigRegistry, p path, config map[string]interface{}) error {
	for key, value := range config {
		keyPath := p.Extend(key)
		if m, ok := value.(map[string]interface{}); ok {
			if err := bindEnv(r, keyPath, m); err != nil {
				return err
			}
			continue
		}
		if err := r.BindEnv(keyPath.String()); err != nil {
			return err
		}
	}
	return nil
}

// AddFlags adds flags to the provided flag set based on the config struct.
// Default are set according to the values present in the config struct.
func AddFlags(r *pflag.FlagSet, config interface{}) error {
	m := map[string]interface{}{}
	if err := mapstructure.Decode(config, &m); err != nil {
		return err
	}
	return addFlags(r, nil, m)
}

// nolint: gocyclo
func addFlags(r *pflag.FlagSet, p path, config map[string]interface{}) error {
	for key, value := range config {
		keyPath := p.Extend(key)
		if m, ok := value.(map[string]interface{}); ok {
			if err := addFlags(r, keyPath, m); err != nil {
				return err
			}
			continue
		}

		if v, ok := value.(pflag.Value); ok {
			r.Var(v, keyPath.String(), "")
			continue
		}

		t, v := reflect.TypeOf(value), reflect.ValueOf(value)
		if t.Kind() == reflect.Ptr || t.Kind() == reflect.Interface {
			t, v = t.Elem(), reflect.Indirect(v)
			if v.Kind() == reflect.Invalid {
				v = reflect.Zero(t)
			}
		}
		switch {
		case t == reflect.TypeOf(time.Duration(1)):
			r.Duration(keyPath.String(), time.Duration(v.Int()), "")
		case t == reflect.TypeOf(net.IP{}):
			r.IP(keyPath.String(), v.Interface().(net.IP), "")
		case t.Kind() == reflect.Bool:
			r.Bool(keyPath.String(), v.Bool(), "")
		case t.Kind() == reflect.Float32:
			r.Float32(keyPath.String(), float32(v.Float()), "")
		case t.Kind() == reflect.Float64:
			r.Float64(keyPath.String(), v.Float(), "")
		case t.Kind() == reflect.Int8:
			r.Int8(keyPath.String(), int8(v.Int()), "")
		case t.Kind() == reflect.Int16:
			r.Int16(keyPath.String(), int16(v.Int()), "")
		case t.Kind() == reflect.Int32:
			r.Int32(keyPath.String(), int32(v.Int()), "")
		case t.Kind() == reflect.Int64:
			r.Int64(keyPath.String(), v.Int(), "")
		case t.Kind() == reflect.Int:
			r.Int(keyPath.String(), int(v.Int()), "")
		case t.Kind() == reflect.Uint8:
			r.Uint8(keyPath.String(), uint8(v.Uint()), "")
		case t.Kind() == reflect.Uint16:
			r.Uint16(keyPath.String(), uint16(v.Uint()), "")
		case t.Kind() == reflect.Uint32:
			r.Uint32(keyPath.String(), uint32(v.Uint()), "")
		case t.Kind() == reflect.Uint64:
			r.Uint64(keyPath.String(), v.Uint(), "")
		case t.Kind() == reflect.Uint:
			r.Uint(keyPath.String(), uint(v.Uint()), "")
		case t.Kind() == reflect.String:
			r.String(keyPath.String(), v.String(), "")
		case t.Kind() == reflect.Slice:
			switch t.Elem().Kind() {
			case reflect.Bool:
				r.BoolSlice(keyPath.String(), v.Interface().([]bool), "")
			case reflect.Int32:
				r.Int32Slice(keyPath.String(), v.Interface().([]int32), "")
			case reflect.Int64:
				r.Int64Slice(keyPath.String(), v.Interface().([]int64), "")
			case reflect.Int:
				r.IntSlice(keyPath.String(), v.Interface().([]int), "")
			case reflect.Uint:
				r.UintSlice(keyPath.String(), v.Interface().([]uint), "")
			case reflect.String:
				r.StringSlice(keyPath.String(), v.Interface().([]string), "")
			}
		default:
			return fmt.Errorf("unsupported value: %s (%T)", keyPath, value)
		}
	}
	return nil
}

type path []string

func (p path) Extend(key string) path {
	return append([]string(p), key)
}

func (p path) String() string {
	return strings.Join([]string(p), ".")
}
