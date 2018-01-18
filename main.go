package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/docktermj/go-hello-dependency-injection/common/runner"
	"github.com/docktermj/go-hello-dependency-injection/service/a"
	"github.com/docktermj/go-hello-dependency-injection/service/b"
	"github.com/docktermj/go-hello-dependency-injection/service/c"
	"github.com/docktermj/go-hello-dependency-injection/service/d"
	"github.com/docktermj/go-logger/logger"
	"github.com/docopt/docopt-go"
	"github.com/karlkfi/inject"
)

var (
	programName    = "go-hello-dependency-injection"
	buildVersion   = "2.0.0"
	buildIteration = "0"
	functions      = map[string]interface{}{}
	topContext     context.Context
	cancel         func()
)

// Dependency Injected variables.

var (
	ctx       context.Context
	myA       *a.A
	myB       *b.B
	myC       *c.C
	myD       *d.D
	waitGroup *sync.WaitGroup
)

// Utility function for dependency injection.
func newWaitGroup() *sync.WaitGroup {
	return &sync.WaitGroup{}
}

// Utility function for dependency injection.
func getTopContext() context.Context {
	return topContext
}

// ----------------------------------------------------------------------------
// Main
// ----------------------------------------------------------------------------

func main() {

	usage := `
Usage:
    go-hello-dependency-injection [<command>] [options]

Options:
   -h, --help                         Show this help
`

	// DocOpt processing.

	commandVersion := fmt.Sprintf("%s %s-%s", programName, buildVersion, buildIteration)
	args, _ := docopt.Parse(usage, nil, true, commandVersion, false)

	// Configure output log.

	log.SetFlags(log.Lshortfile | log.Ldate | log.Lmicroseconds | log.LUTC)
	logger.SetLevel(logger.LevelDebug)

	// Create top-level context.

	topContext, cancel = context.WithCancel(context.Background())
	defer cancel()

	// Handle Ctrl-c interrupts to cancel the context.

	interruptChannel := make(chan os.Signal, 2)
	signal.Notify(interruptChannel, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-interruptChannel
		cancel()
	}()

	// Show debugging information.

	logger.Infof("Starting %s version %s", programName, buildVersion)
	if logger.IsDebug() {
		logger.Debugf("os.Args: %+v\n", os.Args)
		logger.Debugf("args: %+v\n", args)
		logger.Debugf("topContext: %+v\n", topContext)
	}

	// If subcommand was specified, handle it and exit.

	if args["<command>"] != nil {
		_, hasSubcommand := functions[args["<command>"].(string)]
		if hasSubcommand {
			argv := os.Args[1:]
			runner.Run(topContext, argv, functions, usage)
			cancel()
			os.Exit(0)
		}
	}

	// Object creation and Dependency Injection (DI).

	di_container := inject.NewGraph()
	di_container.Define(&waitGroup, inject.NewProvider(newWaitGroup))
	di_container.Define(&ctx, inject.NewProvider(getTopContext))
	di_container.Define(&myD, inject.NewAutoProvider(d.New))
	di_container.Define(&myC, inject.NewAutoProvider(c.New))
	di_container.Define(&myB, inject.NewAutoProvider(b.New))
	di_container.Define(&myA, inject.NewAutoProvider(a.New))
	di_container.ResolveAll()

	// List services to be started.

	startServices := []interface{}{
		myA.Start,
		myB.Start,
		myC.Start,
		myD.Start,
	}

	// Start services.

	for _, service := range startServices {
		waitGroup.Add(1)
		fn := service.(func() error)
		go func() {
			err := fn()
			if err != nil {
				logger.Errorf("error starting service - %s", err.Error())
			}
		}()
	}

	// Wait until all processes are done or termination.

	waitGroup.Wait()

	// List services to be stopped.

	stopServices := []interface{}{
		myD.Stop,
		myC.Stop,
		myB.Stop,
		myA.Stop,
	}

	// Start services.
	// FIXME: For some reason, the object doesn't retain the "iterations".

	for _, service := range stopServices {
		waitGroup.Add(1)
		fn := service.(func() error)
		go func() {
			err := fn()
			if err != nil {
				logger.Infof("error stopping service - %s", err.Error())
			}
		}()
	}

	waitGroup.Wait()

	// Epilog.

	logger.Info("All done.")
}
