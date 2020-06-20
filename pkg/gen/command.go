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
	"os"
	"os/exec"
	"strings"
	"text/template"

	"github.com/oncilla/boa/pkg/tmpl"
)

// Copyright holds the values for the copyright attribution.
type Copyright struct {
	Year   int
	Author string
}

// Command generates new cobra commands.
type Command struct {
	Name      string
	Copyright Copyright
	License   License
	Flags     []Flag
	Imports   []string
}

// UpperName returns the command name starting with a capital letter.
func (c Command) UpperName() string {
	return strings.ToUpper(c.Name[:1]) + c.Name[1:]
}

// Create writes the templated command and formats it using 'gofmt'.
func (c Command) Create(name string) error {
	if err := notExists(name); err != nil {
		return err
	}
	file, err := os.Create(name)
	if err != nil {
		return err
	}
	t := template.Must(template.New("").Parse(tmpl.Command))
	if err := t.Execute(file, c); err != nil {
		return err
	}
	if err := file.Close(); err != nil {
		return err
	}
	cmd := exec.Command("gofmt", "-w", file.Name())
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}
