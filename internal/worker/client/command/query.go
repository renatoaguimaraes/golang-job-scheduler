package command

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/renatoaguimaraes/job-scheduler/internal/worker/proto"
	"google.golang.org/grpc"
)

type QueryCommand struct {
	client proto.WorkerServiceClient
	args   []string
}

func NewQueryCommand(client proto.WorkerServiceClient) Runner {
	return &QueryCommand{
		client: client,
	}
}

func (c *QueryCommand) Init(args []string) {
	c.args = args
}

func (c *QueryCommand) Run() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*100)
	defer cancel()
	command := proto.QueryRequest{
		JobID: c.args[0],
	}
	res, err := c.client.Query(ctx, &command, grpc.WaitForReady(true))
	if err == nil {
		os.Stdout.WriteString(fmt.Sprintf("Pid: %v Exit code: %v Exited: %v\n", res.Pid, res.ExitCode, res.Exited))
	} else {
		os.Stderr.WriteString(fmt.Sprintf("%v\n", err))
	}
}
