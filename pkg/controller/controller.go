package controller

import (
	"context"
	"log"
	"reflect"
	"time"

	batchv1 "k8s.io/api/batch/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
)

const (
	resyncPeriod = 30 * time.Second
)

type Sweeper struct {
	ctx    context.Context
	client *kubernetes.Clientset
	stopCh <-chan struct{}

	jobInformer cache.SharedIndexInformer

	deleteSuccessfulAfter time.Duration
	deleteFailedAfter     time.Duration
}

func NewSweeper(ctx context.Context, client *kubernetes.Clientset, namespace string,
	deleteSuccessfulAfter, deleteFailedAfter time.Duration, stopCh <-chan struct{}) *Sweeper {
	sweeper := &Sweeper{
		ctx:                   ctx,
		client:                client,
		stopCh:                stopCh,
		deleteSuccessfulAfter: deleteSuccessfulAfter,
		deleteFailedAfter:     deleteFailedAfter,
	}

	jobInformer := cache.NewSharedIndexInformer(
		&cache.ListWatch{
			ListFunc: func(opts metav1.ListOptions) (runtime.Object, error) {
				return client.BatchV1().Jobs(namespace).List(ctx, opts)
			},
			WatchFunc: func(opts metav1.ListOptions) (watch.Interface, error) {
				return client.BatchV1().Jobs(namespace).Watch(ctx, opts)
			},
		},
		&batchv1.Job{},
		resyncPeriod,
		cache.Indexers{},
	)

	jobInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		UpdateFunc: func(oldJob, newJob interface{}) {
			if !reflect.DeepEqual(oldJob, newJob) {
				log.Printf("Cleaning up jobs: %v\n", newJob)
			}
		},
	})
	sweeper.jobInformer = jobInformer

	return sweeper
}

// Run listens for job/pod changes and handle accordingly.
func (s *Sweeper) Run() {
	log.Printf("Waiting for change events...")

	go s.jobInformer.Run(s.stopCh)
	<-s.stopCh
}
