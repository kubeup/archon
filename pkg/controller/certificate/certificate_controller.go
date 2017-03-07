/*
Copyright 2016 The Archon Authors.
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package certificate

import (
	"fmt"
	"time"

	archoncache "kubeup.com/archon/pkg/cache"
	"kubeup.com/archon/pkg/clientset"
	"kubeup.com/archon/pkg/controller/instance"

	"k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/client/cache"
	unversionedcore "k8s.io/kubernetes/pkg/client/clientset_generated/internalclientset/typed/core/internalversion"
	"k8s.io/kubernetes/pkg/client/record"
	"k8s.io/kubernetes/pkg/controller"
	"k8s.io/kubernetes/pkg/runtime"
	utilruntime "k8s.io/kubernetes/pkg/util/runtime"
	"k8s.io/kubernetes/pkg/util/wait"
	"k8s.io/kubernetes/pkg/util/workqueue"
	"k8s.io/kubernetes/pkg/watch"

	"github.com/golang/glog"
)

type CertificateController struct {
	kubeClient         clientset.Interface
	namespace          string
	certificateControl instance.CertificateControlInterface

	// CSR framework and store
	csrController *cache.Controller
	csrStore      archoncache.StoreToSecretLister

	syncHandler func(csrKey string) error

	queue workqueue.RateLimitingInterface
}

func New(kubeClient clientset.Interface, syncPeriod time.Duration, caCertFile, caKeyFile string, namespace string) (*CertificateController, error) {
	// Send events to the apiserver
	eventBroadcaster := record.NewBroadcaster()
	eventBroadcaster.StartLogging(glog.Infof)
	eventBroadcaster.StartRecordingToSink(&unversionedcore.EventSinkImpl{Interface: kubeClient.Core().Events("")})

	certControl, err := instance.NewCertificateControl(caCertFile, caKeyFile)
	if err != nil {
		glog.Errorf("WARNING: Unable to start certificate controller: %s", err.Error())
		return nil, err
	}

	cc := &CertificateController{
		kubeClient:         kubeClient,
		namespace:          namespace,
		certificateControl: certControl,
		queue:              workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "certificate"),
	}

	// Manage the addition/update of certificate requests
	cc.csrStore.Indexer, cc.csrController = cache.NewIndexerInformer(
		&cache.ListWatch{
			ListFunc: func(options api.ListOptions) (runtime.Object, error) {
				return cc.kubeClient.Core().Secrets(namespace).List(options)
			},
			WatchFunc: func(options api.ListOptions) (watch.Interface, error) {
				return cc.kubeClient.Core().Secrets(namespace).Watch(options)
			},
		},
		&api.Secret{},
		syncPeriod,
		cache.ResourceEventHandlerFuncs{
			AddFunc: func(obj interface{}) {
				csr := obj.(*api.Secret)
				glog.V(4).Infof("Adding certificate request %s", csr.Name)
				cc.enqueueCertificateRequest(obj)
			},
			UpdateFunc: func(old, new interface{}) {
				cc.enqueueCertificateRequest(new)
			},
			DeleteFunc: func(obj interface{}) {
				csr, ok := obj.(*api.Secret)
				if !ok {
					tombstone, ok := obj.(cache.DeletedFinalStateUnknown)
					if !ok {
						glog.V(2).Infof("Couldn't get object from tombstone %#v", obj)
						return
					}
					csr, ok = tombstone.Obj.(*api.Secret)
					if !ok {
						glog.V(2).Infof("Tombstone contained object that is not a CSR: %#v", obj)
						return
					}
				}
				glog.V(4).Infof("Deleting certificate request %s", csr.Name)
				cc.enqueueCertificateRequest(obj)
			},
		},
		cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc},
	)
	cc.syncHandler = cc.maybeSignCertificate
	return cc, nil
}

// Run the main goroutine responsible for watching and syncing jobs.
func (cc *CertificateController) Run(workers int, stopCh <-chan struct{}) {
	defer utilruntime.HandleCrash()
	defer cc.queue.ShutDown()

	go cc.csrController.Run(stopCh)

	glog.Infof("Starting certificate controller manager")
	for i := 0; i < workers; i++ {
		go wait.Until(cc.worker, time.Second, stopCh)
	}
	<-stopCh
	glog.Infof("Shutting down certificate controller")
}

// worker runs a thread that dequeues CSRs, handles them, and marks them done.
func (cc *CertificateController) worker() {
	for cc.processNextWorkItem() {
	}
}

// processNextWorkItem deals with one key off the queue.  It returns false when it's time to quit.
func (cc *CertificateController) processNextWorkItem() bool {
	cKey, quit := cc.queue.Get()
	if quit {
		return false
	}
	defer cc.queue.Done(cKey)

	err := cc.syncHandler(cKey.(string))
	if err == nil {
		cc.queue.Forget(cKey)
		return true
	}

	cc.queue.AddRateLimited(cKey)
	utilruntime.HandleError(fmt.Errorf("Sync %v failed with : %v", cKey, err))
	return true
}

func (cc *CertificateController) enqueueCertificateRequest(obj interface{}) {
	csr := obj.(*api.Secret)
	if csr.Annotations[instance.ResourceStatusKey] == "Ready" || csr.Annotations[instance.ResourceInstanceKey] != "" {
		return
	}
	if csr.Annotations[instance.ResourceTypeKey] == "csr" || csr.Annotations[instance.ResourceTypeKey] == "ca" {
		key, err := controller.KeyFunc(obj)
		if err != nil {
			utilruntime.HandleError(fmt.Errorf("Couldn't get key for object %+v: %v", obj, err))
			return
		}
		cc.queue.Add(key)
	}
}

// maybeSignCertificate will inspect the certificate request and, if it has
// been approved and meets policy expectations, generate an X509 cert using the
// cluster CA assets. If successful it will update the CSR approve subresource
// with the signed certificate.
func (cc *CertificateController) maybeSignCertificate(key string) error {
	startTime := time.Now()
	defer func() {
		glog.V(4).Infof("Finished syncing certificate request %q (%v)", key, time.Now().Sub(startTime))
	}()
	obj, exists, err := cc.csrStore.Indexer.GetByKey(key)
	if err != nil {
		return err
	}
	if !exists {
		glog.V(3).Infof("csr has been deleted: %v", key)
		return nil
	}
	secret := obj.(*api.Secret)

	switch secret.Annotations[instance.ResourceTypeKey] {
	case "csr":
		certControl := cc.certificateControl
		if caName := secret.Annotations[instance.ResourceCAKey]; caName != "" {
			caSecret, err := cc.kubeClient.Core().Secrets(cc.namespace).Get(caName)
			if err != nil {
				return fmt.Errorf("Failed to get ca certificate %s: %v", secret.Name, err)
			}
			certControl, err = instance.NewCertificateControlFromSecret(caSecret)
			if err != nil {
				return fmt.Errorf("Failed to initialize ca from secret %s: %v", caSecret.Name, err)
			}
		}
		err = certControl.GenerateCertificate(secret, nil)
		if err != nil {
			return fmt.Errorf("Failed to generate certificate %s: %v", secret.Name, err)
		}
	case "ca":
		err = cc.certificateControl.GenerateCA(secret)
		if err != nil {
			return fmt.Errorf("Failed to generate ca certificate %s: %v", secret.Name, err)
		}
	}
	_, err = cc.kubeClient.Core().Secrets(cc.namespace).Update(secret)
	if err != nil {
		return err
	}

	return nil
}
