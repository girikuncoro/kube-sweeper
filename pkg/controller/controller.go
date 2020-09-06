package controller

import (
	"context"
	"log"
	"time"

	"k8s.io/client-go/kubernetes"
)

type Sweeper struct{}

func NewSweeper(ctx context.Context, clientset *kubernetes.Clientset, namespace string,
	deleteSuccessfulAfterSeconds, deleteFailedAfterSeconds time.Duration, stopCh <-chan struct{}) *Sweeper {
	return &Sweeper{}
}

// Run listens for job/pod changes and handle accordingly.
func (s *Sweeper) Run() {
	log.Printf("Waiting for change events...")
}
