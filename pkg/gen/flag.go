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

package gen

import (
	"fmt"
	"sort"
	"strings"
)

// Flag holds the flag description.
type Flag struct {
	Name     string
	Type     string
	Register string
	Default  string
}

// ParseFlags parses a list of flags.
func ParseFlags(inputs []string) ([]Flag, []string, error) {
	flags := make([]Flag, 0, len(inputs))
	imports := map[string]struct{}{}
	for _, input := range inputs {
		f, imp, err := ParseFlag(input)
		if err != nil {
			return nil, nil, err
		}
		if imp != "" {
			imports[imp] = struct{}{}
		}
		flags = append(flags, f)
	}
	var unique []string
	for imp := range imports {
		unique = append(unique, imp)
	}
	sort.Strings(unique)
	return flags, unique, nil

}

// ParseFlag parses a single flag description.
func ParseFlag(input string) (Flag, string, error) {
	s := strings.Split(input, ":")
	if len(s) != 2 {
		return Flag{}, "", fmt.Errorf("malformed flag: %s", input)
	}
	for flagType, v := range supportedFlag {
		if strings.EqualFold(flagType, s[1]) {
			return Flag{
				Name:     s[0],
				Type:     v.Type,
				Default:  v.Default,
				Register: v.Register,
			}, v.Import, nil
		}
	}
	return Flag{}, "", fmt.Errorf("unsupported flag type: %s", s[1])
}

var supportedFlag = map[string]struct {
	Type     string
	Register string
	Default  string
	Import   string
}{
	"bool": {
		Type:     "bool",
		Default:  "false",
		Register: "BoolVar",
	},
	"[]bool": {
		Type:     "[]bool",
		Default:  "nil",
		Register: "BoolSliceVar",
	},
	"bytes": {
		Type:     "[]byte",
		Default:  "nil",
		Register: "BytesBase64Var",
	},
	"hexbytes": {
		Type:     "[]byte",
		Default:  "nil",
		Register: "BytesHexVar",
	},
	"[]duration": {
		Type:     "[]time.Duration",
		Default:  "nil",
		Register: "DurationSliceVar",
		Import:   "time",
	},
	"duration": {
		Type:     "time.Duration",
		Default:  "0",
		Register: "DurationVar",
		Import:   "time",
	},
	"float32": {
		Type:     "float32",
		Default:  "0",
		Register: "Float32Var",
	},
	"[]float32": {
		Type:     "[]float32",
		Default:  "nil",
		Register: "Float32SliceVar",
	},
	"float64": {
		Type:     "float64",
		Default:  "0",
		Register: "Float64Var",
	},
	"[]float64": {
		Type:     "[]float64",
		Default:  "nil",
		Register: "Float64SliceVar",
	},
	"ip": {
		Type:     "net.IP",
		Default:  "nil",
		Register: "IPVar",
		Import:   "net",
	},
	"int": {
		Type:     "int",
		Default:  "0",
		Register: "IntVar",
	},
	"[]int": {
		Type:     "[]int",
		Default:  "nil",
		Register: "IntSliceVar",
	},
	"int8": {
		Type:     "int8",
		Default:  "0",
		Register: "Int8Var",
	},
	"[]int8": {
		Type:     "[]int8",
		Default:  "nil",
		Register: "Int8SliceVar",
	},
	"int16": {
		Type:     "int16",
		Default:  "0",
		Register: "Int16Var",
	},
	"[]int16": {
		Type:     "[]int16",
		Default:  "nil",
		Register: "Int16SliceVar",
	},
	"int32": {
		Type:     "int32",
		Default:  "0",
		Register: "Int32Var",
	},
	"[]int32": {
		Type:     "[]int32",
		Default:  "nil",
		Register: "Int32SliceVar",
	},
	"int64": {
		Type:     "int64",
		Default:  "0",
		Register: "Int64Var",
	},
	"[]int64": {
		Type:     "[]int64",
		Default:  "nil",
		Register: "Int64SliceVar",
	},
	"uint": {
		Type:     "uint",
		Default:  "0",
		Register: "UintVar",
	},
	"[]uint": {
		Type:     "[]uint",
		Default:  "nil",
		Register: "UintSliceVar",
	},
	"uint8": {
		Type:     "uint8",
		Default:  "0",
		Register: "Uint8Var",
	},
	"[]uint8": {
		Type:     "[]uint8",
		Default:  "nil",
		Register: "Uint8SliceVar",
	},
	"uint16": {
		Type:     "uint16",
		Default:  "0",
		Register: "Uint16Var",
	},
	"[]uint16": {
		Type:     "[]uint16",
		Default:  "nil",
		Register: "Uint16SliceVar",
	},
	"uint32": {
		Type:     "uint32",
		Default:  "0",
		Register: "Uint32Var",
	},
	"[]uint32": {
		Type:     "[]uint32",
		Default:  "nil",
		Register: "Uint32SliceVar",
	},
	"uint64": {
		Type:     "uint64",
		Default:  "0",
		Register: "Uint64Var",
	},
	"[]uint64": {
		Type:     "[]uint64",
		Default:  "nil",
		Register: "Uint64SliceVar",
	},
	"string": {
		Type:     "string",
		Default:  `""`,
		Register: "StringVar",
	},
	"[]string": {
		Type:     "[]string",
		Default:  "nil",
		Register: "StringSliceVar",
	},
}
