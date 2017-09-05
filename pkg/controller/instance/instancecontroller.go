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

package instance

import (
	"fmt"
	"sync"
	"time"

	"github.com/golang/glog"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	pkg_runtime "k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	watch "k8s.io/apimachinery/pkg/watch"
	v1core "k8s.io/client-go/kubernetes/typed/core/v1"
	clientv1 "k8s.io/client-go/pkg/api/v1"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/record"
	"k8s.io/client-go/util/workqueue"
	"k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/api/v1"
	"k8s.io/kubernetes/pkg/controller"
	"k8s.io/kubernetes/pkg/util/metrics"
	archoncache "kubeup.com/archon/pkg/cache"
	"kubeup.com/archon/pkg/clientset"
	"kubeup.com/archon/pkg/cloudprovider"
	"kubeup.com/archon/pkg/cluster"
	archoncontroller "kubeup.com/archon/pkg/controller"
	"kubeup.com/archon/pkg/controller/certificate"
	"kubeup.com/archon/pkg/initializer"
	"kubeup.com/archon/pkg/util"
)

const (
	// Interval of synchoronizing instance status from apiserver and cloudprovider
	instanceSyncPeriod = 30 * time.Second
	instanceSyncJitter = 0.1

	minRetryDelay = 5 * time.Second
	maxRetryDelay = 300 * time.Second

	clientRetryCount    = 5
	clientRetryInterval = 5 * time.Second

	retryable    = true
	notRetryable = false

	doNotRetry = time.Duration(0)
)

type cachedInstance struct {
	// The cached state of the instance
	state *cluster.Instance
	// Controls error back-off
	lastRetryDelay time.Duration
	// Lock held by sync routine either with cloud or with k8s
	mu *util.Mutex
}

type instanceCache struct {
	mu          sync.Mutex // protects instanceMap
	instanceMap map[string]*cachedInstance
}

type InstanceController struct {
	cloud       cloudprovider.Interface
	kubeClient  clientset.Interface
	clusterName string
	namespace   string
	archon      cloudprovider.ArchonInterface
	cache       *instanceCache
	// Initializers
	initializerManager *initializer.InitializerManager
	// A store of instances, populated by the instanceController
	instanceStore archoncache.StoreToInstanceLister
	// Watches changes to all instances
	instanceController cache.Controller
	instanceControl    archoncontroller.InstanceControlInterface
	eventBroadcaster   record.EventBroadcaster
	eventRecorder      record.EventRecorder
	// instances that need to be synced
	workingQueue workqueue.DelayingInterface
}

