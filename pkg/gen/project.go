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
	"os"
	"os/exec"
	"text/template"

	"github.com/oncilla/boa/pkg/tmpl"
	"golang.org/x/xerrors"
)

// Project generates new cobra projects.
type Project struct {
	Name      string
	Copyright Copyright
	License   License
}

// Create writes the templated project and formats it using 'gofmt'.
func (p Project) Create(path string) error {
	if err := os.MkdirAll(path, 0755); err != nil {
		return err
	}

	for name, tempStr := range map[string]string{
		fmt.Sprintf("%s/%s.go", path, p.Name): tmpl.Root,
		fmt.Sprintf("%s/completion.go", path): tmpl.Completion,
		fmt.Sprintf("%s/version.go", path):    tmpl.Version,
	} {
		if err := notExists(name); err != nil {
			return err
		}
		file, err := os.Create(name)
		if err != nil {
			return err
		}
		t := template.Must(template.New("").Parse(tempStr))
		if err := t.Execute(file, p); err != nil {
			return err
		}
		if err := file.Close(); err != nil {
			return err
		}
	}

	cmd := exec.Command("gofmt", "-w", path)
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}

func notExists(file string) error {
	_, err := os.Stat(file)
	if xerrors.Is(err, os.ErrNotExist) {
		return nil
	}
	if err != nil {
		return err
	}
	return fmt.Errorf("file already exists: %s", file)
}
