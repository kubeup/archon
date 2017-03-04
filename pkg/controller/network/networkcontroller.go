/*
Copyright 2016 The Kubernetes Authors.
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

// This file is modified from servicecontroller.go in original kubernetes source tree

package network

import (
	"fmt"
	"reflect"
	"sync"
	"time"

	//archoncache "kubeup.com/archon/pkg/cache"
	"kubeup.com/archon/pkg/clientset"
	"kubeup.com/archon/pkg/cloudprovider"
	"kubeup.com/archon/pkg/cluster"

	"github.com/golang/glog"
	"k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/api/errors"
	"k8s.io/kubernetes/pkg/client/cache"
	unversioned_core "k8s.io/kubernetes/pkg/client/clientset_generated/internalclientset/typed/core/internalversion"
	"k8s.io/kubernetes/pkg/client/record"
	"k8s.io/kubernetes/pkg/controller"
	pkg_runtime "k8s.io/kubernetes/pkg/runtime"
	"k8s.io/kubernetes/pkg/util/metrics"
	"k8s.io/kubernetes/pkg/util/runtime"
	"k8s.io/kubernetes/pkg/util/wait"
	"k8s.io/kubernetes/pkg/util/workqueue"
	"k8s.io/kubernetes/pkg/watch"
	"kubeup.com/archon/pkg/util"
)

const (
	// Interval of synchoronizing network status from apiserver
	networkSyncPeriod = 30 * time.Second

	// How long to wait before retrying the processing of a network change.
	minRetryDelay = 5 * time.Second
	maxRetryDelay = 300 * time.Second

	clientRetryCount    = 5
	clientRetryInterval = 5 * time.Second

	retryable    = true
	notRetryable = false

	doNotRetry = time.Duration(0)
)

type cachedNetwork struct {
	// The cached state of the network
	state *cluster.Network
	// Controls error back-off
	lastRetryDelay time.Duration
}

type networkCache struct {
	mu         sync.Mutex // protects networkMap
	networkMap map[string]*cachedNetwork
}

type NetworkController struct {
	cloud       cloudprovider.Interface
	kubeClient  clientset.Interface
	clusterName string
	namespace   string
	archon      cloudprovider.ArchonInterface
	cache       *networkCache
	// A store of networks, populated by the networkController
	networkStore cache.Indexer
	// Watches changes to all networks
	networkController *cache.Controller
	eventBroadcaster  record.EventBroadcaster
	eventRecorder     record.EventRecorder
	// networks that need to be synced
	workingQueue workqueue.DelayingInterface
}

// New returns a new network controller to keep cloud provider network resource
// in sync with registry
// TODO: Currently only registry-to-cloud sync is performed. Is it necessary to
// perform the opposite sync? is network subject to change after creation?
func New(cloud cloudprovider.Interface, kubeClient clientset.Interface, clusterName string, namespace string) (*NetworkController, error) {
	broadcaster := record.NewBroadcaster()
	broadcaster.StartRecordingToSink(&unversioned_core.EventSinkImpl{Interface: kubeClient.Core().Events("")})
	recorder := broadcaster.NewRecorder(api.EventSource{Component: "network-controller"})

	if kubeClient != nil && kubeClient.Core().RESTClient().GetRateLimiter() != nil {
		metrics.RegisterMetricAndTrackRateLimiterUsage("network_controller", kubeClient.Core().RESTClient().GetRateLimiter())
	}

	s := &NetworkController{
		cloud:            cloud,
		kubeClient:       kubeClient,
		clusterName:      clusterName,
		namespace:        namespace,
		cache:            &networkCache{networkMap: make(map[string]*cachedNetwork)},
		eventBroadcaster: broadcaster,
		eventRecorder:    recorder,
		workingQueue:     workqueue.NewDelayingQueue(),
	}
	s.networkStore, s.networkController = cache.NewIndexerInformer(
		&cache.ListWatch{
			ListFunc: func(options api.ListOptions) (pkg_runtime.Object, error) {
				return s.kubeClient.Archon().Networks(namespace).List()
			},
			WatchFunc: func(options api.ListOptions) (watch.Interface, error) {
				return s.kubeClient.Archon().Networks(namespace).Watch()
			},
		},
		&cluster.Network{},
		networkSyncPeriod,
		cache.ResourceEventHandlerFuncs{
			AddFunc: s.enqueueNetwork,
			UpdateFunc: func(old, cur interface{}) {
				n, _ := cur.(*cluster.Network)
				if n != nil {
					glog.V(2).Infof("Network changes in registry are ignored: %v", n.Name)
				}
				// TODO?
				// oldSvc, ok1 := old.(*cluster.network)
				// curSvc, ok2 := cur.(*cluster.network)
				// if ok1 && ok2 && s.needsUpdate(oldSvc, curSvc) {
				// s.enqueuenetwork(cur)
				// }
			},
			DeleteFunc: s.enqueueNetwork,
		},
		cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc},
	)
	if err := s.init(); err != nil {
		return nil, err
	}
	return s, nil
}

// obj could be an *api.network, or a DeletionFinalStateUnknown marker item.
func (s *NetworkController) enqueueNetwork(obj interface{}) {
	key, err := controller.KeyFunc(obj)
	if err != nil {
		glog.Errorf("Couldn't get key for object %#v: %v", obj, err)
		return
	}
	s.workingQueue.Add(key)
}

// Run starts a background goroutine that watches for changes of networks and
// ensures they are created/deleted properly
// networkSyncPeriod controls how often we resync our store with cluster registry
//
// It's an error to call Run() more than once for a given networkController
// object.
func (s *NetworkController) Run(workers int) {
	defer runtime.HandleCrash()
	go s.networkController.Run(wait.NeverStop)
	for i := 0; i < workers; i++ {
		go wait.Until(s.worker, time.Second, wait.NeverStop)
	}
}

// worker runs a worker thread that just dequeues items, processes them, and marks them done.
// It enforces that the syncHandler is never invoked concurrently with the same key.
func (s *NetworkController) worker() {
	for {
		func() {
			key, quit := s.workingQueue.Get()
			if quit {
				return
			}
			defer s.workingQueue.Done(key)
			err := s.syncNetwork(key.(string))
			if err != nil {
				glog.Errorf("Error syncing network: %v", err)
			}
		}()
	}
}

func (s *NetworkController) init() error {
	if s.cloud == nil {
		return fmt.Errorf("WARNING: no cloud provider provided, networks sync will fail.")
	}

	archon, ok := s.cloud.Archon()
	if !ok {
		return fmt.Errorf("the cloud provider does not support archon.")
	}
	s.archon = archon
	return nil
}

// Returns an error if processing the network update failed, along with a time.Duration
// indicating whether processing should be retried; zero means no-retry; otherwise
// we should retry in that Duration.
func (s *NetworkController) processNetworkUpdate(cachedNetwork *cachedNetwork, network *cluster.Network, key string) (error, time.Duration) {

	// cache the network, we need the info for network deletion
	cachedNetwork.state = network
	err, retry := s.createNetworkIfNeeded(key, network)
	if err != nil {
		message := "Error creating network"
		if retry {
			message += " (will retry): "
		} else {
			message += " (will not retry): "
		}
		message += err.Error()
		s.eventRecorder.Event(network, api.EventTypeWarning, "CreatingNetworkFailed", message)

		return err, cachedNetwork.nextRetryDelay()
	}
	// Always update the cache upon success.
	// NOTE: Since we update the cached network if and only if we successfully
	// processed it, a cached network being nil implies that it hasn't yet
	// been successfully processed.
	s.cache.set(key, cachedNetwork)

	cachedNetwork.resetRetryDelay()
	return nil, doNotRetry
}

// Returns whatever error occurred along with a boolean indicator of whether it
// should be retried.
func (s *NetworkController) createNetworkIfNeeded(key string, network *cluster.Network) (error, bool) {
	// Check if network is already there
	if network.Status.Phase == cluster.NetworkRunning {
		return nil, notRetryable
	}

	previousState := *cluster.NetworkStatusDeepCopy(&network.Status)
	previousAnnotations := make(map[string]string)
	util.MapCopy(previousAnnotations, network.Annotations)

	glog.V(2).Infof("Ensuring network %s", key)

	s.eventRecorder.Event(network, api.EventTypeNormal, "CreatingNetwork", "Creating network")
	err := s.createNetwork(network)
	if err != nil {
		return fmt.Errorf("Failed to create network %s: %v", key, err), retryable
	}
	s.eventRecorder.Event(network, api.EventTypeNormal, "CreatedNetwork", "Created network")

	// Write the state if changed
	if !reflect.DeepEqual(previousState, network.Status) || !reflect.DeepEqual(previousAnnotations, network.Annotations) {
		if err := s.persistUpdate(network); err != nil {
			return fmt.Errorf("Failed to persist updated status to apiserver, even after retries. Giving up: %v", err), notRetryable
		}
	} else {
		glog.V(2).Infof("Not persisting unchanged NetworkStatus to registry.")
	}

	return nil, notRetryable
}

func (s *NetworkController) persistUpdate(network *cluster.Network) error {
	var err error
	for i := 0; i < clientRetryCount; i++ {
		_, err = s.kubeClient.Archon().Networks(network.Namespace).Update(network)
		if err == nil {
			return nil
		}
		// If the object no longer exists, we don't want to recreate it. Just bail
		// out so that we can process the delete, which we should soon be receiving
		// if we haven't already.
		if errors.IsNotFound(err) {
			glog.Infof("Not persisting update to network '%s/%s' that no longer exists: %v",
				network.Namespace, network.Name, err)
			return nil
		}
		// TODO: Try to resolve the conflict if the change was unrelated to network
		// status.
		if errors.IsConflict(err) {
			glog.V(4).Infof("Not persisting update to network '%s/%s' that has been changed since we received it: %v",
				network.Namespace, network.Name, err)
			return nil
		}
		glog.Warningf("Failed to persist updated NetworkStatus to network '%s/%s' after creating: %v",
			network.Namespace, network.Name, err)
		time.Sleep(clientRetryInterval)
	}
	return err
}

func (s *NetworkController) createNetwork(network *cluster.Network) error {
	// - Only one protocol supported per network
	// - Not all cloud providers support all protocols and the next step is expected to return
	//   an error for unsupported protocols
	status, err := s.archon.EnsureNetwork(s.clusterName, network)
	if err != nil {
		return err
	} else {
		network.Status = *status
	}

	return nil
}

func (s *networkCache) ListKeys() []string {
	s.mu.Lock()
	defer s.mu.Unlock()
	keys := make([]string, 0, len(s.networkMap))
	for k := range s.networkMap {
		keys = append(keys, k)
	}
	return keys
}

func (s *networkCache) GetByKey(key string) (interface{}, bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if v, ok := s.networkMap[key]; ok {
		return v, true, nil
	}
	return nil, false, nil
}

func (s *networkCache) allNetworks() []*cluster.Network {
	s.mu.Lock()
	defer s.mu.Unlock()
	networks := make([]*cluster.Network, 0, len(s.networkMap))
	for _, v := range s.networkMap {
		networks = append(networks, v.state)
	}
	return networks
}

func (s *networkCache) get(networkName string) (*cachedNetwork, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	network, ok := s.networkMap[networkName]
	return network, ok
}

func (s *networkCache) getOrCreate(networkName string) *cachedNetwork {
	s.mu.Lock()
	defer s.mu.Unlock()
	network, ok := s.networkMap[networkName]
	if !ok {
		network = &cachedNetwork{}
		s.networkMap[networkName] = network
	}
	return network
}

func (s *networkCache) set(networkName string, network *cachedNetwork) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.networkMap[networkName] = network
}

func (s *networkCache) delete(networkName string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.networkMap, networkName)
}

func (s *NetworkController) needsUpdate(oldNetwork *cluster.Network, newNetwork *cluster.Network) bool {
	return true
}

// Computes the next retry, using exponential backoff
// mutex must be held.
func (s *cachedNetwork) nextRetryDelay() time.Duration {
	s.lastRetryDelay = s.lastRetryDelay * 2
	if s.lastRetryDelay < minRetryDelay {
		s.lastRetryDelay = minRetryDelay
	}
	if s.lastRetryDelay > maxRetryDelay {
		s.lastRetryDelay = maxRetryDelay
	}
	return s.lastRetryDelay
}

// Resets the retry exponential backoff.  mutex must be held.
func (s *cachedNetwork) resetRetryDelay() {
	s.lastRetryDelay = time.Duration(0)
}

// syncNetwork will sync the network with the given key if it has had its expectations fulfilled,
func (s *NetworkController) syncNetwork(key string) error {
	startTime := time.Now()
	var cachedNetwork *cachedNetwork
	var retryDelay time.Duration
	defer func() {
		glog.V(4).Infof("Finished syncing network %q (%v)", key, time.Now().Sub(startTime))
	}()
	// obj holds the latest network info from apiserver
	obj, exists, err := s.networkStore.GetByKey(key)
	if err != nil {
		glog.Infof("Unable to retrieve network %v from store: %v", key, err)
		s.workingQueue.Add(key)
		return err
	}
	if !exists {
		glog.Infof("network has been deleted %v", key)
		err, retryDelay = s.processNetworkDeletion(key)
	} else {
		network, ok := obj.(*cluster.Network)
		if ok {
			cachedNetwork = s.cache.getOrCreate(key)
			err, retryDelay = s.processNetworkUpdate(cachedNetwork, network, key)
		} else {
			tombstone, ok := obj.(cache.DeletedFinalStateUnknown)
			if !ok {
				return fmt.Errorf("object contained wasn't a network or a deleted key: %#v", obj)
			}
			glog.Infof("Found tombstone for %v", key)
			err, retryDelay = s.processNetworkDeletion(tombstone.Key)
		}
	}

	if retryDelay != 0 {
		// Add the failed network back to the queue so we'll retry it.
		glog.Errorf("Failed to process network. Retrying in %s: %v", retryDelay, err)
		go func(obj interface{}, delay time.Duration) {
			// put back the network key to working queue, it is possible that more entries of the network
			// were added into the queue during the delay, but it does not mess as when handling the retry,
			// it always get the last network info from network store
			s.workingQueue.AddAfter(obj, delay)
		}(key, retryDelay)
	} else if err != nil {
		runtime.HandleError(fmt.Errorf("Failed to process network. Not retrying: %v", err))
	}
	return nil
}

// Returns an error if processing the network deletion failed, along with a time.Duration
// indicating whether processing should be retried; zero means no-retry; otherwise
// we should retry after that Duration.
func (s *NetworkController) processNetworkDeletion(key string) (error, time.Duration) {
	cachedNetwork, ok := s.cache.get(key)
	if !ok {
		return fmt.Errorf("network %s not in cache even though the watcher thought it was. Ignoring the deletion.", key), doNotRetry
	}
	network := cachedNetwork.state

	s.eventRecorder.Event(network, api.EventTypeNormal, "DeletingNetwork", "Deleting network")
	err := s.archon.EnsureNetworkDeleted(s.clusterName, network)
	if err != nil {
		message := "Error deleting network (will retry): " + err.Error()
		s.eventRecorder.Event(network, api.EventTypeWarning, "DeletingNetworkFailed", message)
		return err, cachedNetwork.nextRetryDelay()
	}
	s.eventRecorder.Event(network, api.EventTypeNormal, "DeletedNetwork", "Deleted network")
	s.cache.delete(key)

	cachedNetwork.resetRetryDelay()
	return nil, doNotRetry
}
