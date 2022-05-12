package match

import (
	"duel-masters/internal"

	"github.com/sirupsen/logrus"
)

type EventExecutionStrategy int

const (
	SequentialEvent EventExecutionStrategy = iota
	ParallelEvent
)

type EventLoop struct {
	stopped bool
	events  chan func()
	exit    chan bool
}

func NewEventLoop() *EventLoop {
	return &EventLoop{
		events: make(chan func(), 1),
		exit:   make(chan bool),
	}
}

func (el *EventLoop) start() {
	defer internal.Recover()
	defer logrus.Debug("Stopped event loop")

	for {
		select {
		case <-el.exit:
			close(el.events)
			return
		case event := <-el.events:
			el.process(event)

			// drain the channel
			for len(el.events) > 0 {
				<-el.events
			}
		}
	}
}

func (el *EventLoop) stop() {

	if el.stopped {
		return
	}

	// We run this in a separate goroutine because an event may currently be
	// processing and blocked by an expected user action. This will wait for
	// the player action channel to be closed without blocking the caller of
	// this function. When the player action channel is closed, the eventloop
	// will be free to receive the stop signal.
	go func() {
		defer internal.Recover()
		el.exit <- true
	}()

}

func (el *EventLoop) schedule(event func(), strategy EventExecutionStrategy) {

	go func() {
		defer internal.Recover()

		switch strategy {
		case ParallelEvent:
			event()
			return

		case SequentialEvent:
			select {
			case el.events <- event:
			default:
				logrus.Debug("Skipped an incoming event")
			}
		}
	}()

}
func (el *EventLoop) process(event func()) {

	defer internal.Recover()

	event()

}
