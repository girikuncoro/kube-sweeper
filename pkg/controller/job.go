package controller

import (
	"time"

	batchv1 "k8s.io/api/batch/v1"
)

// TODO: to be implemented
func isJobExpired(job *batchv1.Job, deleteSuccessfulAfter, deleteFailedAfter time.Duration) bool {
	return true
}
