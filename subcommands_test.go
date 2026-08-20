/*
Copyright 2024 Google Inc. All Rights Reserved.

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

package subcommands

import (
	"bytes"
	"context"
	"flag"
	"testing"
)

type testCmd struct {
	name     string
	synopsis string
	usage    string
	executed bool
}

func (c *testCmd) Name() string     { return c.name }
func (c *testCmd) Synopsis() string { return c.synopsis }
func (c *testCmd) Usage() string    { return c.usage }
func (c *testCmd) SetFlags(f *flag.FlagSet) {}
func (c *testCmd) Execute(ctx context.Context, f *flag.FlagSet, args ...interface{}) ExitStatus {
	c.executed = true
	return ExitSuccess
}

func TestCommander(t *testing.T) {
	fs := flag.NewFlagSet("test", flag.ContinueOnError)
	cdr := NewCommander(fs, "testapp")
	if cdr.Name() != "testapp" {
		t.Errorf("expected commander name %q, got %q", "testapp", cdr.Name())
	}

	cmd := &testCmd{
		name:     "foo",
		synopsis: "do foo",
		usage:    "foo usage",
	}
	cdr.Register(cmd, "groupA")

	var buf bytes.Buffer
	cdr.Output = &buf
	cdr.Error = &buf

	if err := fs.Parse([]string{"foo"}); err != nil {
		t.Fatalf("failed to parse top level flags: %v", err)
	}

	status := cdr.Execute(context.Background())
	if status != ExitSuccess {
		t.Errorf("expected ExitSuccess, got %v", status)
	}
	if !cmd.executed {
		t.Error("expected command to be executed")
	}
}

func TestHelpCommand(t *testing.T) {
	fs := flag.NewFlagSet("test", flag.ContinueOnError)
	cdr := NewCommander(fs, "testapp")
	cdr.Register(cdr.HelpCommand(), "")
	cdr.Register(cdr.FlagsCommand(), "")
	cdr.Register(cdr.CommandsCommand(), "")

	var buf bytes.Buffer
	cdr.Output = &buf
	cdr.Error = &buf

	if err := fs.Parse([]string{"help"}); err != nil {
		t.Fatalf("failed to parse: %v", err)
	}

	status := cdr.Execute(context.Background())
	if status != ExitSuccess {
		t.Errorf("expected ExitSuccess for help, got %v", status)
	}
	if buf.Len() == 0 {
		t.Error("expected non-empty help output")
	}
}

func TestVisitGroups(t *testing.T) {
	fs := flag.NewFlagSet("test", flag.ContinueOnError)
	cdr := NewCommander(fs, "testapp")
	cmd1 := &testCmd{name: "c1"}
	cmd2 := &testCmd{name: "c2"}
	cdr.Register(cmd1, "group1")
	cdr.Register(cmd2, "group2")

	var groups []string
	cdr.VisitGroups(func(g *CommandGroup) {
		groups = append(groups, g.Name())
	})

	if len(groups) != 2 || groups[0] != "group1" || groups[1] != "group2" {
		t.Errorf("unexpected groups: %v", groups)
	}
}

