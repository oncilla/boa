// Copyright 2020 my-name
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"fmt"
	"net"

	"github.com/spf13/cobra"
)

func newServe(pather CommandPather) *cobra.Command {
	var flags struct {
		addr net.IP
		port uint16
	}

	var cmd = &cobra.Command{
		Use:     "serve <arg>",
		Short:   "serve does amazing work!",
		Example: fmt.Sprintf("  %[1]s serve --sample", pather.CommandPath()),
		RunE: func(cmd *cobra.Command, args []string) error {
			// Add basic sanity checks, where the usage help message should be
			// printed on error, before this line. After this line, the usage
			// message is no longer printed on error.
			cmd.SilenceUsage = true

			// TODO: Amazing work goes here!
			return nil
		},
	}

	cmd.Flags().IPVar(&flags.addr, "addr", nil, "addr description")
	cmd.Flags().Uint16Var(&flags.port, "port", 0, "port description")
	return cmd
}
