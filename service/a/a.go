package a

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/docktermj/go-logger/logger"
)

type A struct {
	Context   context.Context
	Greetings string
	WaitGroup *sync.WaitGroup
	iteration int
}

func New(ctx context.Context, waitGroup *sync.WaitGroup) *A {
	return &A{
		Context:   ctx,
		Greetings: "Hello world!",
		WaitGroup: waitGroup,
	}
}

// ----------------------------------------------------------------------------
// Utility methods
// ----------------------------------------------------------------------------

func (a A) Speak() string {
	return fmt.Sprintf("A says: %s", a.Greetings)
}

// ----------------------------------------------------------------------------
// Service interface
// ----------------------------------------------------------------------------

func (a A) Start() error {

	// Synchronize the services at shutdown.

	if a.WaitGroup != nil {
		defer a.WaitGroup.Done()
	}

	// Loop.

PrintLoop:
	for {
		time.Sleep(2 * time.Second)

		select {
		case <-a.Context.Done():
			break PrintLoop
		default:
		}

		a.iteration = a.iteration + 1
		logger.Infof("%s (iteration  %d)", a.Speak(), a.iteration)
	}

	// Epilog.

	logger.Infof("A is done.")
	return nil
}

func (a A) Stop() error {

	// Synchronize the services at shutdown.

	if a.WaitGroup != nil {
		defer a.WaitGroup.Done()
	}

	logger.Infof("A iterations: %d", a.iteration)
	return nil
}
