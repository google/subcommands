/*
Copyright 2018 Google Inc. All Rights Reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"bytes"
	"flag"
	"io/ioutil"
	"testing"

	"github.com/google/go-cmp/cmp"
)

var update = flag.Bool("update", false, "update .golden files")

func Test(t *testing.T) {
	tests := []struct {
		desc       string
		goldenFile string
		params     tmplParams
	}{{
		"fooCmd",
		"testdata/foo.golden",
		tmplParams{
			Cmd:      "foo",
			Pkg:      "main",
			Synopsis: "lorem ipsum dolor",
			Usage: `usage: foo [flags]

wopr foo -q=42`,
			Username: "Alice",
		},
	}, {
		"BarCmd",
		"testdata/bar.golden",
		tmplParams{
			Cmd:      "Bar",
			Pkg:      "bar",
			Synopsis: "sit amet, consectetur",
			Usage: `usage: bar [flags]

wopr bar -x`,
			Username: "Bob",
		},
	}}

	for _, test := range tests {
		var got bytes.Buffer

		if err := tmpl.Execute(&got, test.params); err != nil {
			t.Errorf("%s: failed to execute template: %v", test.desc, err)
			continue
		}

		if *update {
			if err := ioutil.WriteFile(test.goldenFile, got.Bytes(), 0644); err != nil {
				t.Errorf("%s: failed to update golden file (%s): %v", test.desc, test.goldenFile, err)
			}
			continue
		}

		want, err := ioutil.ReadFile(test.goldenFile)
		if err != nil {
			t.Errorf("%s: failed to open golden file (%s): %v", test.desc, test.goldenFile, err)
			continue
		}

		if diff := cmp.Diff(got.String(), string(want)); diff != "" {
			t.Errorf("%s: files differ (-got +want)\n%s", test.desc, diff)
		}
	}
}
