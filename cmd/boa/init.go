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

package main

import (
	"fmt"
	"os"
	"time"

	"github.com/oncilla/boa/pkg/gen"
	"github.com/spf13/cobra"
)

func newInit(pather CommandPather) *cobra.Command {
	var flags struct {
		author  string
		license string
		path    string
	}

	var cmd = &cobra.Command{
		Use:     "init <project name>",
		Aliases: []string{"initialize", "create"},
		Short:   "Initialize a venom-free cobra application",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.SilenceUsage = true
			path := flags.path
			if path == "" {
				var err error
				if path, err = os.Getwd(); err != nil {
					return err
				}
			}
			license, err := gen.FindLicense(flags.license)
			if err != nil {
				return err
			}
			p := gen.Project{
				Name: args[0],
				Copyright: gen.Copyright{
					Year:   time.Now().Year(),
					Author: flags.author,
				},
				License: license,
			}
			if err := p.Create(path); err != nil {
				return err
			}
			fmt.Println("Created project at", path)
			return nil
		},
	}
	cmd.Flags().StringVarP(&flags.author, "author", "a", "YOUR_NAME", "author name for copyright attribution")
	cmd.Flags().StringVarP(&flags.license, "license", "l", "apache", "name of license for the project")
	cmd.Flags().StringVarP(&flags.path, "path", "p", "", "path to main package")
	return cmd
}
