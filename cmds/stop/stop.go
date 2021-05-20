package stop

import (
	"context"
	"flag"

	"github.com/google/subcommands"
	"github.com/keyneston/fscachemonitor/internal/shared"
)

type Command struct {
	*shared.Config

	root     string
	filename string
	sql      bool
	dirOnly  bool
}

func (*Command) Name() string     { return "stop" }
func (*Command) Synopsis() string { return "Stop running fscachemonitor" }
func (*Command) Usage() string {
	return `stop:
`
}

func (c *Command) SetFlags(f *flag.FlagSet) {
	c.Config.SetFlags(f)

	f.StringVar(&c.root, "r", "", "Root directory to monitor")
	f.StringVar(&c.root, "root", "", "Alias for -r")
	f.StringVar(&c.filename, "c", "", "File to output cache to")
	f.StringVar(&c.filename, "cache", "", "Alias for -c")
	f.BoolVar(&c.sql, "s", false, "Use SQLite3 backed file")
}

func (c *Command) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	if c.root == "" {
		return shared.Exitf("Must specify root to watch")
	}

	if c.filename == "" {
		return shared.Exitf("Must specify file to output cache to")
	}

	pid, err := shared.NewPID(c.PIDFile, c.root, c.filename)
	if err != nil {
		return shared.Exitf("Error creating pid file: %v", err)
	}

	defer pid.Release()
	if ok, err := pid.Acquire(); err != nil {
		return shared.Exitf("Error checking pid: %v", err)
	} else if !ok {
		if err := pid.Stop(); err != nil {
			shared.Exitf("Error stopping: %v", err)
		}
	}

	return subcommands.ExitSuccess
}
