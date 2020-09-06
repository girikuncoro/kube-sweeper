package kube

import (
	"github.com/pkg/errors"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

type RESTClientGetter interface {
	ToRESTConfig() (*rest.Config, error)
}

func NewClientSet(getter *rest.Config) (*kubernetes.Clientset, error) {
	conf, err := getter.ToRESTConfig()
	if err != nil {
		return nil, errors.Wrap(err, "unable to generate config for kubernetes client")
	}
	return kubernetes.NewForConfig(conf)
}
