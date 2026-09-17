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

type dummyCmd struct {
	executed bool
	flagVal  string
}

func (d *dummyCmd) Name() string             { return "dummy" }
func (d *dummyCmd) Synopsis() string         { return "dummy command for test" }
func (d *dummyCmd) Usage() string            { return "dummy [-flag val]:\n" }
func (d *dummyCmd) SetFlags(f *flag.FlagSet) { f.StringVar(&d.flagVal, "flag", "", "a test flag") }
func (d *dummyCmd) Execute(ctx context.Context, f *flag.FlagSet, args ...interface{}) ExitStatus {
	d.executed = true
	return ExitSuccess
}

func TestExecuteFlagSetOutput(t *testing.T) {
	flags := flag.NewFlagSet("test", flag.ContinueOnError)
	cdr := NewCommander(flags, "test")
	cmd := &dummyCmd{}
	cdr.Register(cmd, "")

	var errBuf bytes.Buffer
	cdr.Error = &errBuf

	if err := flags.Parse([]string{"dummy", "-undefined_flag"}); err != nil {
		t.Fatalf("unexpected error parsing top-level flags: %v", err)
	}

	status := cdr.Execute(context.Background())
	if status != ExitUsageError {
		t.Errorf("expected ExitUsageError, got %v", status)
	}

	out := errBuf.String()
	if !strings.Contains(out, "flag provided but not defined: -undefined_flag") {
		t.Errorf("expected error output in cdr.Error, got: %q", out)
	}
}

func TestExecuteSuccess(t *testing.T) {
	flags := flag.NewFlagSet("test", flag.ContinueOnError)
	cdr := NewCommander(flags, "test")
	cmd := &dummyCmd{}
	cdr.Register(cmd, "")

	var errBuf bytes.Buffer
	cdr.Error = &errBuf

	if err := flags.Parse([]string{"dummy", "-flag", "hello"}); err != nil {
		t.Fatalf("unexpected error parsing top-level flags: %v", err)
	}

	status := cdr.Execute(context.Background())
	if status != ExitSuccess {
		t.Errorf("expected ExitSuccess, got %v", status)
	}
	if !cmd.executed {
		t.Errorf("expected dummyCmd to be executed")
	}
	if cmd.flagVal != "hello" {
		t.Errorf("expected flagVal to be 'hello', got %q", cmd.flagVal)
	}
}

func TestExecuteHelpFlag(t *testing.T) {
	flags := flag.NewFlagSet("test", flag.ContinueOnError)
	cdr := NewCommander(flags, "test")
	cmd := &dummyCmd{}
	cdr.Register(cmd, "")

	var errBuf bytes.Buffer
	cdr.Error = &errBuf

	if err := flags.Parse([]string{"dummy", "-h"}); err != nil {
		t.Fatalf("unexpected error parsing top-level flags: %v", err)
	}

	status := cdr.Execute(context.Background())
	if status != ExitSuccess {
		t.Errorf("expected ExitSuccess on -h, got %v", status)
	}
	if !strings.Contains(errBuf.String(), "dummy [-flag val]:") {
		t.Errorf("expected usage info in cdr.Error, got: %q", errBuf.String())
	}
}
