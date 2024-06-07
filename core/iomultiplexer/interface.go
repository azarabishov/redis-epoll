package iomultiplexer

import "time"

type IOMultiplexer interface {
	Subscribe(event Event) error
	Poll(timeout time.Duration) ([]Event, error)
	Close() error
}
