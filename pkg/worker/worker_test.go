package worker

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var w = NewWorker(NewConfig())

func TestStartExistingCommand(t *testing.T) {
	jobID, err := w.Start(Command{Name: "ls"})

	assert.NotEmpty(t, jobID, "jobID shouldn't be empty")
	assert.Nil(t, err, "err should be nil")
}

func TestStartNotExisingCommand(t *testing.T) {
	jobID, err := w.Start(Command{Name: "xpto"})

	assert.NotEmpty(t, jobID, "jobID shouldn't be empty")
	assert.NotNil(t, err, "err shouldn't be nil")
}

func TestStopProcess(t *testing.T) {
	err := w.Stop("")

	assert.Equal(t, "Job not found", err.Error(), "")
}
func TestStopExistingProcess(t *testing.T) {
	jobID, _ := w.Start(Command{Name: "sleep", Args: []string{"2m"}})
	err := w.Stop(jobID)

	assert.NoError(t, err, "err should be nil, process can't stop")
}

func TestQueryExistingProcess(t *testing.T) {
	jobID, _ := w.Start(Command{Name: "sleep", Args: []string{"2"}})
	st, err := w.Query(jobID)

	assert.NotNil(t, st)
	assert.False(t, st.Exited)
	assert.Zero(t, st.ExitCode)
	assert.NoError(t, err, "err should be nil, process can't stop")

	time.Sleep(time.Second * 3)
	st, err = w.Query(jobID)

	assert.NotNil(t, st)
	assert.NoError(t, err, "err should be nil, process can't stop")
	assert.True(t, st.Exited)
	assert.Zero(t, st.ExitCode)
}

func TestQueryStoppedProcess(t *testing.T) {
	jobID, _ := w.Start(Command{Name: "sleep", Args: []string{"1"}})
	err := w.Stop(jobID)

	assert.NoError(t, err, "err should be nil, process can't stop")

	time.Sleep(time.Second * 2)
	st, err := w.Query(jobID)

	assert.NotNil(t, st)
	assert.NoError(t, err, "err should be nil, process can't stop")
	assert.False(t, st.Exited)
	assert.Equal(t, -1, st.ExitCode)
}

func TestStreamExistingCommand(t *testing.T) {
	jobID, err := w.Start(Command{Name: "top"})
	assert.Nil(t, err, "err should be nil")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*1)
	logchan, err := w.Stream(ctx, jobID)

	assert.Nil(t, err, "err should be nil")
	assert.NotNil(t, <-logchan)

	cancel()
	w.Stop(jobID)
}
