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

package clientset

import (
	"kubeup.com/archon/pkg/clientset/archon"

	"github.com/golang/glog"
	restclient "k8s.io/client-go/rest"
	kubernetes "k8s.io/kubernetes/pkg/client/clientset_generated/clientset"

	_ "kubeup.com/archon/pkg/clientset/archon/install"
)

type Interface interface {
	kubernetes.Interface
	Archon() archon.ArchonInterface
}

type Clientset struct {
	kubernetes.Clientset
	*archon.ArchonClient

	config *restclient.Config
}

// Archon retrieves the ArchonClient
func (c *Clientset) Archon() archon.ArchonInterface {
	if c == nil {
		return nil
	}
	return c.ArchonClient
}

// NewForConfig creates a new Clientset for the given config.
func NewForConfig(c *restclient.Config) (*Clientset, error) {
	var clientset Clientset
	var err error
	cs, err := kubernetes.NewForConfig(c)
	if err != nil {
		return nil, err
	}
	clientset.config = c
	clientset.Clientset = *cs
	clientset.ArchonClient, err = archon.NewForConfig(c)
	if err != nil {
		glog.Errorf("failed to create the ArchonClient: %v", err)
		return nil, err
	}
	return &clientset, nil
}

// NewForConfigOrDie creates a new Clientset for the given config and
// panics if there is an error in the config.
func NewForConfigOrDie(c *restclient.Config) *Clientset {
	var clientset Clientset
	clientset.Clientset = *kubernetes.NewForConfigOrDie(c)
	clientset.ArchonClient = archon.NewForConfigOrDie(c)
	return &clientset
}

// New creates a new Clientset for the given RESTClient.
func New(c *restclient.RESTClient) *Clientset {
	var clientset Clientset
	clientset.Clientset = *kubernetes.New(c)
	clientset.ArchonClient = archon.New(c)
	return &clientset
}
