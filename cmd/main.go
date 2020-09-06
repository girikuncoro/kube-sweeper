package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/girikuncoro/kube-sweeper/cmd/options"
	"github.com/girikuncoro/kube-sweeper/pkg/cli"
	"github.com/girikuncoro/kube-sweeper/pkg/controller"
	"github.com/girikuncoro/kube-sweeper/pkg/kube"
	"github.com/spf13/pflag"
	"k8s.io/klog"
)

var settings = cli.New()

func initKubeLogs() {
	log.SetOutput(os.Stdout)
	gofs := flag.NewFlagSet("klog", flag.ExitOnError)
	klog.InitFlags(gofs)
	pflag.CommandLine.Set("logtostderr", "true")
}

func main() {
	opts := options.NewKubeSweeperOptions()
	opts.AddFlags(pflag.CommandLine)

	pflag.Parse()
	initKubeLogs()

	log.Printf("Starting kube-sweeper operator...")

	// Channel to receive OS signals
	sigsCh := make(chan os.Signal, 1)
	// Channel to receive stop signal
	stopCh := make(chan struct{})
	// Register sigsCH to receive SIGTERM signal
	signal.Notify(sigsCh, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	wg := &sync.WaitGroup{}

	clientset, err := kube.NewClientSet(settings.RESTClientGetter())
	if err != nil {
		log.Fatal(err.Error())
	}

	wg.Add(1)
	go func() {
		controller.NewSweeper(
			context.Background(),
			clientset,
			opts.Namespace,
			time.Duration(opts.DeleteSuccessAfterSeconds)*time.Second,
			time.Duration(opts.DeleteFailedAfterSeconds)*time.Second,
			stopCh,
		).Run()
		wg.Done()
	}()
	log.Printf("kube-sweeper controller started...")

	// TODO: expose metrics
	<-sigsCh
	log.Printf("received termination signal...")
	close(stopCh)
	wg.Wait()
}