// New returns a new instance controller to keep cloud provider instance resources
// in sync with the registry.
func New(cloud cloudprovider.Interface, kubeClient clientset.Interface, clusterName, namespace string) (*InstanceController, error) {
	broadcaster := record.NewBroadcaster()
	broadcaster.StartRecordingToSink(&v1core.EventSinkImpl{Interface: v1core.New(kubeClient.Core().RESTClient()).Events("")})
	recorder := broadcaster.NewRecorder(api.Scheme, clientv1.EventSource{Component: "instance-controller"})

	if kubeClient != nil && kubeClient.Core().RESTClient().GetRateLimiter() != nil {
		metrics.RegisterMetricAndTrackRateLimiterUsage("instance_controller", kubeClient.Core().RESTClient().GetRateLimiter())
	}

	// Initializers
	csrInit, err := NewCSRInitializer(kubeClient)
	if err != nil {
		return nil, err
	}
	publicIpInit, err := NewPublicIPInitializer(kubeClient, cloud, clusterName)
	if err != nil {
		return nil, err
	}
	privateIpInit, err := NewPrivateIPInitializer(kubeClient, cloud, clusterName)
	if err != nil {
		return nil, err
	}
	inits := []initializer.Initializer{
		publicIpInit,
		privateIpInit,
		csrInit,
	}
	manager, err := initializer.NewInitializerManager(inits, kubeClient)
	if err != nil {
		return nil, err
	}

	s := &InstanceController{
		cloud:              cloud,
		kubeClient:         kubeClient,
		initializerManager: manager,
		instanceControl: archoncontroller.RealInstanceControl{
			KubeClient: kubeClient,
			Recorder:   recorder,
		},
		clusterName:      clusterName,
		namespace:        namespace,
		cache:            &instanceCache{instanceMap: make(map[string]*cachedInstance)},
		eventBroadcaster: broadcaster,
		eventRecorder:    recorder,
		workingQueue:     workqueue.NewDelayingQueue(),
	}
	s.instanceStore.Indexer, s.instanceController = cache.NewIndexerInformer(
		&cache.ListWatch{
			ListFunc: func(options metav1.ListOptions) (pkg_runtime.Object, error) {
				options.IncludeUninitialized = true
				return s.kubeClient.Archon().Instances(namespace).List(options)
			},
			WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
				options.IncludeUninitialized = true
				return s.kubeClient.Archon().Instances(namespace).Watch(options)
			},
		},
		&cluster.Instance{},
		instanceSyncPeriod,
		cache.ResourceEventHandlerFuncs{
			AddFunc: s.enqueueInstance,
			UpdateFunc: func(old, cur interface{}) {
				oldInstance, ok1 := old.(*cluster.Instance)
				curInstance, ok2 := cur.(*cluster.Instance)
				if ok1 && ok2 && s.needsUpdate(oldInstance, curInstance) {
					s.enqueueInstance(cur)
				}
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

// obj could be an *api.instance, or a DeletionFinalStateUnknown marker item.
func (s *InstanceController) enqueueInstance(obj interface{}) {
	key, err := controller.KeyFunc(obj)
	if err != nil {
		glog.Errorf("Couldn't get key for object %#v: %v", obj, err)
		return
	}
	s.workingQueue.Add(key)
}

// Run starts a watcher who watches changes of instances in registry, serveral workers
// who process instance updates, and a cloud instance syncer who syncs instance statuses
// from cloud back to registry
func (s *InstanceController) Run(workers int) {
	defer runtime.HandleCrash()
	go s.instanceController.Run(wait.NeverStop)
	go wait.JitterUntil(s.syncCloudInstances, instanceSyncPeriod, instanceSyncJitter, true, wait.NeverStop)
	for i := 0; i < workers; i++ {
		go wait.Until(s.worker, time.Second, wait.NeverStop)
	}
}

func (s *InstanceController) syncCloudInstances() {
	glog.V(4).Infof("Syncing cloud instances")

	// List all networks in our working namespace
	networks, err := s.kubeClient.Archon().Networks(s.namespace).List()
	if err != nil {
		glog.Errorf("Couldn't get all networks from k8s. Cloud instance sync failed: %+v", err)
		return
	}

	// List all instances in all networks in the cloud
	cloudStatus := make(map[string]*cluster.InstanceStatus)
	for _, n := range networks.Items {
		names, statuses, err := s.archon.ListInstances(s.clusterName, &n, nil)
		if err != nil {
			glog.Errorf("Couldn't get instances for network %s. Skipping: %+v", n.Name, err)
			continue
		}

		for i, name := range names {
			cloudStatus[name] = statuses[i]
		}
	}
	glog.Infof("Cloud statuses: %+v", cloudStatus)

	// List all instances in the registry
	instances, err := s.kubeClient.Archon().Instances(s.namespace).List(metav1.ListOptions{})
	if err != nil {
		glog.Errorf("Couldn't get instances from k8s. Cloud instance sync failed: %+v", err)
		return
	}

	// Update the statuses of existing instances in the registry
	for _, instance := range instances.Items {
		key, err := controller.KeyFunc(instance)
		if err != nil {
			glog.Errorf("Couldn't get key for instance %#v: %v", instance, err)
			continue
		}

		glog.Infof("Registry instance key: %v, name: %v", key, instance.Name)
		status, ok := cloudStatus[instance.Name]
		if !ok {
			// Exists in registry but not in the cloud. Either waiting to be created or
			// terminated unexpectedly.
			switch instance.Status.Phase {
			case cluster.InstanceInitializing, cluster.InstancePending:
				continue
			default:
				status = cluster.InstanceStatusDeepCopy(&instance.Status)
				status.Phase = cluster.InstanceUnknown
			}
		}

		if !cluster.InstanceStatusEqual(*status, instance.Status) {
			// Update only when there's no current operation going on with this instance
			cachedInstance := s.cache.getOrCreate(key)
			if cachedInstance.mu.TryLock() {
				glog.V(4).Infof("Syncing instance %s from cloud", key)
				instance.Status = *status
				updatedInstance, err := s.kubeClient.Archon().Instances(instance.Namespace).Update(&instance)
				if err != nil {
					glog.Errorf("Couldn't update instance %s in k8s: %+v", key, err)
					continue
				}

				cachedInstance.state = updatedInstance
				s.cache.set(key, cachedInstance)
				cachedInstance.mu.Unlock()
			} else {
				glog.Errorf("Couldn't update instance %s. Coz there's an ongoing operation", key)
			}
		} else {
			glog.V(4).Infof("Ignore syncing for instance %s. Status unchanged", key)
		}
	}

	glog.V(4).Infof("Done syncing cloud instances.")
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
				glog.Errorf("Error syncing instance: %v", err)
			}
		}()
	}
}

