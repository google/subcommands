package subcommands

import (
	"bytes"
	"context"
	"flag"
	"strings"
	"testing"
)

type dummyCommand struct{}

func (d *dummyCommand) Name() string           { return "dummy" }
func (d *dummyCommand) Synopsis() string       { return "dummy command" }
func (d *dummyCommand) SetFlags(f *flag.FlagSet) {
	f.Bool("foo", false, "foo flag")
}
func (d *dummyCommand) Usage() string          { return "dummy usage\n" }
func (d *dummyCommand) Execute(ctx context.Context, f *flag.FlagSet, args ...interface{}) ExitStatus {
	return ExitSuccess
}

func TestCommanderFlagOutput(t *testing.T) {
	var out, errOut bytes.Buffer
	topFlags := flag.NewFlagSet("test", flag.ContinueOnError)
	topFlags.Parse([]string{"dummy", "--bar"})
	cdr := NewCommander(topFlags, "test")
	cdr.Output = &out
	cdr.Error = &errOut
	cdr.Register(&dummyCommand{}, "")

	status := cdr.Execute(context.Background())
	if status != ExitUsageError {
		t.Fatalf("expected ExitUsageError, got %v", status)
	}

	errOutput := errOut.String()
	if !strings.Contains(errOutput, "flag provided but not defined: -bar") {
		t.Errorf("expected flag error output to be written to commander error writer, got: %q", errOutput)
	}
}
