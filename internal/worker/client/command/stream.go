package command

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"github.com/renatoaguimaraes/job-scheduler/internal/worker/proto"
	"google.golang.org/grpc"
)

type StreamCommand struct {
	client proto.WorkerServiceClient
	args   []string
}

func NewStreamCommand(client proto.WorkerServiceClient) Runner {
	return &StreamCommand{
		client: client,
	}
}

func (c *StreamCommand) Run(args []string) error {
	ctx, cancel := context.WithCancel(context.Background())
	command := proto.StreamRequest{
		JobID: c.args[0],
	}
	res, err := c.client.Stream(ctx, &command, grpc.WaitForReady(true))
	if err != nil {
		cancel()
		return err
	}
	// runs the streaming in backgroud
	go func() {
		for {
			out, err := res.Recv()
			if err != nil {
				os.Stderr.WriteString(fmt.Sprintf("%v\n", err))
				return
			}
			os.Stdout.WriteString(out.Output)
		}
	}()
	// waits for os signal to terminate the streaming
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, os.Interrupt)
	defer func() {
		cancel()
		signal.Stop(sigchan)
	}()
	<-sigchan
	return nil
}