func (s *InstanceController) init() error {
	if s.cloud == nil {
		return fmt.Errorf("WARNING: no cloud provider provided, InstanceController won't work.")
	}

	archon, ok := s.cloud.Archon()
	if !ok {
		return fmt.Errorf("the cloud provider does not support archon.")
	}
	s.archon = archon
	return nil
}

func (s *InstanceController) ensureDependency(key string, instance *cluster.Instance) (error, *cluster.InstanceDependency) {
	// Check Network availability before creating instance
	network, err := s.kubeClient.Archon().Networks(instance.Namespace).Get(instance.Spec.NetworkName)
	if err != nil {
		return fmt.Errorf("Failed to get network %s: %v", instance.Spec.NetworkName, err), nil
	}

	if network.Status.Phase != cluster.NetworkRunning {
		return fmt.Errorf("Network is not ready %s: %v", instance.Spec.NetworkName, network.Status.Phase), nil
	}

	users := make([]cluster.User, 0)

	for _, u := range instance.Spec.Users {
		user, err := s.kubeClient.Archon().Users(instance.Namespace).Get(u.Name)
		if err != nil {
			return fmt.Errorf("Failed to get user resource %s: %v", u.Name, err), nil
		}
		users = append(users, *user)
	}

	deps := &cluster.InstanceDependency{
		Network: *network,
		Users:   users,
	}

	if instance.Spec.ReservedInstanceRef != nil {
		ri, err := s.kubeClient.Archon().ReservedInstances(instance.Namespace).Get(instance.Spec.ReservedInstanceRef.Name)
		if err != nil {
			return fmt.Errorf("Failed to get reserved instance %s: %v", instance.Spec.ReservedInstanceRef.Name, err), nil
		}
		deps.ReservedInstance = *ri
	}
	return err, deps
}

func (s *InstanceController) ensureSecrets(key string, instance *cluster.Instance) error {
	secrets := make([]v1.Secret, 0)

	for _, n := range instance.Spec.Secrets {
		secret, err := s.kubeClient.Core().Secrets(instance.Namespace).Get(n.Name, metav1.GetOptions{})
		if err != nil {
			return fmt.Errorf("Failed to get secret resource %s: %v", n.Name, err)
		}
		if status, ok := secret.Annotations[certificate.ResourceStatusKey]; ok {
			if status != "Ready" {
				return fmt.Errorf("Secret resource %s is not ready", n.Name)
			}
		}
		secrets = append(secrets, *secret)
	}

	instance.Dependency.Secrets = secrets
	return nil
}

