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

// Command is the template for creating new command files.
const Command = `{{ if .Copyright.Author }}// Copyright {{.Copyright.Year}} {{ .Copyright.Author }}{{ end }}
{{ if .License.Commented }}{{ .License.Commented }}{{ end }}

package main

import (
	"fmt"{{ range $element := .Imports }}
	"{{$element}}"{{end}}

	"github.com/spf13/cobra"
)

func new{{ .UpperName }}(pather CommandPather) *cobra.Command {
	{{ if eq (len .Flags) 0 }}var flags struct{
		sample bool
	} {{ else }} var flags struct { {{ range .Flags }}
		{{.Name}} {{.Type}} {{ end }}
	} {{ end }}

	var cmd = &cobra.Command{
		Use:     "{{.Name}} <arg>",
		Short:   "{{.Name}} does amazing work!",
		Example: fmt.Sprintf("  %[1]s {{.Name}} --sample", pather.CommandPath()),
		RunE: func(cmd *cobra.Command, args []string) error {
			// Add basic sanity checks, where the usage help message should be
			// printed on error, before this line. After this line, the usage
			// message is no longer printed on error.
			cmd.SilenceUsage = true

			// TODO: Amazing work goes here!
			return nil
		},
	}
	{{ if eq (len .Flags) 0 }}cmd.Flags().BoolVarP(&flags.sample, "sample", "s", false, "sample flag"){{else}} {{ range .Flags }}
	cmd.Flags().{{.Register}}(&flags.{{.Name}}, "{{.Name}}", {{.Default}}, "{{.Name}} description") {{ end }} {{ end }}
	return cmd
}
`
