package controller

import (
	"context"
	"log"
	"reflect"
	"time"

	batchv1 "k8s.io/api/batch/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
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
		UpdateFunc: func(oldObj, newObj interface{}) {
			if !reflect.DeepEqual(oldObj, newObj) {
				sweeper.Process(newObj.(*batchv1.Job))
			}
		},
	})
	sweeper.jobInformer = jobInformer

	return sweeper
}

func (s *Sweeper) reconcile() {
	ticker := time.NewTicker(2 * resyncPeriod)
	for {
		select {
		case <-s.stopCh:
			ticker.Stop()
			return
		case <-ticker.C:
			for _, obj := range s.jobInformer.GetStore().List() {
				s.Process(obj.(*batchv1.Job))
			}
		}
	}
}

func (s *Sweeper) Run() {
	log.Printf("Waiting for change events...")

	go s.jobInformer.Run(s.stopCh)
	go s.reconcile()
	<-s.stopCh
}

func (s *Sweeper) Process(job *batchv1.Job) {
	if !isBeingDeleted(job) {
		return
	}
	if !isJobExpired(job, s.deleteSuccessfulAfter, s.deleteFailedAfter) {
		return
	}

	log.Printf("Deleting job %q in namespace %q", job.Name, job.Namespace)
	// Set cascading delete
	propagationPolicy := metav1.DeletePropagationForeground
	opts := metav1.DeleteOptions{
		PropagationPolicy: &propagationPolicy,
	}

	err := s.client.BatchV1().Jobs(job.Namespace).Delete(s.ctx, job.Name, opts)
	if err != nil && !apierrors.IsNotFound(err) {
		log.Printf("Failed to delete job %q in namespace %q: %v", job.Name, job.Namespace, err)
		return
	}
}
