// Copyright 2026 Google Inc. All Rights Reserved.
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

package subcommands

import (
	"bytes"
	"context"
	"flag"
	"strings"
	"testing"
)

type flaggedCommand struct{}

func (flaggedCommand) Name() string             { return "flagged" }
func (flaggedCommand) Synopsis() string         { return "a command with one flag" }
func (flaggedCommand) Usage() string            { return "flagged [-known]\n" }
func (flaggedCommand) SetFlags(f *flag.FlagSet) { f.Bool("known", false, "a defined flag") }
func (flaggedCommand) Execute(context.Context, *flag.FlagSet, ...interface{}) ExitStatus {
	return ExitSuccess
}

// TestExecuteFlagErrorOutput checks that a command flag parse error goes to
// Commander.Error, not to os.Stderr (#38).
func TestExecuteFlagErrorOutput(t *testing.T) {
	top := flag.NewFlagSet("top", flag.ContinueOnError)
	cdr := NewCommander(top, "top")
	var stdout, stderr bytes.Buffer
	cdr.Output, cdr.Error = &stdout, &stderr
	cdr.Register(flaggedCommand{}, "")

	if err := top.Parse([]string{"flagged", "-unknown"}); err != nil {
		t.Fatalf("parsing top-level flags: %v", err)
	}
	if status := cdr.Execute(t.Context()); status != ExitUsageError {
		t.Errorf("Execute: got status %v, want %v", status, ExitUsageError)
	}
	const want = "flag provided but not defined: -unknown"
	if got := stderr.String(); !strings.Contains(got, want) {
		t.Errorf("Commander.Error: got %q, want it to contain %q", got, want)
	}
	if stdout.Len() != 0 {
		t.Errorf("Commander.Output: got %q, want nothing", stdout.String())
	}
}
