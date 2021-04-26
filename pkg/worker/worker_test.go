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

	assert.NotEmpty(t, jobID)
	assert.Nil(t, err)
}

func TestStartNotExistingCommand(t *testing.T) {
	jobID, err := w.Start(Command{Name: "notexists"})

	assert.NotEmpty(t, jobID)
	assert.NotNil(t, err)
}

func TestStopNotExistingProcess(t *testing.T) {
	err := w.Stop("")
	assert.Equal(t, "Job not found", err.Error(), "")
}
func TestStopExistingProcess(t *testing.T) {
	jobID, err := w.Start(Command{Name: "sleep", Args: []string{"1"}})
	assert.NoError(t, err)

	err = w.Stop(jobID)
	assert.NoError(t, err)
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
	jobID, err := w.Start(Command{Name: "ls", Args: []string{"1"}})
	assert.NoError(t, err)

	err = w.Stop(jobID)
	assert.NoError(t, err)

	time.Sleep(time.Second * 2)

	st, err := w.Query(jobID)
	assert.NoError(t, err)
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

	err = w.Stop(jobID)
	assert.NoError(t, err)
}
