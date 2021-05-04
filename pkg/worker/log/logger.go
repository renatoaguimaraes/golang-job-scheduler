package log

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/renatoaguimaraes/job-scheduler/pkg/worker/conf"
)

// Logger implementation.
type Logger struct {
	config  conf.Config
	watcher Watcher
}

// NewLogger returns a new Logger instance.
func NewLogger(config conf.Config) Logger {
	return Logger{
		config:  config,
		watcher: NewWatcher(),
	}
}

// Path returns a absolute file path given a name.
func (l *Logger) Path(name string) string {
	return filepath.Join(l.config.LogFolder, fmt.Sprintf("%s.log", name))
}

// Create creates and return a os.File under log folder.
// If the file can't be created an error will be returned.
func (l *Logger) Create(name string) (*os.File, error) {
	return os.Create(l.Path(name))
}

// Remove removes the named file under the log folder.
func (l *Logger) Remove(name string) error {
	return os.Remove(l.Path(name))
}

// Tailf watch a named log file under the log folder, and
// streams his content through a channel.
func (l *Logger) Tailf(ctx context.Context, name string) (chan string, error) {
	file, err := os.OpenFile(l.Path(name), os.O_RDONLY, 0644)
	if err != nil {
		return nil, err
	}
	logchan := make(chan string)
	go func() {
		defer func() {
			if err := file.Close(); err != nil {
				log.Printf("fail to close the log file: %v", err)
			}
			close(logchan)
		}()
		// reads file from the begin
		if err := l.streamFile(ctx, file, logchan); err != nil && err != io.EOF {
			log.Printf("fail to read the log file: %v", err)
			return
		}
		// watching modify and close events
		eventchan, err := l.watcher.Watch(ctx, l.Path(name))
		if err != nil {
			log.Printf("fail to watch the log file events: %v", err)
			return
		}
		// reads file changes
		for {
			if err := waitForChange(ctx, eventchan); err != nil {
				log.Printf("%v", err)
				return
			}
			if err := l.streamFile(ctx, file, logchan); err != nil && err != io.EOF {
				log.Printf("fail to read the log file: %v", err)
				return
			}
		}
	}()
	return logchan, nil
}

// streamFile reads chunks from log file given a specific offset and send them throught the channel.
func (l *Logger) streamFile(ctx context.Context, file *os.File, logchan chan string) error {
	for {
		chunck := make([]byte, l.config.LogChunckSize)
		nbytes, err := file.Read(chunck)
		if err != nil {
			return err
		}
		select {
		case logchan <- string(chunck[:nbytes]):
		case <-ctx.Done():
			return errors.New("log file stream cancelled")
		}
	}
}

// waitForChange waits for file system change events and them through the channel.
func waitForChange(ctx context.Context, eventchan chan FileEvent) error {
	for {
		select {
		case event, ok := <-eventchan:
			if !ok {
				return errors.New("log file event channel closed")
			}
			if event.Modified() {
				return nil
			}
			if event.Closed() {
				return errors.New("log file closed")
			}
		case <-ctx.Done():
			return errors.New("log file` watcher cancelled")
		}
	}
}
