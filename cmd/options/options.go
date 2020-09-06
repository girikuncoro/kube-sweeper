package options

import (
	"flag"

	"github.com/spf13/pflag"
)

// KubeSweeperOptions contains kube sweeper command line options.
type KubeSweeperOptions struct {
	// Namespace determines the namespace scope to perform the resource cleanup for.
	Namespace string
	// DeleteSuccessAfterSeconds determines the number of seconds needed before deleting successful Jobs.
	DeleteSuccessAfterSeconds int
	// DeleteFailedAfterSeconds determines the number of seconds needed before deleting failed Jobs.
	DeleteFailedAfterSeconds int
}

func NewKubeSweeperOptions() *KubeSweeperOptions {
	return &KubeSweeperOptions{}
}

// AddFlags adds kube sweeper command line options to pflag.
func (opts *KubeSweeperOptions) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&opts.Namespace, "namespace", "", "Limit scope to single namespace")
	fs.IntVar(&opts.DeleteSuccessAfterSeconds, "delete-success-after-seconds", 900, "Delete successful jobs after X seconds")
	fs.IntVar(&opts.DeleteFailedAfterSeconds, "delete-failed-after-seconds", 0, "Delete failed jobs after X seconds")
}

func init() {
	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
}
