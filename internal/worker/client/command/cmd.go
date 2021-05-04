package command

import (
	"errors"
	"fmt"

	"github.com/renatoaguimaraes/job-scheduler/internal/worker/client"
	"github.com/renatoaguimaraes/job-scheduler/pkg/worker/conf"
)

type Runner interface {
	// Run runs a initialized runner.
	Run(args []string) error
}

func Execute(config conf.Config, args []string) error {
	if len(args) < 1 {
		return errors.New("you must pass a command")
	}
	client, err := client.NewWorkerClient(config)
	if err != nil {
		return err
	}
	cmds := map[string]Runner{
		"start":  NewStartCommand(client),
		"query":  NewQueryCommand(client),
		"stop":   NewStopCommand(client),
		"stream": NewStreamCommand(client),
	}
	cmd, ok := cmds[args[0]]
	if ok {
		return cmd.Run(args[1:])
	}
	return fmt.Errorf("unknown command: %s", cmd)
}
