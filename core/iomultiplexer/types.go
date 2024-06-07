package iomultiplexer

type Event struct {
	Fd int
	Op Operations
}

type Operations uint32
