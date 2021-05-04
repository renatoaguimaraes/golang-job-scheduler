package log

import "context"

// FileEvent wrap an inotify event.
type FileEvent interface {
	// Modified returns true if the file was changed.
	Modified() bool
	// Deleted returns true if the file was closed.
	Closed() bool
}

// Watcher
type Watcher interface {
	// Watch monitoring file system changes and send the events throught a channel.
	Watch(ctx context.Context, path string) (eventchan chan FileEvent, err error)
}
