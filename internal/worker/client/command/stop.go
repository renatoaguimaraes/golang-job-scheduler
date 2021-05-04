package command

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/renatoaguimaraes/job-scheduler/internal/worker/proto"
	"google.golang.org/grpc"
)

type StopCommand struct {
	client proto.WorkerServiceClient
	args   []string
}

func NewStopCommand(client proto.WorkerServiceClient) Runner {
	return &StopCommand{
		client: client,
	}
}

func (c *StopCommand) Run(args []string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*100)
	defer cancel()
	command := proto.StopRequest{
		JobID: c.args[0],
	}
	_, err := c.client.Stop(ctx, &command, grpc.WaitForReady(true))
	if err != nil {
		return err
	}
	os.Stdout.WriteString(fmt.Sprintf("Job %v has been stopped\n", command.JobID))
	return nil
}
