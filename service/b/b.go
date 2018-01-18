package b

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/docktermj/go-hello-dependency-injection/service/a"
	"github.com/docktermj/go-logger/logger"
)

type B struct {
	Context   context.Context
	A         *a.A
	WaitGroup *sync.WaitGroup
	iteration int
}

func New(ctx context.Context, myA *a.A, waitGroup *sync.WaitGroup) *B {
	return &B{
		Context:   ctx,
		A:         myA,
		WaitGroup: waitGroup,
	}
}

// ----------------------------------------------------------------------------
// Utility methods
// ----------------------------------------------------------------------------

func (b B) Speak() string {
	return fmt.Sprintf("B says that: %s", b.A.Speak())
}

// ----------------------------------------------------------------------------
// Service interface
// ----------------------------------------------------------------------------

func (b B) Start() error {

	// Synchronize the services at shutdown.

	if b.WaitGroup != nil {
		defer b.WaitGroup.Done()
	}

	// Loop.

PrintLoop:
	for {
		time.Sleep(3 * time.Second)

		select {
		case <-b.Context.Done():
			break PrintLoop
		default:
		}

		b.iteration = b.iteration + 1
		logger.Infof("%s (iteration  %d)", b.Speak(), b.iteration)
	}

	// Epilog.

	logger.Infof("B is done.")
	return nil
}

func (b B) Stop() error {

	// Synchronize the services at shutdown.

	if b.WaitGroup != nil {
		defer b.WaitGroup.Done()
	}

	logger.Infof("B iterations: %d", b.iteration)
	return nil
}
