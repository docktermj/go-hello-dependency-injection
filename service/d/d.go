package d

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/docktermj/go-hello-dependency-injection/service/a"
	"github.com/docktermj/go-hello-dependency-injection/service/b"
	"github.com/docktermj/go-hello-dependency-injection/service/c"
	"github.com/docktermj/go-logger/logger"
)

type D struct {
	Context   context.Context
	A         *a.A
	B         *b.B
	C         *c.C
	WaitGroup *sync.WaitGroup
	iteration int
}

func (d D) Speak() string {
	return fmt.Sprintf("D says that: %s", d.C.Speak())
}

func (d D) Speak2() string {
	return fmt.Sprintf("D says that: %s", d.A.Speak())
}

// ----------------------------------------------------------------------------
// Service interface
// ----------------------------------------------------------------------------

func (d D) Start() error {

	// Synchronize the services at shutdown.

	if d.WaitGroup != nil {
		defer d.WaitGroup.Done()
	}

	// Loop.

PrintLoop:
	for {
		time.Sleep(8 * time.Second)

		select {
		case <-d.Context.Done():
			break PrintLoop
		default:
		}

		d.iteration = d.iteration + 1
		logger.Infof("%s (iteration  %d)", d.Speak(), d.iteration)
		logger.Infof("%s", d.Speak2())
	}

	// Epilog.

	logger.Infof("D is done.")
	return nil
}

func (d D) Stop() error {

	// Synchronize the services at shutdown.

	if d.WaitGroup != nil {
		defer d.WaitGroup.Done()
	}

	logger.Infof("D iterations: %d", d.iteration)
	return nil
}
