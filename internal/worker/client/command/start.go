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

func (c *StartCommand) Init(args []string) {
	c.args = args
}

func (c *StartCommand) Run() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*100)
	defer cancel()
	command := proto.StartRequest{
		Name: c.args[0],
		Args: c.args[1:],
	}
	res, err := c.client.Start(ctx, &command, grpc.WaitForReady(true))
	if err == nil {
		os.Stdout.WriteString(fmt.Sprintf("Job %v is started\n", res.JobID))
	} else {
		os.Stderr.WriteString(fmt.Sprintf("%v\n", err))
	}
}
