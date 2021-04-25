package worker

import (
	"context"
	"fmt"
	"os"

	"github.com/hpcloud/tail"
)

type Logger struct {
	config Config
}

func NewLogger(config Config) Logger {
	return Logger{
		config: config,
	}
}

func (l *Logger) Path(name string) (logpath string) {
	return fmt.Sprintf("%s%s.log", l.config.LogFolder(), name)
}

func (l *Logger) Create(name string) (f *os.File, err error) {
	return os.Create(l.Path(name))
}

func (l *Logger) Remove(name string) (err error) {
	return os.Remove(l.Path(name))
}

func (l *Logger) Tailf(ctx context.Context, name string) (logchan chan string, err error) {
	// tailling the log output file
	t, err := tail.TailFile(l.Path(name), tail.Config{Follow: true})
	if err != nil {
		return
	}
	logchan = make(chan string)
	// sends the log tail to the log channel
	go func() {
		defer close(logchan)
		for line := range t.Lines {
			logchan <- line.Text
		}
	}()
	// waits context cancelation to stop the tail log
	go func() {
		<-ctx.Done()
		err = t.Stop()
	}()
	return
}
