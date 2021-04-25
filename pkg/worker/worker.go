package worker

import (
	"context"
	"errors"
	"os"
	"os/exec"
	"sync"
	"syscall"

	"github.com/google/uuid"
)

// Command is a job request with the program name and arguments.
type Command struct {
	// Command name
	Name string
	// Command arguments
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
	Query(jobID string) (status Status, err error)
	// Streams the process output.
	//    - ctx: context to cancel the log stream
	//    - ID: Job identifier
	// It returns read chan to stream process stdout/stderr and the
	// execution error encountered.
	Stream(ctx context.Context, jobID string) (logchan chan string, err error)
}

// NewLinuxWorker creates a new Worker instance for Linux programs.
func NewWorker(config Config) Worker {
	return &worker{
		logger: NewLogger(config),
	}
}

// worker implementation.
type worker struct {
	logger Logger
	jobs   sync.Map
}

func (w *worker) Start(command Command) (jobID string, err error) {
	cmd := exec.Command(command.Name, command.Args...)
	jobID = uuid.NewString()
	logfile, err := w.logger.Create(jobID)
	if err != nil {
		return
	}
	// redirect the stdout and stderr to the log file
	cmd.Stdout = logfile
	cmd.Stderr = logfile
	if err = cmd.Start(); err != nil {
		return
	}
	// create and store the job
	job := w.createJob(jobID, cmd)
	// update the job status in background after
	// the process exit
	go w.updateJob(job.ID)
	return
}

func (w *worker) Stop(jobID string) (err error) {
	job, err := w.getJob(jobID)
	if err != nil {
		return
	}
	proc, err := os.FindProcess(job.Status.Pid)
	if err != nil {
		return
	}
	// kills the process gracefully
	return proc.Signal(syscall.SIGTERM)
}

func (w *worker) Query(jobID string) (status Status, err error) {
	job, err := w.getJob(jobID)
	if err != nil {
		return
	}
	status = *job.Status
	return
}

func (w *worker) Stream(ctx context.Context, jobID string) (logchan chan string, err error) {
	job, err := w.getJob(jobID)
	if err != nil {
		return
	}
	return w.logger.Tailf(ctx, job.ID)
}

func (w *worker) createJob(jobID string, cmd *exec.Cmd) *Job {
	job := &Job{ID: jobID, Cmd: cmd, Status: &Status{Pid: cmd.Process.Pid}}
	w.jobs.Store(jobID, job)
	return job
}

func (w *worker) updateJob(jobID string) {
	job, err := w.getJob(jobID)
	if err != nil {
		return
	}
	cmd := job.Cmd
	if err := cmd.Wait(); err == nil {
		// can be a good place to delete the log file
		// after process exit with success
		w.logger.Remove(jobID)
		// if the job is removed from map the query
		// will no longer find the job to return it's
		// status to the caller
		w.jobs.Delete(jobID)
	}
	// update the job status with information about
	// the exited process
	job.Status.ExitCode = cmd.ProcessState.ExitCode()
	job.Status.Exited = cmd.ProcessState.Exited()
}

func (w *worker) getJob(jobID string) (*Job, error) {
	value, ok := w.jobs.Load(jobID)
	if !ok {
		return nil, errors.New("Job not found")
	}
	return value.(*Job), nil
}
