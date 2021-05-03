package log

import (
	"context"
	"log"
	"syscall"
	"unsafe"
)

type fileLinuxEvent struct {
	mask uint32
}

// Modified returns true if the event is syscall.IN_MODIFY.
func (e *fileLinuxEvent) Modified() bool {
	return e.mask&syscall.IN_MODIFY == syscall.IN_MODIFY
}

// Closed returns true if the event is syscall.IN_CLOSE.
func (e *fileLinuxEvent) Closed() bool {
	return e.mask&syscall.IN_CLOSE == syscall.IN_CLOSE
}

// wrapLinuxEvent create a new File event given an event mask.
func wrapLinuxEvent(mask uint32) FileEvent {
	return &fileLinuxEvent{
		mask: mask,
	}
}

// linuxWatcher Watcher implementation for Linux
type linuxWatcher struct {
	mask uint32
}

// NewWatcher instance
func NewWatcher() Watcher {
	return &linuxWatcher{
		mask: syscall.IN_MODIFY | syscall.IN_CLOSE,
	}
}

// Watcher watch file system events and send them throught a channel.
// See https://linux.die.net/man/1/inotifywait
func (w *linuxWatcher) Watch(ctx context.Context, path string) (chan FileEvent, error) {
	filed, err := syscall.InotifyInit()
	if err != nil {
		return nil, err
	}
	// watching for file system events given a file path
	watched, err := syscall.InotifyAddWatch(filed, path, w.mask)
	if err != nil {
		fderr := syscall.Close(filed)
		if fderr != nil {
			log.Printf("Fail to close file descriptor: %v", fderr)
		}
		return nil, err
	}
	eventchan := make(chan FileEvent)
	go func() {
		// remove watched from the watch list
		defer func() {
			success, err := syscall.InotifyRmWatch(filed, uint32(watched))
			if success == -1 || err != nil {
				log.Printf("Fail to remove the file watch: %v", err)
			}
			close(eventchan)
		}()
		// buffer to store events sent by OS
		buf := make([]byte, syscall.SizeofInotifyEvent*4096)
		for {
			// read events from file descriptor and fill the event buffer
			nbytes, err := syscall.Read(filed, buf[:])
			if err != nil {
				nbytes = 0
				log.Printf("Fail to read events: %v", err)
				return
			}
			// iterate over all buffered events and send them through the channel
			offset := 0
			for offset <= nbytes-syscall.SizeofInotifyEvent {
				// read raw event from the buffer given a position
				raw := (*syscall.InotifyEvent)(unsafe.Pointer(&buf[offset]))
				// calculate the next offset
				offset += syscall.SizeofInotifyEvent + int(raw.Len)
				// filtering events
				mask := raw.Mask
				if !w.acceptEvent(mask) {
					continue
				}
				select {
				// send the event throught the channel
				case eventchan <- wrapLinuxEvent(mask):
				// stop the watcher
				case <-ctx.Done():
					return
				}
			}
		}
	}()
	return eventchan, nil
}

// acceptEvent checks if the event can be accepted.
func (w *linuxWatcher) acceptEvent(mask uint32) bool {
	return w.mask&mask != 0
}
