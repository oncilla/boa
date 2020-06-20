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
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/oncilla/boa/pkg/boa"
)

func TestInit(t *testing.T) {
	var dir string
	if !*update {
		var err error
		dir, err = ioutil.TempDir("", "init")
		require.NoError(t, err)
		defer func() {
			require.NoError(t, os.RemoveAll(dir))
		}()
	} else {
		dir = "testdata/init"
		require.NoError(t, os.RemoveAll(dir))
		require.NoError(t, os.MkdirAll(dir, 0755))
	}

	cmd := newInit(boa.Pather("parent path"))
	cmd.SetArgs([]string{
		"--author", "my-name",
		"--license", "apache",
		"--path", dir,
		"server",
	})
	err := cmd.Execute()
	require.NoError(t, err)

	files, err := filepath.Glob("testdata/init/*")
	require.NoError(t, err)

	for _, file := range files {
		t.Log("Checking:", file)
		created, err := ioutil.ReadFile(filepath.Join(dir, filepath.Base(file)))
		require.NoError(t, err)
		golden, err := ioutil.ReadFile(file)
		require.NoError(t, err)
		assert.Equal(t, string(golden), string(created))
	}
}
