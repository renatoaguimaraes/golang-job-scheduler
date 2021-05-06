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

type StopCommand struct {
	client proto.WorkerServiceClient
}

func NewStopCommand(client proto.WorkerServiceClient) Runner {
	return &StopCommand{
		client: client,
	}
}

func (c *StopCommand) Run(args []string) error {
	if len(args) < 1 {
		return errors.New("you must pass an argument")
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	command := proto.StopRequest{
		JobID: args[0],
	}
	_, err := c.client.Stop(ctx, &command, grpc.WaitForReady(true))
	if err != nil {
		return err
	}
	os.Stdout.WriteString(fmt.Sprintf("Job %v has been stopped\n", command.JobID))
	return nil
}