// Do a initial instance sync from cloudprovider. Ignore failures on the cloudpovide side
// (maybe populate defaults too?)
func (s *InstanceController) initializeNewInstance(instance *cluster.Instance) (updatedInstance *cluster.Instance, err error) {
	options := cluster.InstanceOptions{}
	err = util.MapToStruct(instance.Labels, &options, cluster.AnnotationPrefix)
	if err != nil {
		err = fmt.Errorf("Can't get instance options: %s", err.Error())
		return
	}

	if options.UseInstanceID != "" {
		instance.Status.InstanceID = options.UseInstanceID
	}

	status, err := s.archon.GetInstance(s.clusterName, instance)
	if err == nil {
		instance.Status = *status
		err = s.persistUpdate(instance)
		if err != nil {
			err = fmt.Errorf("Unable to persist initial instance status sync: %v", err)
			return
		}
	}

	return instance, nil
}

// Returns an error if processing the instance update failed, along with a time.Duration
// indicating whether processing should be retried; zero means no-retry; otherwise
// we should retry in that Duration.
func (s *InstanceController) processInstanceUpdate(cachedInstance *cachedInstance, instance *cluster.Instance, key string) (err error, retryDelay time.Duration) {
	cachedInstance.mu.Lock()
	defer func() {
		cachedInstance.mu.Unlock()
	}()

	// Instance is being initialized by InstanceGroup controller. so abort now.
	if instance.Status.Phase == cluster.InstanceInitializing {
		return nil, doNotRetry
	}

	// Update instance if it's newly created
	new := cachedInstance.state == nil
	if new && instance.GetDeletionTimestamp() == nil {
		instance, err = s.initializeNewInstance(instance)
		if err != nil {
			return fmt.Errorf("Failed to initialize new instance: %v", err), cachedInstance.nextRetryDelay()
		}
		cachedInstance.state = instance
		return
	}

	// cache the instance, we need the info for instance deletion
	cachedInstance.state = instance

	var (
		retry bool
		obj   metav1.Object
	)
	glog.Infof("Updating instance: %+v", key)
	err, deps := s.ensureDependency(key, instance)
	if err != nil {
		return fmt.Errorf("Failed to ensure all dependencies %s: %v", key, err), cachedInstance.nextRetryDelay()
	}
	instance.Dependency = *deps

	if instance.GetDeletionTimestamp() != nil {
		obj, err, retry = s.initializerManager.FinalizeAll(instance)
		if obj != nil && err == nil {
			instance, _ = obj.(*cluster.Instance)
			cachedInstance.state = instance
		}
	} else if s.initializerManager.NeedsInitialization(instance) {
		obj, err, retry = s.initializerManager.Initialize(instance)
		if obj != nil && err == nil {
			instance, _ = obj.(*cluster.Instance)
			cachedInstance.state = instance
		}
	} else {
		err, retry = s.createInstanceIfNeeded(key, instance)
	}

	if err != nil {
		message := "Error updating the instance"
		if retry {
			message += " (will retry): "
		} else {
			message += " (will not retry): "
		}
		message += err.Error()
		s.eventRecorder.Event(instance, api.EventTypeWarning, "CreatingInstanceFailed", message)

		if retry {
			return err, cachedInstance.nextRetryDelay()
		} else {
			return err, doNotRetry
		}
	}
	// Always update the cache upon success.
	// NOTE: Since we update the cached instance if and only if we successfully
	// processed it, a cached instance being nil implies that it hasn't yet
	// been successfully processed.
	s.cache.set(key, cachedInstance)

	cachedInstance.resetRetryDelay()
	return nil, doNotRetry
}

