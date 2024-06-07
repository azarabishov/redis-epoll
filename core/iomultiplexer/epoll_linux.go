package iomultiplexer

import (
	"fmt"
	"syscall"
	"time"
)

type Epoll struct {
	fd int
	ePollEvents []syscall.EpollEvent
	diceEvents []Event
}

func New(maxClients int) (*Epoll, error) {
	if maxClients < 0 {
		return nil, ErrInvalidMaxClients
	}

	fd, err := syscall.EpollCreate1(0)
	if err != nil {
		return nil, err
	}

	return &Epoll{
		fd:          fd,
		ePollEvents: make([]syscall.EpollEvent, maxClients),
		diceEvents:  make([]Event, maxClients),
	}, nil
}

func (ep *Epoll) Subscribe(event Event) error {
	nativeEvent := event.toNative()
	if err := syscall.EpollCtl(ep.fd, syscall.EPOLL_CTL_ADD, event.Fd, &nativeEvent); err != nil {
		return fmt.Errorf("epoll subscribe: %w", err)
	}
	return nil
}

func (ep *Epoll) Poll(timeout time.Duration) ([]Event, error) {
	nEvents, err := syscall.EpollWait(ep.fd, ep.ePollEvents, newTime(timeout))
	if err != nil {
		return nil, fmt.Errorf("epoll poll: %w", err)
	}

	for i := 0; i < nEvents; i++ {
		ep.diceEvents[i] = newEvent(ep.ePollEvents[i])
	}

	return ep.diceEvents[:nEvents], nil
}

func (ep *Epoll) Close() error {
	return syscall.Close(ep.fd)
}
