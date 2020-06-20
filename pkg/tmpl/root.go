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

package tmpl

// Root is the template for creating a new root command file.
const Root = `{{ if .Copyright.Author }}// Copyright {{.Copyright.Year}} {{ .Copyright.Author }}{{ end }}
{{ if .License.Commented }}{{ .License.Commented }}{{ end }}



package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// CommandPather returns the path to a command.
type CommandPather interface {
	CommandPath() string
}

func main() {
	cmd := &cobra.Command{
		Use:           "{{.Name}}",
		Short:         "{{.Name}} does amazing work!",
		SilenceErrors: true,
	}
	cmd.AddCommand(
		newCompletion(cmd),
		newVersion(cmd),
	)
	if err := cmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s", err)
		os.Exit(1)
	}
}
`
