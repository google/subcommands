/*
Copyright 2016 Google Inc. All Rights Reserved.

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

type dummyCmd struct {
	name      string
	executed  bool
	tag       string
}

func (d *dummyCmd) Name() string     { return d.name }
func (d *dummyCmd) Synopsis() string { return "dummy synopsis" }
func (d *dummyCmd) Usage() string    { return "dummy usage" }
func (d *dummyCmd) SetFlags(f *flag.FlagSet) {}
func (d *dummyCmd) Execute(ctx context.Context, f *flag.FlagSet, args ...interface{}) ExitStatus {
	d.executed = true
	return ExitSuccess
}

type taggedDummyCmd struct {
	name string
	tag  string
}

func (d *taggedDummyCmd) Name() string     { return d.name }
func (d *taggedDummyCmd) Synopsis() string { return d.tag }
func (d *taggedDummyCmd) Usage() string    { return d.tag }
func (d *taggedDummyCmd) SetFlags(f *flag.FlagSet) {}
func (d *taggedDummyCmd) Execute(ctx context.Context, f *flag.FlagSet, args ...interface{}) ExitStatus {
	return ExitSuccess
}

func TestRegisterHomonymousCommands(t *testing.T) {
	cdr := NewCommander(flag.NewFlagSet("test", flag.ContinueOnError), "test")

	cmd1 := &dummyCmd{name: "foo", tag: "first"}
	cmd2 := &dummyCmd{name: "foo", tag: "second"}

	cdr.Register(cmd1, "group1")
	cdr.Register(cmd2, "group1")

	// Verify that only one command exists and it's cmd2
	count := 0
	cdr.VisitCommands(func(g *CommandGroup, c Command) {
		count++
		if g.Name() != "group1" {
			t.Errorf("expected group 'group1', got %q", g.Name())
		}
		if c != cmd2 {
			t.Errorf("expected registered command to be cmd2, got %v", c)
		}
	})
	if count != 1 {
		t.Errorf("expected exactly 1 command, got %d", count)
	}
}

func TestRegisterHomonymousCommandsCrossGroup(t *testing.T) {
	cdr := NewCommander(flag.NewFlagSet("test", flag.ContinueOnError), "test")

	cmd1 := &dummyCmd{name: "foo"}
	cmd2 := &dummyCmd{name: "foo"}

	cdr.Register(cmd1, "group1")
	cdr.Register(cmd2, "group2")

	// Verify that cmd1 is removed from group1 and cmd2 is in group2
	groupCount := 0
	cdr.VisitGroups(func(g *CommandGroup) {
		groupCount++
		if g.Name() == "group1" {
			if len(g.commands) != 0 {
				t.Errorf("expected group1 to be empty, got %d commands", len(g.commands))
			}
		} else if g.Name() == "group2" {
			if len(g.commands) != 1 || g.commands[0] != cmd2 {
				t.Errorf("expected group2 to contain cmd2")
			}
		} else {
			t.Errorf("unexpected group %q", g.Name())
		}
	})
}

func TestExecuteHomonymousCommand(t *testing.T) {
	cdr := NewCommander(flag.NewFlagSet("test", flag.ContinueOnError), "test")

	cmd1 := &dummyCmd{name: "foo"}
	cmd2 := &dummyCmd{name: "foo"}

	cdr.Register(cmd1, "")
	cdr.Register(cmd2, "")

	ctx := context.Background()
	// Set args to execute "foo"
	cdr.topFlags.Parse([]string{"foo"})
	status := cdr.Execute(ctx)
	if status != ExitSuccess {
		t.Errorf("expected ExitSuccess, got %v", status)
	}

	if cmd1.executed {
		t.Errorf("cmd1 should not have been executed")
	}
	if !cmd2.executed {
		t.Errorf("cmd2 should have been executed")
	}
}

func TestHelpReflectsUpdatedCommand(t *testing.T) {
	cdr := NewCommander(flag.NewFlagSet("test", flag.ContinueOnError), "test")
	buf := &bytes.Buffer{}
	cdr.Output = buf
	cdr.Error = buf

	cmd1 := &taggedDummyCmd{name: "foo", tag: "old-synopsis"}
	cmd2 := &taggedDummyCmd{name: "foo", tag: "new-synopsis"}

	cdr.Register(cmd1, "group1")
	cdr.Register(cmd2, "group1")

	ctx := context.Background()
	cdr.topFlags.Parse([]string{"help", "foo"})
	cdr.Execute(ctx)

	output := buf.String()
	if bytes.Contains([]byte(output), []byte("old-synopsis")) {
		t.Errorf("help output contains stale 'old-synopsis': %s", output)
	}
	if !bytes.Contains([]byte(output), []byte("new-synopsis")) {
		t.Errorf("help output missing 'new-synopsis': %s", output)
	}
}
