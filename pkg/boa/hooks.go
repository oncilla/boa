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
	"net"
	"reflect"

	"github.com/mitchellh/mapstructure"
	"github.com/oncilla/boa/pkg/boa/flag"
)

// DefaultDecodeHooks returns a list of useful decoding hooks.
func DefaultDecodeHooks() []mapstructure.DecodeHookFunc {
	return []mapstructure.DecodeHookFunc{
		mapstructure.StringToIPHookFunc(),
		mapstructure.StringToTimeDurationHookFunc(),
		mapstructure.StringToSliceHookFunc(","),
		StringToTCPAddrHookFunc(),
		StringToUDPAddrHookFunc(),
	}
}

// StringToTCPAddrHookFunc returns a DecodeHookFunc that converts
// strings to net.TCPAddr
func StringToTCPAddrHookFunc() mapstructure.DecodeHookFunc {
	return func(
		f reflect.Type,
		t reflect.Type,
		data interface{}) (interface{}, error) {
		if f.Kind() != reflect.String {
			return data, nil
		}
		if t != reflect.TypeOf(net.TCPAddr{}) && t != reflect.TypeOf(flag.TCPAddr{}) {
			return data, nil
		}
		addr, err := net.ResolveTCPAddr("tcp", data.(string))
		if err != nil {
			return nil, err
		}
		return addr, nil
	}
}

// StringToUDPAddrHookFunc returns a DecodeHookFunc that converts
// strings to net.UDPAddr
func StringToUDPAddrHookFunc() mapstructure.DecodeHookFunc {
	return func(
		f reflect.Type,
		t reflect.Type,
		data interface{}) (interface{}, error) {
		if f.Kind() != reflect.String {
			return data, nil
		}
		if t != reflect.TypeOf(net.UDPAddr{}) && t != reflect.TypeOf(flag.UDPAddr{}) {
			return data, nil
		}
		addr, err := net.ResolveUDPAddr("udp", data.(string))
		if err != nil {
			return nil, err
		}
		return addr, nil
	}
}
