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
	"bufio"
	"fmt"
	"strings"
)

// License holds the software license agreement.
type License struct {
	Name   string
	Alias  []string
	Header string
	Text   string
}

// Commented returns the commented header.
func (l License) Commented() string {
	var builder strings.Builder
	scanner := bufio.NewScanner(strings.NewReader(l.Header))
	for scanner.Scan() {
		builder.WriteString("// ")
		builder.Write(scanner.Bytes())
		builder.WriteString("\n")
	}
	return builder.String()
}

// FindLicense finds the license by name from the supported set.
func FindLicense(name string) (License, error) {
	none := License{
		Name:   "none",
		Alias:  []string{"none"},
		Header: "",
		Text:   "",
	}
	for _, l := range []License{AGPL, Apache, FreeBSD, BSD, GPL2, GPL3, LGPL, MIT, none} {
		for _, alias := range l.Alias {
			if strings.EqualFold(alias, name) {
				return l, nil
			}
		}
	}
	return License{}, fmt.Errorf("license not found: %s", name)
}
