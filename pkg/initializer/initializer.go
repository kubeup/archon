package initializer

import (
	"fmt"
	"github.com/golang/glog"
	//"k8s.io/kubernetes/pkg/api"
	//"k8s.io/apimachinery/pkg/runtime"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"kubeup.com/archon/pkg/clientset"
)

var (
	ErrSkip = fmt.Errorf("This initializer will not process this resource")
)

type Object interface {
	metav1.Object

	GetInitializersInAnnotations() []string
	SetInitializersInAnnotations([]string)

	//GetFinalizers() []string
	//SetFinalizers([]string)
}

type Initializer interface {
	Token() string
	Initialize(Object) (Object, error, bool)
	Finalize(Object) (Object, error, bool)
	//ResolveConflict(old runtime.Object, new runtime.Object) (runtime.Object, error)
}

type InitializerMap map[string]Initializer

type InitializerManager struct {
	initializers InitializerMap
	kubeClient   clientset.Interface
}

func NewInitializerManager(list []Initializer, kubeClient clientset.Interface) (im *InitializerManager, err error) {
	imap := make(InitializerMap)
	for _, i := range list {
		imap[i.Token()] = i
	}

	return &InitializerManager{
		initializers: imap,
		kubeClient:   kubeClient,
	}, nil
}

func (im *InitializerManager) NeedsInitialization(obj Object) bool {
	if obj != nil && len(obj.GetInitializersInAnnotations()) > 0 {
		return true
	}

	return false
}

func (im *InitializerManager) Initialize(obj Object) (updatedObj Object, err error, retry bool) {
	tokens := obj.GetInitializersInAnnotations()
	glog.Infof("Initializer manager is initializing %v", tokens)
	for _, token := range tokens {
		init, ok := im.initializers[token]
		if !ok {
			glog.Infof("%s is not managed by initializer manager. will skip", token)
			continue
		}

		glog.Infof("Initializer manager is initializing %v", token)
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

func (im *InitializerManager) Finalize(obj Object) (updatedObj Object, err error, retry bool) {
	tokens := obj.GetFinalizers()
	glog.Infof("Initializer manager is finalizing %v", tokens)
	for _, token := range tokens {
		init, ok := im.initializers[token]
		if !ok {
			glog.Infof("%s is not managed by initializer manager. will skip", token)
			continue
		}

		glog.Infof("Initializer manager is finalizing %v", token)
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

func (im *InitializerManager) FinalizeAll(obj Object) (updatedObj Object, err error, retry bool) {
	retrys := 3
	var obj2 Object
	for retrys > 0 {
		tokens := obj.GetFinalizers()
		if len(tokens) == 0 {
			break
		}

		obj2, err, retry = im.Finalize(obj)
		if err != nil {
			retrys -= 1
			glog.Errorf("Unable to finalize object: %v", err)
		} else if obj2 == nil {
			glog.Warningf("Finalizer didn't return an update Object")
		} else {
			obj = obj2
		}
	}

	if err == nil {
		updatedObj = obj
	}
	return
}
