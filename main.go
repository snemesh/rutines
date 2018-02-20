package main

import (
	"time"
	"github.com/inn4sc/vcg-go-common/routines"
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"sync"
	"os"
	"os/signal"
	"syscall"
	log "github.com/sirupsen/logrus"
)

var Revizor routines.Chief

func init() {
	Revizor = routines.Chief{}
	Revizor.AddWorkman("MyWorker", &WorckMan{})
	Revizor.InitWorkers(nil)
}

func main() {

	println("hello world")
	wg:= sync.WaitGroup{}
	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		defer wg.Done()
		//
		wg.Add(1)
		Revizor.Start(ctx)
		fmt.Print("Start 1st group")
	}()


	go func() {
		defer wg.Done()
		//
		wg.Add(1)
		Revizor.Start(ctx)
		fmt.Print("Start 2st group")
	}()

	var gracefulStop = make(chan os.Signal)
	signal.Notify(gracefulStop, syscall.SIGTERM, syscall.SIGINT)

	<-gracefulStop

	cancel()
	done := make(chan struct{})

	// Close done after wg finishes.
	go func() {
		defer close(done)
		wg.Wait()
	}()

	select {
	case <-done:
		//log.Default.Info("Graceful exit.")
		//return nil
	case <-time.NewTimer(60 * time.Second).C:
		//log.Default.Warn("Graceful exit timeout exceeded. Force exit.")
		//return cli.NewExitError("Graceful exit timeout exceeded", 1)
	}

}

type WorckMan struct {
	id int
	name string
	timeout time.Duration
	logger *logrus.Entry
	ctx context.Context
}

func (w *WorckMan)Run()  {
	defer w.Recover()

	ticker := time.NewTicker(w.timeout * time.Second)
	for {
		select {
		case <-ticker.C:
			w.applyOperations()
		case <-w.ctx.Done():
			ticker.Stop()
			w.logger.Info("End job")
			return
		}
	}
}

func (w *WorckMan)New(parentCtx context.Context)routines.Workman  {

	return &WorckMan{
		name: "Grabber",
		logger: log.WithField("worker", "RateGrabber"),
		ctx: parentCtx,
		timeout: time.Duration(15),
	}
}

func (s *WorckMan) Recover() {
	err := recover()
	if err == nil {
		return
	}
	s.logger.WithField("panic", err).Error("Caught panic")

}

func (w *WorckMan) applyOperations() {

	fmt.Println("Runn job")

}
