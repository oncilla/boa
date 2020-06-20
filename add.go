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

func newAdd(pather CommandPather) *cobra.Command {
	var flags struct {
		author  string
		license string
		path    string
		flags   []string
	}

	help := `  %[1]s add ping --flags count:int,interval:duration
  %[1]s add pong --license apache`

	var cmd = &cobra.Command{
		Use:     "add <command name>",
		Aliases: []string{"cmd", "command"},
		Short:   "Add a command to a cobra application",
		Long: `Add a command to the cobra application.

The command is added to the main package of the cobra application. If main is
not located in the current working directory, supply the path with the appropriate
flag.

The license header is automatically added to the file. If you choose to not add
a license header, supply 'none' to the flag. Likewise, the copyright line can
be disabled by supplying an empty author name.

The generated command needs to be registered with the root command in the main
function.

This command supports adding flags to the generated command. To do so, specify
the desired flags as a comma separated list of 'name:type' pairs. In addition
to the basic go types, 'net.IP' and 'time.Duration' are supported with the type
identifiers 'ip' and 'duration'.

For example:

  names:[]string,addr:ip,interval:duration

Creates a command that supports the following flags:

  Flags:
        --names    strings    names description
        --addr     ip         addr description
        --interval duration   interval description
		`,
		Args:    cobra.ExactArgs(1),
		Example: fmt.Sprintf(help, pather.CommandPath()),
		RunE: func(cmd *cobra.Command, args []string) error {
			path := flags.path
			if path == "" {
				var err error
				if path, err = os.Getwd(); err != nil {
					return err
				}
			}
			cmd.SilenceUsage = true
			license, err := gen.FindLicense(flags.license)
			if err != nil {
				return err
			}
			cmdFlags, imports, err := gen.ParseFlags(flags.flags)
			if err != nil {
				return err
			}
			g := gen.Command{
				Copyright: gen.Copyright{
					Year:   time.Now().Year(),
					Author: flags.author,
				},
				License: license,
				Name:    args[0],
				Flags:   cmdFlags,
				Imports: imports,
			}
			name := fmt.Sprintf("%s/%s.go", path, args[0])
			if err := g.Create(name); err != nil {
				return err
			}
			fmt.Println("Created command at", name)
			fmt.Println("Make sure to register it with the root command")
			return nil
		},
	}
	cmd.Flags().StringVarP(&flags.author, "author", "a", "YOUR_NAME", "author name for copyright attribution")
	cmd.Flags().StringVarP(&flags.license, "license", "l", "apache", "name of license for the project")
	cmd.Flags().StringVarP(&flags.path, "path", "p", "", "path to main package")
	cmd.Flags().StringSliceVar(&flags.flags, "flags", nil, `flags to generate as comma separated list`)
	return cmd
}
