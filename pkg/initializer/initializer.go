package initializer

import (
	"fmt"
	"github.com/golang/glog"
	"k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/runtime"
	"kubeup.com/archon/pkg/clientset"
	"kubeup.com/archon/pkg/cloudprovider"
)

var (
	ErrSkip = fmt.Errorf("This initializer will not process this resource")
)

type Initializer interface {
	Token() string
	Initialize(runtime.Object) (runtime.Object, error, bool)
	Finalize(runtime.Object) (runtime.Object, error, bool)
	//ResolveConflict(old runtime.Object, new runtime.Object) (runtime.Object, error)
}

type InitializerMap map[string]Initializer
type InitializerFactory func(clientset.Interface, cloudprovider.Interface, string) (Initializer, error)
type Factories []InitializerFactory

type InitializerManager struct {
	initializers InitializerMap
	kubeClient   clientset.Interface
}

func NewInitializerManager(factories Factories, kubeClient clientset.Interface, cloud cloudprovider.Interface, clusterName string) (im *InitializerManager, err error) {
	var i Initializer
	imap := make(InitializerMap)
	for _, f := range factories {
		i, err = f(kubeClient, cloud, clusterName)
		if err != nil {
			return
		}
		imap[i.Token()] = i
	}

	return &InitializerManager{
		initializers: imap,
		kubeClient:   kubeClient,
	}, nil
}

func (im *InitializerManager) Initialize(obj runtime.Object) (updatedObj runtime.Object, err error, retry bool) {
	m, err := api.ObjectMetaFor(obj)
	if err != nil {
		return
	}

	tokens := m.Initializers
	for _, token := range tokens {
		init, ok := im.initializers[token]
		if !ok {
			glog.Infof("%s is not managed by initializer manager. will skip", token)
			continue
		}

		updatedObj, err, retry = init.Initialize(obj)
		if err == ErrSkip {
			continue
		} else if err != nil {
			glog.Errorf("Initializer error: %v", err)
			return
		}

		break
	}

	return
}

func (im *InitializerManager) Finalize(obj runtime.Object) (updatedObj runtime.Object, err error, retry bool) {
	m, err := api.ObjectMetaFor(obj)
	if err != nil {
		return
	}

	tokens := m.Finalizers
	for _, token := range tokens {
		init, ok := im.initializers[token]
		if !ok {
			glog.Infof("%s is not managed by initializer manager. will skip", token)
			continue
		}

		updatedObj, err, retry = init.Finalize(obj)
		if err == ErrSkip {
			continue
		} else if err != nil {
			glog.Errorf("Finalizer error: %v", err)
			return nil, err, true
		}

		break
	}

	return
}
