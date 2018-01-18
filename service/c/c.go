package c

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/docktermj/go-hello-dependency-injection/service/a"
	"github.com/docktermj/go-hello-dependency-injection/service/b"
	"github.com/docktermj/go-logger/logger"
)

type C struct {
	Context   context.Context
	A         a.A
	B         b.B
	WaitGroup *sync.WaitGroup
	iteration int
}

func New(ctx context.Context, myA *a.A, myB *b.B, waitGroup sync.WaitGroup) *C {
	return &C{
		Context:   ctx,
		A:        *myA,
		B:        *myB,
		WaitGroup: &waitGroup,
	}
}

// ----------------------------------------------------------------------------
// Utility methods
// ----------------------------------------------------------------------------

func (c C) Speak() string {
	return fmt.Sprintf("C says that: %s", c.B.Speak())
}

// ----------------------------------------------------------------------------
// Service interface
// ----------------------------------------------------------------------------

func (c C) Start() error {

	// Synchronize the services at shutdown.

	if c.WaitGroup != nil {
		defer c.WaitGroup.Done()
	}

	// Loop.

PrintLoop:
	for {
		time.Sleep(5 * time.Second)

		select {
		case <-c.Context.Done():
			break PrintLoop
		default:
		}

		c.iteration = c.iteration + 1
		logger.Infof("%s (iteration  %d)", c.Speak(), c.iteration)
	}

	// Epilog.

	logger.Infof("C is done.")
	return nil
}

func (c C) Stop() error {

	// Synchronize the services at shutdown.

	if c.WaitGroup != nil {
		defer c.WaitGroup.Done()
	}

	logger.Infof("C iterations: %d", c.iteration)
	return nil
}
