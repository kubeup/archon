/*
Copyright 2015 The Kubernetes Authors.

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
	"fmt"
	"sync"
	"time"

	archoncache "kubeup.com/archon/pkg/cache"
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
)

const (
	// Interval of synchoronizing instance status from apiserver
	instanceSyncPeriod = 30 * time.Second

	// How long to wait before retrying the processing of a service change.
	// If this changes, the sleep in hack/jenkins/e2e.sh before downing a cluster
	// should be changed appropriately.
	minRetryDelay = 5 * time.Second
	maxRetryDelay = 300 * time.Second

	clientRetryCount    = 5
	clientRetryInterval = 5 * time.Second

	retryable    = true
	notRetryable = false

	doNotRetry = time.Duration(0)
)

type cachedInstance struct {
	// The cached state of the service
	state *cluster.Instance
	// Controls error back-off
	lastRetryDelay time.Duration
}

type instanceCache struct {
	mu          sync.Mutex // protects serviceMap
	instanceMap map[string]*cachedInstance
}

type InstanceController struct {
	cloud       cloudprovider.Interface
	kubeClient  clientset.Interface
	clusterName string
	archon      cloudprovider.ArchonInterface
	cache       *instanceCache
	// A store of services, populated by the serviceController
	instanceStore archoncache.StoreToInstanceLister
	// Watches changes to all services
	instanceController *cache.Controller
	eventBroadcaster   record.EventBroadcaster
	eventRecorder      record.EventRecorder
	// services that need to be synced
	workingQueue workqueue.DelayingInterface
}

// New returns a new service controller to keep cloud provider service resources
// (like load balancers) in sync with the registry.
func New(cloud cloudprovider.Interface, kubeClient clientset.Interface, clusterName string) (*InstanceController, error) {
	broadcaster := record.NewBroadcaster()
	broadcaster.StartRecordingToSink(&unversioned_core.EventSinkImpl{Interface: kubeClient.Core().Events("")})
	recorder := broadcaster.NewRecorder(api.EventSource{Component: "instance-controller"})

	if kubeClient != nil && kubeClient.Core().RESTClient().GetRateLimiter() != nil {
		metrics.RegisterMetricAndTrackRateLimiterUsage("instance_controller", kubeClient.Core().RESTClient().GetRateLimiter())
	}

	s := &InstanceController{
		cloud:            cloud,
		kubeClient:       kubeClient,
		clusterName:      clusterName,
		cache:            &instanceCache{instanceMap: make(map[string]*cachedInstance)},
		eventBroadcaster: broadcaster,
		eventRecorder:    recorder,
		workingQueue:     workqueue.NewDelayingQueue(),
	}
	s.instanceStore.Indexer, s.instanceController = cache.NewIndexerInformer(
		&cache.ListWatch{
			ListFunc: func(options api.ListOptions) (pkg_runtime.Object, error) {
				return s.kubeClient.Archon().Instances(api.NamespaceAll).List(options)
			},
			WatchFunc: func(options api.ListOptions) (watch.Interface, error) {
				return s.kubeClient.Archon().Instances(api.NamespaceAll).Watch(options)
			},
		},
		&cluster.Instance{},
		instanceSyncPeriod,
		cache.ResourceEventHandlerFuncs{
			AddFunc: s.enqueueInstance,
			UpdateFunc: func(old, cur interface{}) {
				// oldSvc, ok1 := old.(*cluster.Instance)
				// curSvc, ok2 := cur.(*cluster.Instance)
				// if ok1 && ok2 && s.needsUpdate(oldSvc, curSvc) {
				// s.enqueueInstance(cur)
				// }
			},
			DeleteFunc: s.enqueueInstance,
		},
		cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc},
	)
	if err := s.init(); err != nil {
		return nil, err
	}
	return s, nil
}

// obj could be an *api.Service, or a DeletionFinalStateUnknown marker item.
func (s *InstanceController) enqueueInstance(obj interface{}) {
	key, err := controller.KeyFunc(obj)
	if err != nil {
		glog.Errorf("Couldn't get key for object %#v: %v", obj, err)
		return
	}
	s.workingQueue.Add(key)
}

// Run starts a background goroutine that watches for changes to services that
// have (or had) LoadBalancers=true and ensures that they have
// load balancers created and deleted appropriately.
// serviceSyncPeriod controls how often we check the cluster's services to
// ensure that the correct load balancers exist.
// nodeSyncPeriod controls how often we check the cluster's nodes to determine
// if load balancers need to be updated to point to a new set.
//
// It's an error to call Run() more than once for a given ServiceController
// object.
func (s *InstanceController) Run(workers int) {
	defer runtime.HandleCrash()
	go s.instanceController.Run(wait.NeverStop)
	for i := 0; i < workers; i++ {
		go wait.Until(s.worker, time.Second, wait.NeverStop)
	}
}

// worker runs a worker thread that just dequeues items, processes them, and marks them done.
// It enforces that the syncHandler is never invoked concurrently with the same key.
func (s *InstanceController) worker() {
	for {
		func() {
			key, quit := s.workingQueue.Get()
			if quit {
				return
			}
			defer s.workingQueue.Done(key)
			err := s.syncInstance(key.(string))
			if err != nil {
				glog.Errorf("Error syncing service: %v", err)
			}
		}()
	}
}

func (s *InstanceController) init() error {
	if s.cloud == nil {
		return fmt.Errorf("WARNING: no cloud provider provided, services of type LoadBalancer will fail.")
	}

	archon, ok := s.cloud.Archon()
	if !ok {
		return fmt.Errorf("the cloud provider does not support archon.")
	}
	s.archon = archon
	return nil
}

// Returns an error if processing the service update failed, along with a time.Duration
// indicating whether processing should be retried; zero means no-retry; otherwise
// we should retry in that Duration.
func (s *InstanceController) processInstanceUpdate(cachedInstance *cachedInstance, instance *cluster.Instance, key string) (error, time.Duration) {

	// cache the service, we need the info for service deletion
	cachedInstance.state = instance
	err, retry := s.createInstanceIfNeeded(key, instance)
	if err != nil {
		message := "Error creating instance"
		if retry {
			message += " (will retry): "
		} else {
			message += " (will not retry): "
		}
		message += err.Error()
		s.eventRecorder.Event(instance, api.EventTypeWarning, "CreatingInstanceFailed", message)

		return err, cachedInstance.nextRetryDelay()
	}
	// Always update the cache upon success.
	// NOTE: Since we update the cached service if and only if we successfully
	// processed it, a cached service being nil implies that it hasn't yet
	// been successfully processed.
	s.cache.set(key, cachedInstance)

	cachedInstance.resetRetryDelay()
	return nil, doNotRetry
}

// Returns whatever error occurred along with a boolean indicator of whether it
// should be retried.
func (s *InstanceController) createInstanceIfNeeded(key string, instance *cluster.Instance) (error, bool) {
	// Check Network availability before creating instance
	network, err := s.kubeClient.Archon().Networks(instance.Namespace).Get(instance.Spec.NetworkName)

	if err != nil {
		return fmt.Errorf("Failed to get network %s: %v", instance.Spec.NetworkName, err), retryable
	}

	if network.Status.Phase != cluster.NetworkRunning {
		return fmt.Errorf("Network is not ready %s: %v", instance.Spec.NetworkName, network.Status.Phase), retryable
	}

	err = s.archon.AddNetworkAnnotation(s.clusterName, instance, network)
	if err != nil {
		return fmt.Errorf("Failed to add network annotation %s: %v", instance.Spec.NetworkName, err), retryable
	}

	users := make([]cluster.User, 0)

	for _, u := range instance.Spec.Users {
		user, err := s.kubeClient.Archon().Users(instance.Namespace).Get(u.Name)
		if err != nil {
			return fmt.Errorf("Failed to get user resource %s: %v", u.Name, err), retryable
		}
		users = append(users, *user)
	}

	deps := cluster.InstanceDependency{
		Network: *network,
		Users:   users,
	}

	instance.Dependency = deps

	// Note: It is safe to just call EnsureLoadBalancer.  But, on some clouds that requires a delete & create,
	// which may involve service interruption.  Also, we would like user-friendly events.

	// Save the state so we can avoid a write if it doesn't change
	previousState := instance.Status

	glog.V(2).Infof("Ensuring instance %s", key)

	// TODO: We could do a dry-run here if wanted to avoid the spurious cloud-calls & events when we restart

	// The load balancer doesn't exist yet, so create it.
	s.eventRecorder.Event(instance, api.EventTypeNormal, "CreatingInstance", "Creating instance")
	err = s.createInstance(instance)
	if err != nil {
		return fmt.Errorf("Failed to create instance %s: %v", key, err), retryable
	}
	s.eventRecorder.Event(instance, api.EventTypeNormal, "CreatedInstance", "Created instance")

	// Write the state if changed
	// TODO: Be careful here ... what if there were other changes to the service?
	if previousState != instance.Status {
		if err := s.persistUpdate(instance); err != nil {
			return fmt.Errorf("Failed to persist updated status to apiserver, even after retries. Giving up: %v", err), notRetryable
		}
	} else {
		glog.V(2).Infof("Not persisting unchanged LoadBalancerStatus to registry.")
	}

	return nil, notRetryable
}

func (s *InstanceController) persistUpdate(instance *cluster.Instance) error {
	var err error
	for i := 0; i < clientRetryCount; i++ {
		_, err = s.kubeClient.Archon().Instances(instance.Namespace).UpdateStatus(instance)
		if err == nil {
			return nil
		}
		// If the object no longer exists, we don't want to recreate it. Just bail
		// out so that we can process the delete, which we should soon be receiving
		// if we haven't already.
		if errors.IsNotFound(err) {
			glog.Infof("Not persisting update to service '%s/%s' that no longer exists: %v",
				instance.Namespace, instance.Name, err)
			return nil
		}
		// TODO: Try to resolve the conflict if the change was unrelated to load
		// balancer status. For now, just rely on the fact that we'll
		// also process the update that caused the resource version to change.
		if errors.IsConflict(err) {
			glog.V(4).Infof("Not persisting update to instance '%s/%s' that has been changed since we received it: %v",
				instance.Namespace, instance.Name, err)
			return nil
		}
		glog.Warningf("Failed to persist updated InstanceStatus to instance '%s/%s' after creating its load balancer: %v",
			instance.Namespace, instance.Name, err)
		time.Sleep(clientRetryInterval)
	}
	return err
}

func (s *InstanceController) createInstance(instance *cluster.Instance) error {
	// - Only one protocol supported per service
	// - Not all cloud providers support all protocols and the next step is expected to return
	//   an error for unsupported protocols
	status, err := s.archon.EnsureInstance(s.clusterName, instance)
	if err != nil {
		return err
	} else {
		instance.Status = *status
	}

	return nil
}

// ListKeys implements the interface required by DeltaFIFO to list the keys we
// already know about.
func (s *instanceCache) ListKeys() []string {
	s.mu.Lock()
	defer s.mu.Unlock()
	keys := make([]string, 0, len(s.instanceMap))
	for k := range s.instanceMap {
		keys = append(keys, k)
	}
	return keys
}

// GetByKey returns the value stored in the serviceMap under the given key
func (s *instanceCache) GetByKey(key string) (interface{}, bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if v, ok := s.instanceMap[key]; ok {
		return v, true, nil
	}
	return nil, false, nil
}

// ListKeys implements the interface required by DeltaFIFO to list the keys we
// already know about.
func (s *instanceCache) allInstances() []*cluster.Instance {
	s.mu.Lock()
	defer s.mu.Unlock()
	instances := make([]*cluster.Instance, 0, len(s.instanceMap))
	for _, v := range s.instanceMap {
		instances = append(instances, v.state)
	}
	return instances
}

func (s *instanceCache) get(instanceName string) (*cachedInstance, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	instance, ok := s.instanceMap[instanceName]
	return instance, ok
}

func (s *instanceCache) getOrCreate(instanceName string) *cachedInstance {
	s.mu.Lock()
	defer s.mu.Unlock()
	instance, ok := s.instanceMap[instanceName]
	if !ok {
		instance = &cachedInstance{}
		s.instanceMap[instanceName] = instance
	}
	return instance
}

func (s *instanceCache) set(instanceName string, instance *cachedInstance) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.instanceMap[instanceName] = instance
}

func (s *instanceCache) delete(instanceName string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.instanceMap, instanceName)
}

func (s *InstanceController) needsUpdate(oldInstance *cluster.Instance, newInstance *cluster.Instance) bool {
	return true
}

// Computes the next retry, using exponential backoff
// mutex must be held.
func (s *cachedInstance) nextRetryDelay() time.Duration {
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
func (s *cachedInstance) resetRetryDelay() {
	s.lastRetryDelay = time.Duration(0)
}

// syncService will sync the Service with the given key if it has had its expectations fulfilled,
// meaning it did not expect to see any more of its pods created or deleted. This function is not meant to be
// invoked concurrently with the same key.
func (s *InstanceController) syncInstance(key string) error {
	startTime := time.Now()
	var cachedInstance *cachedInstance
	var retryDelay time.Duration
	defer func() {
		glog.V(4).Infof("Finished syncing service %q (%v)", key, time.Now().Sub(startTime))
	}()
	// obj holds the latest service info from apiserver
	obj, exists, err := s.instanceStore.Indexer.GetByKey(key)
	if err != nil {
		glog.Infof("Unable to retrieve service %v from store: %v", key, err)
		s.workingQueue.Add(key)
		return err
	}
	if !exists {
		// service absence in store means watcher caught the deletion, ensure LB info is cleaned
		glog.Infof("Instance has been deleted %v", key)
		err, retryDelay = s.processInstanceDeletion(key)
	} else {
		instance, ok := obj.(*cluster.Instance)
		if ok {
			cachedInstance = s.cache.getOrCreate(key)
			err, retryDelay = s.processInstanceUpdate(cachedInstance, instance, key)
		} else {
			tombstone, ok := obj.(cache.DeletedFinalStateUnknown)
			if !ok {
				return fmt.Errorf("object contained wasn't a service or a deleted key: %#v", obj)
			}
			glog.Infof("Found tombstone for %v", key)
			err, retryDelay = s.processInstanceDeletion(tombstone.Key)
		}
	}

	if retryDelay != 0 {
		// Add the failed service back to the queue so we'll retry it.
		glog.Errorf("Failed to process service. Retrying in %s: %v", retryDelay, err)
		go func(obj interface{}, delay time.Duration) {
			// put back the service key to working queue, it is possible that more entries of the service
			// were added into the queue during the delay, but it does not mess as when handling the retry,
			// it always get the last service info from service store
			s.workingQueue.AddAfter(obj, delay)
		}(key, retryDelay)
	} else if err != nil {
		runtime.HandleError(fmt.Errorf("Failed to process service. Not retrying: %v", err))
	}
	return nil
}

// Returns an error if processing the service deletion failed, along with a time.Duration
// indicating whether processing should be retried; zero means no-retry; otherwise
// we should retry after that Duration.
func (s *InstanceController) processInstanceDeletion(key string) (error, time.Duration) {
	cachedInstance, ok := s.cache.get(key)
	if !ok {
		return fmt.Errorf("Service %s not in cache even though the watcher thought it was. Ignoring the deletion.", key), doNotRetry
	}
	instance := cachedInstance.state

	s.eventRecorder.Event(instance, api.EventTypeNormal, "DeletingInstance", "Deleting instance")
	err := s.archon.EnsureInstanceDeleted(s.clusterName, instance)
	if err != nil {
		message := "Error deleting instance (will retry): " + err.Error()
		s.eventRecorder.Event(instance, api.EventTypeWarning, "DeletingInstanceFailed", message)
		return err, cachedInstance.nextRetryDelay()
	}
	s.eventRecorder.Event(instance, api.EventTypeNormal, "DeletedInstance", "Deleted instance")
	s.cache.delete(key)

	cachedInstance.resetRetryDelay()
	return nil, doNotRetry
}
