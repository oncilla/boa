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
	"net"
	"os"
	"strings"

	"github.com/mitchellh/mapstructure"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"

	"github.com/oncilla/boa/pkg/boa"
	"github.com/oncilla/boa/pkg/boa/flag"
)

func defaultConfig() *Config {
	return &Config{
		DB: DB{
			User:     "user",
			Password: "password",
		},
		Addr: &flag.TCPAddr{IP: net.ParseIP("127.0.0.1"), Port: 8080},
	}
}

func main() {
	v := viper.New()

	cmd := &cobra.Command{
		Use:   "config <config-file>",
		Short: "A sample application with config parsing",
		Long: `This is a sample application that showcases config parsign with the help of boa.

This application loads the configuration based on the following precedence:

1. Command line flag
2. Environment variable
3. Configuration file
4. Default value

The default configuration is:

    db:
	  user: user
	  password: password
	addr: 127.0.0.1:8080

The config file can be passed as as a command line argument.

Environment variables are prefixed With 'SAMPLE_':

SAMPLE_DB_USER=secure
`,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg := defaultConfig()
			v.SetEnvPrefix("sample")
			v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
			v.BindPFlags(cmd.Flags())

			// Set the default values for the configuration.
			if err := boa.SetDefaults(v, cfg); err != nil {
				return err
			}

			// Bind the environment variables based on the configuration struct.
			if err := boa.BindEnv(v, cfg); err != nil {
				return err
			}

			// Set hooks to parse the TCP address.
			hooks := viper.DecodeHook(mapstructure.ComposeDecodeHookFunc(boa.DefaultDecodeHooks()...))

			// If any configuration files are passed, parse them.
			if len(args) > 0 {
				for _, file := range args {
					v.SetConfigFile(file)
				}

				if err := v.ReadInConfig(); err != nil {
					return err
				}
			}

			// Pares the configuration.
			var out Config
			if err := v.Unmarshal(&out, hooks); err != nil {
				return err
			}

			// Display the parsed configuration.
			enc := yaml.NewEncoder(os.Stdout)
			return enc.Encode(out)
		},
	}
	if err := boa.AddFlags(cmd.Flags(), defaultConfig()); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		os.Exit(1)
	}
	if err := cmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		os.Exit(1)
	}
}

type Config struct {
	DB   DB            `mapstructure:"db"`
	Addr *flag.TCPAddr `mapstructure:"addr"`
}

type DB struct {
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
}
