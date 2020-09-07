# Kube Sweeper

Kubernetes controller that listens to completed Jobs and Pods and automatically delete them after X seconds (default to 1 hour). This project is inspired from [kube-job-cleaner](https://github.com/hjacobs/kube-job-cleaner) and [kube-cleanup-operator](https://github.com/lwolf/kube-cleanup-operator).

Kubernetes Jobs are not cleaned up by default and completed Pods are never deleted. Jobs that are run frequently causing unnecessary Pod resources which significantly slowdown the Kubernetes API server. This controller listens and cleans up the completed Jobs/Pods, as well as perform periodic cleanup for existing resources.

## Development

Building binary:
```sh
$ make
```

Running the binary:
```sh
./bin/kubesweeper --delete-successful-after-seconds 300
```

## Usage

Deploying:
```sh
$ kubectl apply -f deploy/
```

There are few options:
| flag                               | description                                                           |
| ---------------------------------- | --------------------------------------------------------------------- |
| delete-successful-after-seconds    | number of seconds after successful job completion to remove the job   |
| delete-failed-after-seconds        | number of seconds after failed job completion to remove the job       |
| namespace                          | only cleanup jobs in single namespace                                 |
