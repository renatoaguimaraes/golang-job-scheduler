package command

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/renatoaguimaraes/job-scheduler/internal/worker/proto"
	"google.golang.org/grpc"
)

type StartCommand struct {
	client proto.WorkerServiceClient
	args   []string
}

func NewStartCommand(client proto.WorkerServiceClient) Runner {
	return &StartCommand{
		client: client,
	}
}

func (c *StartCommand) Run(args []string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*100)
	defer cancel()
	command := proto.StartRequest{
		Name: args[0],
		Args: args[1:],
	}
	res, err := c.client.Start(ctx, &command, grpc.WaitForReady(true))
	if err != nil {
		return err
	}
	os.Stdout.WriteString(fmt.Sprintf("Job %v is started\n", res.JobID))
	return nil
}
