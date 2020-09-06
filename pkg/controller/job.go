package controller

import (
	"time"

	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
)

func isJobExpired(job *batchv1.Job, deleteSuccessfulAfter, deleteFailedAfter time.Duration) bool {
	completionTime := getJobCompletionTime(job)
	if completionTime.IsZero() {
		return false
	}

	timeSinceCompletion := time.Since(completionTime)
	if isSuccessfulJob(job) {
		if deleteSuccessfulAfter > 0 && timeSinceCompletion > deleteSuccessfulAfter {
			return true
		}
	}

	if isFailedJob(job) {
		if deleteFailedAfter > 0 && timeSinceCompletion >= deleteFailedAfter {
			return true
		}
	}

	return true
}

func getJobCompletionTime(job *batchv1.Job) time.Time {
	if !job.Status.CompletionTime.IsZero() {
		return job.Status.CompletionTime.Time
	}
	for _, j := range job.Status.Conditions {
		if j.Type == batchv1.JobFailed && j.Status == corev1.ConditionTrue {
			return j.LastTransitionTime.Time
		}
	}
	return time.Time{}
}

func isSuccessfulJob(job *batchv1.Job) bool {
	return job.Status.Succeeded > 0
}

func isFailedJob(job *batchv1.Job) bool {
	if job.Status.Failed > 0 {
		return true
	}

	for _, j := range job.Status.Conditions {
		if j.Type == batchv1.JobFailed && j.Status == corev1.ConditionTrue {
			return true
		}
	}
	return false
}

func isBeingDeleted(job *batchv1.Job) bool {
	return job.DeletionTimestamp.IsZero()
}
