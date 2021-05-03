package worker

import (
	"context"
	"errors"
	"fmt"
	logger "log"
	"os/exec"
	"sync"
	"syscall"

	"github.com/google/uuid"
	"github.com/renatoaguimaraes/job-scheduler/pkg/worker/conf"
	"github.com/renatoaguimaraes/job-scheduler/pkg/worker/log"
)

// Command is a job request with the program name and arguments.
type Command struct {
	// Name program path/name
	Name string
	// Args program arguments
	Args []string
}

// Job represents an arbitrary Linux process schedule by the Worker.
type Job struct {
	// ID job identifier
	ID string
	// Command pipeline
	Cmd *exec.Cmd
	// Status of the process.
	Status *Status
}

// IsRunning checks if the process still running.
func (j *Job) IsRunning() bool {
	return j.Status.ExitCode != 0 || !j.Status.Exited
}

// Status of the process.
type Status struct {
	// Process identifier
	Pid int
	// ExitCode of the exited process, or -1 if the process hasn't
	// exited or was terminated by a signal
	ExitCode int
	// Exited reports whether the program has exited
	Exited bool
}

// Worker defines the basic operations to manage Jobs.
type Worker interface {
	// Start creates a Linux process.
	//    - command: command to be executed
	// It returns the job ID and the execution error encountered.
	Start(command Command) (jobID string, err error)
	// Stop a running Job which kills a running process.
	//    - ID: Job identifier
	// It returns the execution error encountered.
	Stop(jobID string) (err error)
	// Query a Job to check the current status.
	//    - ID: Job identifier
	// It returns process status and the execution error
	// encountered.
	Query(jobID string) (status *Status, err error)
	// Streams the process output.
	//    - ctx: context to cancel the log stream
	//    - ID: Job identifier
	// It returns read chan to stream process stdout/stderr and the
	// execution error encountered.
	Stream(ctx context.Context, jobID string) (logchan chan string, err error)
}

// NewWorker creates a new Worker instance.
func NewWorker(config conf.Config) Worker {
	return &worker{
		logger: log.NewLogger(config),
		jobs:   make(map[string]*Job),
	}
}

// worker implementation.
type worker struct {
	// logger is responsible to handle the
	// stdout and stderr of a running process
	logger log.Logger
	// jobs is concurrency safe map to store
	// the requested jobs
	jobs map[string]*Job
	// mtx to control jobs concurret access
	mtx sync.RWMutex
}

// Start runs a Linux single command with arguments.
// If the command runs with success, a Job identifier will be returned.
// A log file will be created with the Job id name to capture the stdout and stderr
// of a running process.
// To get the process status, the Job request will be stored in memory,
// and a goroutine will be launched to update the job status when the process is finished.
func (w *worker) Start(command Command) (string, error) {
	cmd := exec.Command(command.Name, command.Args...)
	jobID := uuid.NewString()
	logfile, err := w.logger.Create(jobID)
	if err != nil {
		return jobID, err
	}
	// redirect the stdout and stderr to the log file
	cmd.Stdout = logfile
	cmd.Stderr = logfile
	if err = cmd.Start(); err != nil {
		w.logger.Remove(jobID)
		return jobID, err
	}
	// create and store the job
	w.mtx.Lock()
	defer w.mtx.Unlock()
	job := Job{ID: jobID, Cmd: cmd, Status: &Status{Pid: cmd.Process.Pid}}
	w.jobs[jobID] = &job
	// update the job status in background
	go func() {
		if err := job.Cmd.Wait(); err != nil {
			logger.Printf("Command execution fails, %v", err)
		}
		// update the job status with information about
		// the exited process
		status := Status{
			Pid:      job.Cmd.ProcessState.Pid(),
			ExitCode: job.Cmd.ProcessState.ExitCode(),
			Exited:   job.Cmd.ProcessState.Exited(),
		}
		w.mtx.Lock()
		job.Status = &status
		w.mtx.Unlock()
	}()
	return jobID, nil
}

// Stop terminates a running Job gracefully sending a SIGTERM to the process.
// If the job doesn't exitis an error will be returned.
func (w *worker) Stop(jobID string) error {
	w.mtx.RLock()
	defer w.mtx.RUnlock()
	job, err := w.getJob(jobID)
	if err != nil {
		return err
	}
	if job.IsRunning() {
		return job.Cmd.Process.Signal(syscall.SIGTERM)
	}
	return errors.New("the process is already finished")
}

// Query returns the process status of a specific Job.
// If the job doesn't exitis an error will be returned.
func (w *worker) Query(jobID string) (Status, error) {
	w.mtx.RLock()
	defer w.mtx.RUnlock()
	job, err := w.getJob(jobID)
	if err != nil {
		return nil, err
	}
	return *job.Status, nil
}

// Stream reads from the log file, like 'tail -f' through
// a channel. If the context is canceled the channel will
// be closed and the tailing will be stopped.
func (w *worker) Stream(ctx context.Context, jobID string) (chan string, error) {
	w.mtx.RLock()
	job, err := w.getJob(jobID)
	w.mtx.RUnlock()
	if err != nil {
		return nil, err
	}
	return w.logger.Tailf(ctx, job.ID)
}

// getJob helper to get a job given an id.
func (w *worker) getJob(jobID string) (*Job, error) {
	job, ok := w.jobs[jobID]
	if !ok {
		return nil, fmt.Errorf("Job %v not found", jobID)
	}
	return job, nil
}
