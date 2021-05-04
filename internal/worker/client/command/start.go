package command

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/renatoaguimaraes/job-scheduler/internal/worker/proto"
	"google.golang.org/grpc"
)

type StartCommand struct {
	client proto.WorkerServiceClient
}

func NewStartCommand(client proto.WorkerServiceClient) Runner {
	return &StartCommand{
		client: client,
	}
}

func (c *StartCommand) Run(args []string) error {
	if len(args) < 1 {
		return errors.New("you must pass a program name")
	}
	var cargs []string
	if len(args) > 1 {
		cargs = append(cargs, args[1:]...)
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	command := proto.StartRequest{
		Name: args[0],
		Args: cargs,
	}
	res, err := c.client.Start(ctx, &command, grpc.WaitForReady(true))
	if err != nil {
		return err
	}
	os.Stdout.WriteString(fmt.Sprintf("Job %v is started\n", res.JobID))
	return nil
}
