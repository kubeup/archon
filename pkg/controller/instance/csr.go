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

package instance

import (
	"github.com/golang/glog"
	"k8s.io/kubernetes/pkg/api"
	"kubeup.com/archon/pkg/clientset"
	"kubeup.com/archon/pkg/cluster"
	"kubeup.com/archon/pkg/initializer"

	"fmt"
)

var (
	ResourceStatusKey   = "archon.kubeup.com/status"
	ResourceTypeKey     = "archon.kubeup.com/type"
	ResourceInstanceKey = "archon.kubeup.com/instance"
	CSRToken            = "csr"
)

type CSRInitializer struct {
	certificateControl CertificateControlInterface
	kubeClient         clientset.Interface
}

var _ initializer.Initializer = &CSRInitializer{}

func NewCSRInitializer(kubeClient clientset.Interface, caCertFile, caKeyFile string) (initializer.Initializer, error) {
	certControl, err := NewCertificateControl(caCertFile, caKeyFile)
	if err != nil {
		glog.Errorf("WARNING: Unable to start certificate controller: %s", err.Error())
		//return
	}

	c := &CSRInitializer{
		certificateControl: certControl,
		kubeClient:         kubeClient,
	}
	return c, nil
}

func (ci *CSRInitializer) Token() string {
	return CSRToken
}

func (ci *CSRInitializer) Initialize(obj initializer.Object) (updatedObj initializer.Object, err error, retryable bool) {
	instance, ok := obj.(*cluster.Instance)
	if !ok {
		err = fmt.Errorf("expecting Instance. got %v", obj)
		return
	}

	if initializer.HasInitializer(instance, PublicIPToken, PrivateIPToken) {
		err = initializer.ErrSkip
		return
	}

	var secret *api.Secret
	for _, n := range instance.Spec.Secrets {
		secret, err = ci.kubeClient.Core().Secrets(instance.Namespace).Get(n.Name)
		if err != nil {
			err = fmt.Errorf("Failed to get secret resource %s: %v", n.Name, err)
			return
		}
		if status, ok := secret.Annotations[ResourceStatusKey]; ok {
			if status != "Ready" {
				switch secret.Annotations[ResourceTypeKey] {
				case "csr":
					if secret.Annotations[ResourceInstanceKey] != instance.Name {
						err = fmt.Errorf("Failed to generate certificate. CSR doesn't belong to this instance")
						return
					}
					if ci.certificateControl == nil {
						err = fmt.Errorf("Failed to generated certificate. Certficate control is not there")
						return
					}
					err = ci.certificateControl.GenerateCertificate(secret, instance)
					if err != nil {
						err = fmt.Errorf("Failed to generate certificate %s: %v", n.Name, err)
						return
					}
					_, err = ci.kubeClient.Core().Secrets(instance.Namespace).Update(secret)
				}
			}
		}
	}

	initializer.RemoveInitializer(instance, CSRToken)
	updatedObj, err = ci.kubeClient.Archon().Instances(instance.Namespace).Update(instance)
	if err != nil {
		retryable = true
	}
	return
}

func (ci *CSRInitializer) Finalize(obj initializer.Object) (updatedObj initializer.Object, err error, retryable bool) {
	return
}
