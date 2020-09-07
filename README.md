[![CI](https://github.com/girikuncoro/kube-sweeper/workflows/master/badge.svg)][ci]
[![Go Report Card](https://goreportcard.com/badge/github.com/girikuncoro/kube-sweeper)][goreportcard]
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)][license]

[ci]: https://github.com/girikuncoro/kube-sweeper/actions?query=master+branch%3Amaster
[goreportcard]: https://goreportcard.com/report/github.com/girikuncoro/kube-sweeper
[license]: https://opensource.org/licenses/Apache-2.0

# Kube Sweeper

Kubernetes controller that listens to completed Jobs and Pods and automatically delete them after X seconds (default to 15 minutes). This project is inspired from [kube-job-cleaner](https://github.com/hjacobs/kube-job-cleaner) and [kube-cleanup-operator](https://github.com/lwolf/kube-cleanup-operator).

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
