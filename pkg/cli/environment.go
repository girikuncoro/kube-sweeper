package cli

import (
	"os"

	"k8s.io/cli-runtime/pkg/genericclioptions"
)

// EnvSettings describes all of the environment settings.
type EnvSettings struct {
	namespace string
	config    *genericclioptions.ConfigFlags

	KubeConfig    string
	KubeContext   string
	KubeAPIServer string
}

func New() *EnvSettings {
	env := &EnvSettings{
		namespace:     os.Getenv("KUBESWEEPER_NAMESPACE"),
		KubeContext:   os.Getenv("KUBESWEEPER_KUBECONTEXT"),
		KubeAPIServer: os.Getenv("KUBESWEEPER_KUBEAPISERVER"),
	}

	env.config = &genericclioptions.ConfigFlags{
		Namespace:  &env.namespace,
		Context:    &env.KubeContext,
		APIServer:  &env.KubeAPIServer,
		KubeConfig: &env.KubeConfig,
	}

	return env
}

func (s *EnvSettings) RESTClientGetter() genericclioptions.RESTClientGetter {
	return s.config
}