// Returns whatever error occurred along with a boolean indicator of whether it
// should be retried.
func (s *InstanceController) createInstanceIfNeeded(key string, instance *cluster.Instance) (error, bool) {
	// We will not recreate an instance
	if instance.Status.Phase == cluster.InstanceRunning {
		return nil, notRetryable
	}

	err := s.ensureSecrets(key, instance)
	if err != nil {
		return fmt.Errorf("Failed to ensure all secrets %s: %v", key, err), retryable
	}

	previousState := *cluster.InstanceStatusDeepCopy(&instance.Status)

	glog.V(2).Infof("Ensuring instance %s", key)

	// The instance doesn't exist yet, so create it.
	s.eventRecorder.Event(instance, api.EventTypeNormal, "CreatingInstance", "Creating instance")
	err = s.createInstance(instance)
	if err != nil {
		return fmt.Errorf("Failed to create instance %s: %v", key, err), retryable
	}
	s.eventRecorder.Event(instance, api.EventTypeNormal, "CreatedInstance", "Created instance")

	// Write the status if changed.
	if !cluster.InstanceStatusEqual(previousState, instance.Status) {
		if err := s.persistUpdate(instance); err != nil {
			return fmt.Errorf("Failed to persist updated status to apiserver, even after retries. Giving up: %v", err), notRetryable
		}
	} else {
		glog.V(2).Infof("Not persisting unchanged InstanceStatus to registry.")
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
			glog.Infof("Not persisting update to instance '%s/%s' that no longer exists: %v",
				instance.Namespace, instance.Name, err)
			return nil
		}
		// TODO: Try to resolve the conflict if the change was unrelated to instance
		if errors.IsConflict(err) {
			glog.V(4).Infof("Not persisting update to instance '%s/%s' that has been changed since we received it: %v",
				instance.Namespace, instance.Name, err)
			return nil
		}
		glog.Warningf("Failed to persist updated InstanceStatus to instance '%s/%s' after creating: %v",
			instance.Namespace, instance.Name, err)
		time.Sleep(clientRetryInterval)
	}
	return err
}

func (s *InstanceController) createInstance(instance *cluster.Instance) error {
	status, err := s.archon.EnsureInstance(s.clusterName, instance)
	if err != nil {
		return err
	} else if status != nil {
		instance.Status = *status
	}

	return nil
}

func (s *instanceCache) ListKeys() []string {
	s.mu.Lock()
	defer s.mu.Unlock()
	keys := make([]string, 0, len(s.instanceMap))
	for k := range s.instanceMap {
		keys = append(keys, k)
	}
	return keys
}

func (s *instanceCache) GetByKey(key string) (interface{}, bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if v, ok := s.instanceMap[key]; ok {
		return v, true, nil
	}
	return nil, false, nil
}

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
		instance = &cachedInstance{
			mu: util.NewMutex(),
		}
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
	if newInstance.GetDeletionTimestamp() != nil {
		return true
	}
	if newInstance.Status.Phase == cluster.InstancePending ||
		(oldInstance.Status.Phase == cluster.InstancePending && newInstance.Status.Phase == cluster.InstanceInitializing) {
		return true
	}
	return false
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

