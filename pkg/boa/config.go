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
	"strings"

	"github.com/mitchellh/mapstructure"
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

type path []string

func (p path) Extend(key string) path {
	return append([]string(p), key)
}

func (p path) String() string {
	return strings.Join([]string(p), ".")
}
