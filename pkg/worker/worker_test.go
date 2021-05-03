package worker

import (
	"context"
	"testing"
	"time"

	"github.com/renatoaguimaraes/job-scheduler/pkg/worker/conf"
	"github.com/stretchr/testify/assert"
)

var w = NewWorker(conf.NewConfig())

func TestStartExistingCommand(t *testing.T) {
	jobID, err := w.Start(Command{Name: "ls"})

	assert.NotEmpty(t, jobID)
	assert.Nil(t, err)
}

func TestStartNotExistingCommand(t *testing.T) {
	jobID, err := w.Start(Command{Name: "notexists"})

	assert.NotEmpty(t, jobID)
	assert.NotNil(t, err)
}

func TestStopNotExistingProcess(t *testing.T) {
	err := w.Stop("notexists")
	assert.Equal(t, "Job notexists not found", err.Error())
}

func TestStopExistingProcess(t *testing.T) {
	jobID, err := w.Start(Command{Name: "sleep", Args: []string{"2"}})
	assert.NoError(t, err)

	err = w.Stop(jobID)
	assert.NoError(t, err)
}

func TestStopStoppedProcess(t *testing.T) {
	jobID, err := w.Start(Command{Name: "sleep", Args: []string{"1"}})
	assert.NoError(t, err)

	time.Sleep(time.Second * 2)

	err = w.Stop(jobID)
	assert.Error(t, err)
}

func TestQueryExistingProcess(t *testing.T) {
	jobID, err := w.Start(Command{Name: "sleep", Args: []string{"1"}})
	assert.NoError(t, err)

	st, err := w.Query(jobID)
	assert.False(t, st.Exited)
	assert.Zero(t, st.ExitCode)
	assert.NoError(t, err)

	time.Sleep(time.Second * 2)

	st, err = w.Query(jobID)
	assert.NoError(t, err)
	assert.True(t, st.Exited)
	assert.Zero(t, st.ExitCode)
}

func TestQueryStoppedProcess(t *testing.T) {
	jobID, err := w.Start(Command{Name: "sleep", Args: []string{"1"}})
	assert.NoError(t, err)

	err = w.Stop(jobID)
	assert.NoError(t, err)

	time.Sleep(time.Second * 2)

	st, err := w.Query(jobID)
	assert.NoError(t, err)
	assert.False(t, st.Exited)
	assert.Equal(t, -1, st.ExitCode)
}

func TestQueryNotExistingProcess(t *testing.T) {
	status, err := w.Query("not-exists-job-id")

	assert.Error(t, err)
	assert.Equal(t, Status{}, status)
}

func TestStreamExistingProcess(t *testing.T) {
	jobID, err := w.Start(Command{Name: "bash", Args: []string{"-c", "while true; do date; sleep 1; done"}})
	assert.Nil(t, err, "err should be nil")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	logchan, err := w.Stream(ctx, jobID)
	assert.Nil(t, err, "err should be nil")
	assert.NotNil(t, <-logchan)
	cancel()

	err = w.Stop(jobID)
	assert.NoError(t, err)
}

func TestStreamNotExistingProcess(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	logchan, err := w.Stream(ctx, "not-exists-job-id")
	cancel()
	assert.Nil(t, logchan)
	assert.Error(t, err)
}