// syncInstance will sync the instance with the given key if it has had its expectations fulfilled,
func (s *InstanceController) syncInstance(key string) error {
	startTime := time.Now()
	var cachedInstance *cachedInstance
	var retryDelay time.Duration
	defer func() {
		glog.V(4).Infof("Finished syncing instance %q (%v)", key, time.Now().Sub(startTime))
	}()
	// obj holds the latest instance info from apiserver
	obj, exists, err := s.instanceStore.Indexer.GetByKey(key)
	if err != nil {
		glog.Infof("Unable to retrieve instance %v from store: %v", key, err)
		s.workingQueue.Add(key)
		return err
	}
	if !exists {
		// instance absence in store means watcher caught the deletion, ensure instance is cleaned
		glog.Infof("Instance has been deleted %v", key)
		err, retryDelay = s.processInstanceDeletion(key)
	} else {
		instance, ok := obj.(*cluster.Instance)
		if ok {
			if instance.GetDeletionTimestamp() != nil {
				err, retryDelay = s.processInstanceDeletion(key)
			} else {
				cachedInstance = s.cache.getOrCreate(key)
				err, retryDelay = s.processInstanceUpdate(cachedInstance, instance, key)
			}
		} else {
			tombstone, ok := obj.(cache.DeletedFinalStateUnknown)
			if !ok {
				return fmt.Errorf("object contained wasn't a instance or a deleted key: %#v", obj)
			}
			glog.Infof("Found tombstone for %v", key)
			err, retryDelay = s.processInstanceDeletion(tombstone.Key)
		}
	}

	if retryDelay != 0 {
		// Add the failed instance back to the queue so we'll retry it.
		glog.Errorf("Failed to process instance. Retrying in %s: %v", retryDelay, err)
		go func(obj interface{}, delay time.Duration) {
			// put back the instance key to working queue, it is possible that more entries of the instance
			// were added into the queue during the delay, but it does not mess as when handling the retry,
			// it always get the last instance info from instance store
			s.workingQueue.AddAfter(obj, delay)
		}(key, retryDelay)
	} else if err != nil {
		runtime.HandleError(fmt.Errorf("Failed to process instance. Not retrying: %v", err))
	}
	return nil
}

// Returns an error if processing the instance deletion failed, along with a time.Duration
// indicating whether processing should be retried; zero means no-retry; otherwise
// we should retry after that Duration.
func (s *InstanceController) processInstanceDeletion(key string) (error, time.Duration) {
	cachedInstance, ok := s.cache.get(key)
	if !ok {
		return fmt.Errorf("instance %s not in cache even though the watcher thought it was. Ignoring the deletion.", key), doNotRetry
	}
	cachedInstance.mu.Lock()
	defer func() {
		cachedInstance.mu.Unlock()
	}()
	instance := cachedInstance.state

	obj, exists, err := s.instanceStore.Indexer.GetByKey(key)
	if err != nil {
		glog.Infof("Unable to retrieve instance %v from store: %v", key, err)
		return err, cachedInstance.nextRetryDelay()
	}
	if exists {
		instance = obj.(*cluster.Instance)
	}

	if instance == nil {
		return fmt.Errorf("Unable to process deletion when instance = nil"), doNotRetry
	}

	err, deps := s.ensureDependency(key, instance)
	if err != nil {
		return fmt.Errorf("Failed to ensure all dependencies %s: %v", key, err), cachedInstance.nextRetryDelay()
	}
	instance.Dependency = *deps

	// Instance
	//	glog.Infof("deleting instance %v", instance)
	s.eventRecorder.Event(instance, api.EventTypeNormal, "DeletingInstance", "Deleting instance")
	err = s.archon.EnsureInstanceDeleted(s.clusterName, instance)
	if err != nil {
		message := "Error deleting instance:"
		message += err.Error()
		s.eventRecorder.Event(instance, api.EventTypeWarning, "DeletingInstanceFailed", message)
		return err, cachedInstance.nextRetryDelay()
	}

	obj, err, retry := s.initializerManager.FinalizeAll(instance)
	if err != nil {
		message := "Error finalizing instance:"
		message += err.Error()
		s.eventRecorder.Event(instance, api.EventTypeWarning, "DeletingInstanceFailed", message)
		if retry {
			return err, cachedInstance.nextRetryDelay()
		}
		return err, doNotRetry
	}

	if obj != nil && err == nil {
		instance, _ = obj.(*cluster.Instance)
		cachedInstance.state = instance
	}

	err = s.instanceControl.UnbindInstanceWithReservedInstance(instance)
	if err != nil {
		return err, cachedInstance.nextRetryDelay()
	}

	s.eventRecorder.Event(instance, api.EventTypeNormal, "DeletedInstance", "Deleted instance")
	s.cache.delete(key)

	cachedInstance.resetRetryDelay()
	return nil, doNotRetry
}
