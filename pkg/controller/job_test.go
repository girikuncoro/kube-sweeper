package controller

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	batchv1 "k8s.io/api/batch/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func mockJob(completed time.Time, active, succeeded, failed int32) *batchv1.Job {
	t := metav1.NewTime(completed)
	return &batchv1.Job{
		Status: batchv1.JobStatus{
			CompletionTime: &t,
			Active:         active,
			Succeeded:      succeeded,
			Failed:         failed,
		},
	}
}

func TestDeleteJob(t *testing.T) {
	ts := time.Now()
	testcases := map[string]struct {
		job        *batchv1.Job
		successful time.Duration
		failed     time.Duration
		expected   bool
	}{
		"jobs with active pods should not be deleted": {
			job:        mockJob(ts.Add(-time.Minute), 1, 0, 0),
			successful: time.Second,
			failed:     time.Second,
			expected:   false,
		},
		"expired successful jobs should be deleted": {
			job:        mockJob(ts.Add(-time.Minute), 0, 1, 0),
			successful: time.Second,
			failed:     time.Second,
			expected:   true,
		},
		"non-expired successful jobs should not be deleted": {
			job:        mockJob(ts.Add(-time.Minute), 0, 1, 0),
			successful: 5 * time.Minute,
			failed:     time.Second,
			expected:   false,
		},
		"expired failed jobs should be deleted": {
			job:        mockJob(ts.Add(-time.Minute), 0, 0, 1),
			successful: time.Second,
			failed:     time.Second,
			expected:   true,
		},
		"non-expired failed jobs should not be deleted": {
			job:        mockJob(ts.Add(-time.Minute), 0, 0, 1),
			successful: time.Second,
			failed:     5 * time.Minute,
			expected:   false,
		},
	}
	for name, tc := range testcases {
		t.Run(name, func(t *testing.T) {
			res := isJobExpired(tc.job, tc.successful, tc.failed)
			assert.Equal(t, tc.expected, res, name)
		})
	}
}
