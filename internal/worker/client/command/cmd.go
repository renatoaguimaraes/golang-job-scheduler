package command

import (
	"fmt"
	"os"

	"github.com/renatoaguimaraes/job-scheduler/internal/worker/client"
	"github.com/renatoaguimaraes/job-scheduler/pkg/worker/conf"
)

type Runner interface {
	// Init initialize a runner with execution arguments.
	//    - args: runner arguments
	Init(args []string)
	// Run runs a initialized runner.
	Run()
}

func Execute(config conf.Config, args []string) {
	if len(args) < 1 {
		os.Stderr.WriteString("you must pass a command")
		return
	}
	c := client.NewWorkerClient(config)
	cmds := map[string]Runner{
		"start":  NewStartCommand(c),
		"query":  NewQueryCommand(c),
		"stop":   NewStopCommand(c),
		"stream": NewStreamCommand(c),
	}
	cmd, ok := cmds[args[0]]
	if ok {
		cmd.Init(args[1:])
		cmd.Run()
		return
	}
	os.Stderr.WriteString(fmt.Sprintf("unknown command: %s", cmd))
}
