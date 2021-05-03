package api

import (
	"context"

	"github.com/renatoaguimaraes/job-scheduler/internal/worker/proto"
	"github.com/renatoaguimaraes/job-scheduler/pkg/worker"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type workerServer struct {
	proto.UnimplementedWorkerServiceServer
	Worker worker.Worker
}

func (s *workerServer) Start(ctx context.Context, r *proto.StartRequest) (*proto.StartResponse, error) {
	jobID, err := s.Worker.Start(worker.Command{Name: r.Name, Args: r.Args})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	res := proto.StartResponse{
		JobID: jobID,
	}
	return &res, nil
}

func (s *workerServer) Stop(ctx context.Context, r *proto.StopRequest) (*proto.StopResponse, error) {
	err := s.Worker.Stop(r.JobID)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &proto.StopResponse{}, nil
}

func (s *workerServer) Query(ctx context.Context, r *proto.QueryRequest) (*proto.QueryResponse, error) {
	jobstatus, err := s.Worker.Query(r.JobID)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	res := proto.QueryResponse{
		Pid:      int32(jobstatus.Pid),
		ExitCode: int32(jobstatus.ExitCode),
		Exited:   jobstatus.Exited,
	}
	return &res, nil
}

func (s *workerServer) Stream(r *proto.StreamRequest, stream proto.WorkerService_StreamServer) error {
	logchan, err := s.Worker.Stream(stream.Context(), r.JobID)
	if err != nil {
		return status.Error(codes.Internal, err.Error())
	}
	for {
		select {
		case <-stream.Context().Done():
			return stream.Context().Err()
		case log, ok := <-logchan:
			if !ok {
				return status.Error(codes.Internal, "the log channel was closed")
			}
			if err := stream.SendMsg(&proto.StreamResponse{Output: log}); err != nil {
				return status.Error(codes.Internal, err.Error())
			}
		}
	}
}
