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

package boa_test

import (
	"net"
	"os"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/oncilla/boa/pkg/boa"
	"github.com/oncilla/boa/pkg/boa/mock_boa"
)

type Config struct {
	DB struct {
		User     string `mapstructure:"user"`
		Password string `mapstructure:"password"`
	} `mapstructure:"db"`
	Addr  *net.TCPAddr `mapstructure:"addr"`
	Token `mapstructure:",squash"`
}

type Token struct {
	Token string `mapstructure:"token"`
}

func TestBindEnv(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	r := mock_boa.NewMockConfigRegistry(ctrl)
	r.EXPECT().BindEnv([]string{"db.user"})
	r.EXPECT().BindEnv([]string{"db.password"})
	r.EXPECT().BindEnv([]string{"addr"})
	r.EXPECT().BindEnv([]string{"token"})

	var config Config
	err := boa.BindEnv(r, &config)
	require.NoError(t, err)
}

func TestSetDefault(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	var config Config
	config.DB.User = "oncilla"
	config.DB.Password = "password"
	config.Addr = &net.TCPAddr{IP: net.ParseIP("127.0.0.1"), Port: 8080}
	config.Token.Token = "token"

	r := mock_boa.NewMockConfigRegistry(ctrl)
	r.EXPECT().SetDefault("db.user", "oncilla")
	r.EXPECT().SetDefault("db.password", "password")
	r.EXPECT().SetDefault("addr", &net.TCPAddr{IP: net.ParseIP("127.0.0.1"), Port: 8080})
	r.EXPECT().SetDefault("token", "token")

	err := boa.SetDefaults(r, &config)
	require.NoError(t, err)
}

func TestViperBindEnv(t *testing.T) {
	var config Config
	v := viper.New()
	v.SetEnvPrefix("boa")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	err := boa.BindEnv(v, &config)
	require.NoError(t, err)

	os.Setenv("BOA_DB_USER", "oncilla")
	os.Setenv("BOA_DB_PASSWORD", "password")
	os.Setenv("BOA_ADDR", "127.0.0.1:8080")
	os.Setenv("BOA_TOKEN", "token")

	err = v.Unmarshal(&config, viper.DecodeHook(
		mapstructure.ComposeDecodeHookFunc(
			boa.DefaultDecodeHooks()...,
		),
	))
	require.NoError(t, err)

	assert.Equal(t, "oncilla", config.DB.User)
	assert.Equal(t, "password", config.DB.Password)
	assert.Equal(t, "127.0.0.1:8080", config.Addr.String())
	assert.Equal(t, "token", config.Token.Token)
}
