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

// This file is modified from replica_set.go in the original kubernetes source tree

package instancegroup

import (
	"fmt"
	"reflect"
	"sort"
	"sync"
	"time"

	"k8s.io/apimachinery/pkg/labels"
	v1core "k8s.io/client-go/kubernetes/typed/core/v1"
	archoncache "kubeup.com/archon/pkg/cache"
	"kubeup.com/archon/pkg/clientset"
	"kubeup.com/archon/pkg/cluster"
	archoncontroller "kubeup.com/archon/pkg/controller"

	"github.com/golang/glog"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	pkg_runtime "k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	utilerrors "k8s.io/apimachinery/pkg/util/errors"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	watch "k8s.io/apimachinery/pkg/watch"
	clientv1 "k8s.io/client-go/pkg/api/v1"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/record"
	"k8s.io/client-go/util/workqueue"
	"k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/controller"
	"k8s.io/kubernetes/pkg/util/metrics"
)

const (
	// We'll attempt to recompute the required replicas of all InstanceGroups
	// that have fulfilled their expectations at least this often. This recomputation
	// happens based on contents in local instance storage.
	FullControllerResyncPeriod = 30 * time.Second

	// Realistic value of the burstReplica field for the replica set manager based off
	// performance requirements for kubernetes 1.0.
	// TODO: adjust for instancegroups
	BurstReplicas = 500

	// The number of times we retry updating a InstanceGroup's status.
	statusUpdateRetries = 1
)

func getIGKind() schema.GroupVersionKind {
	return cluster.SchemeGroupVersion.WithKind("InstanceGroup")
}

// InstanceGroupController is responsible for synchronizing InstanceGroup objects stored
// in the system with actual running Instances.
type InstanceGroupController struct {
	kubeClient      clientset.Interface
	instanceControl archoncontroller.InstanceControlInterface
	namespace       string

	// A InstanceGroup is temporarily suspended after creating/deleting these many replicas.
	// It resumes normal action after observing the watch events for them.
	burstReplicas int
	// To allow injection of syncInstanceGroup for testing.
	syncHandler func(igKey string) error

	// A TTLCache of instance creates/deletes each rc expects to see.
	expectations *controller.UIDTrackingControllerExpectations

	instanceStore           archoncache.StoreToInstanceLister
	instanceGroupStore      archoncache.StoreToInstanceGroupLister
	instanceController      cache.Controller
	instanceGroupController cache.Controller

	lookupCache *controller.MatchingCache

	// Controllers that need to be synced
	queue workqueue.RateLimitingInterface

	// garbageCollectorEnabled denotes if the garbage collector is enabled. RC
	// manager behaves differently if GC is enabled.
	garbageCollectorEnabled bool
}

// NewInstanceGroupController configures a replica set controller with the specified event recorder
func NewInstanceGroupController(kubeClient clientset.Interface, namespace string, burstReplicas int, lookupCacheSize int, garbageCollectorEnabled bool) *InstanceGroupController {
	if kubeClient != nil && kubeClient.Core().RESTClient().GetRateLimiter() != nil {
		metrics.RegisterMetricAndTrackRateLimiterUsage("instance_group_controller", kubeClient.Core().RESTClient().GetRateLimiter())
	}
	eventBroadcaster := record.NewBroadcaster()
	eventBroadcaster.StartLogging(glog.Infof)
	eventBroadcaster.StartRecordingToSink(&v1core.EventSinkImpl{Interface: v1core.New(kubeClient.Core().RESTClient()).Events("")})

	rsc := &InstanceGroupController{
		kubeClient: kubeClient,
		instanceControl: archoncontroller.RealInstanceControl{
			KubeClient: kubeClient,
			Recorder:   eventBroadcaster.NewRecorder(api.Scheme, clientv1.EventSource{Component: "instance-group-controller"}),
		},
		namespace:     namespace,
		burstReplicas: burstReplicas,
		expectations:  controller.NewUIDTrackingControllerExpectations(controller.NewControllerExpectations()),
		queue:         workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "instancegroup"),
		garbageCollectorEnabled: garbageCollectorEnabled,
	}

	rsc.instanceStore.Indexer, rsc.instanceController = cache.NewIndexerInformer(
		&cache.ListWatch{
			ListFunc: func(options metav1.ListOptions) (pkg_runtime.Object, error) {
				options.IncludeUninitialized = true
				return rsc.kubeClient.Archon().Instances(namespace).List(options)
			},
			WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
				options.IncludeUninitialized = true
				return rsc.kubeClient.Archon().Instances(namespace).Watch(options)
			},
		},
		&cluster.Instance{},
		FullControllerResyncPeriod,
		cache.ResourceEventHandlerFuncs{
			AddFunc:    rsc.addInstance,
			UpdateFunc: rsc.updateInstance,
			DeleteFunc: rsc.deleteInstance,
		},
		cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc},
	)

	rsc.instanceGroupStore.Indexer, rsc.instanceGroupController = cache.NewIndexerInformer(
		&cache.ListWatch{
			ListFunc: func(options metav1.ListOptions) (pkg_runtime.Object, error) {
				return rsc.kubeClient.Archon().InstanceGroups(namespace).List(options)
			},
			WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
				return rsc.kubeClient.Archon().InstanceGroups(namespace).Watch(options)
			},
		},
		&cluster.InstanceGroup{},
		FullControllerResyncPeriod,
		cache.ResourceEventHandlerFuncs{
			AddFunc:    rsc.enqueueInstanceGroup,
			UpdateFunc: rsc.updateIG,
			DeleteFunc: rsc.enqueueInstanceGroup,
		},
		cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc},
	)

	rsc.syncHandler = rsc.syncInstanceGroup
	rsc.lookupCache = controller.NewMatchingCache(lookupCacheSize)
	return rsc
}

// SetEventRecorder replaces the event recorder used by the InstanceGroupController
// with the given recorder. Only used for testing.
func (rsc *InstanceGroupController) SetEventRecorder(recorder record.EventRecorder) {
	// TODO: Hack. We can't cleanly shutdown the event recorder, so benchmarks
	// need to pass in a fake.
	rsc.instanceControl = archoncontroller.RealInstanceControl{KubeClient: rsc.kubeClient, Recorder: recorder}
}

// Run begins watching and syncing.
func (rsc *InstanceGroupController) Run(workers int, stopCh <-chan struct{}) {
	defer utilruntime.HandleCrash()
	defer rsc.queue.ShutDown()

	glog.Infof("Starting InstanceGroup controller")

	go rsc.instanceController.Run(wait.NeverStop)
	go rsc.instanceGroupController.Run(wait.NeverStop)

	if !cache.WaitForCacheSync(stopCh, rsc.instanceController.HasSynced) {
		return
	}

	for i := 0; i < workers; i++ {
		go wait.Until(rsc.worker, time.Second, stopCh)
	}

	<-stopCh
	glog.Infof("Shutting down InstanceGroup Controller")
}

// getInstanceGroup returns the replica set managing the given instance.
func (rsc *InstanceGroupController) getInstanceGroup(instance *cluster.Instance) *cluster.InstanceGroup {
	// look up in the cache, if cached and the cache is valid, just return cached value
	if obj, cached := rsc.lookupCache.GetMatchingObject(instance); cached {
		ig, ok := obj.(*cluster.InstanceGroup)
		if !ok {
			// This should not happen
			utilruntime.HandleError(fmt.Errorf("lookup cache does not return a InstanceGroup object"))
			return nil
		}
		if cached && rsc.isCacheValid(instance, ig) {
			return ig
		}
	}

	// if not cached or cached value is invalid, search all the rs to find the matching one, and update cache
	igs, err := rsc.instanceGroupStore.GetInstanceGroup(instance)
	if err != nil {
		glog.V(4).Infof("No InstanceGroups found for instance %v, InstanceGroup controller will avoid syncing", instance.Name)
		return nil
	}
	// In theory, overlapping InstanceGroup's user error. This sorting will not prevent
	// oscillation of replicas in all cases, eg:
	// ig1 (older ig): [(k1=v1)], replicas=1 ig2: [(k2=v2)], replicas=2
	// instance: [(k1:v1), (k2:v2)] will wake both ig1 and ig2, and we will sync ig1.
	// instance: [(k2:v2)] will wake ig2 which creates a new replica.
	if len(igs) > 1 {
		// More than two items in this list indicates user error. If two instancegroups
		// overlap, sort by creation timestamp, subsort by name, then pick
		// the first.
		utilruntime.HandleError(fmt.Errorf("user error! more than one InstanceGroup is selecting instances with labels: %+v", instance.Labels))
		sort.Sort(overlappingInstanceGroups(igs))
	}

	// update lookup cache
	rsc.lookupCache.Update(instance, igs[0])

	return igs[0]
}

// callback when InstanceGroup is updated
func (rsc *InstanceGroupController) updateIG(old, cur interface{}) {
	oldIG := old.(*cluster.InstanceGroup)
	curIG := cur.(*cluster.InstanceGroup)

	// We should invalidate the whole lookup cache if a IG's selector has been updated.
	//
	// Imagine that you have two IGs:
	// * old IG1
	// * new IG2
	// You also have a instance that is attached to IG2 (because it doesn't match IG1 selector).
	// Now imagine that you are changing IG1 selector so that it is now matching that instance,
	// in such case we must invalidate the whole cache so that instance could be adopted by IG1
	//
	// This makes the lookup cache less helpful, but selector update does not happen often,
	// so it's not a big problem
	if !reflect.DeepEqual(oldIG.Spec.Selector, curIG.Spec.Selector) {
		rsc.lookupCache.InvalidateAll()
	}

	// You might imagine that we only really need to enqueue the
	// replica set when Spec changes, but it is safer to sync any
	// time this function is triggered. That way a full informer
	// resync can requeue any replica set that don't yet have instances
	// but whose last attempts at creating a instance have failed (since
	// we don't block on creation of instances) instead of those
	// replica sets stalling indefinitely. Enqueueing every time
	// does result in some spurious syncs (like when Status.Replica
	// is updated and the watch notification from it retriggers
	// this function), but in general extra resyncs shouldn't be
	// that bad as InstanceGroups that haven't met expectations yet won't
	// sync, and all the listing is done using local stores.
	if oldIG.Status.Replicas != curIG.Status.Replicas {
		glog.V(4).Infof("Observed updated replica count for InstanceGroup: %v, %d->%d", curIG.Name, oldIG.Status.Replicas, curIG.Status.Replicas)
	}
	rsc.enqueueInstanceGroup(cur)
}

// isCacheValid check if the cache is valid
func (rsc *InstanceGroupController) isCacheValid(instance *cluster.Instance, cachedIG *cluster.InstanceGroup) bool {
	_, err := rsc.instanceGroupStore.InstanceGroups(cachedIG.Namespace).Get(cachedIG.Name)
	// rs has been deleted or updated, cache is invalid
	if err != nil || !isInstanceGroupMatch(instance, cachedIG) {
		return false
	}
	return true
}

// isInstanceGroupMatch take a Instance and InstanceGroup, return whether the Instance and InstanceGroup are matching
// TODO(mqliang): This logic is a copy from GetInstanceGroups(), remove the duplication
func isInstanceGroupMatch(instance *cluster.Instance, ig *cluster.InstanceGroup) bool {
	if ig.Namespace != instance.Namespace {
		return false
	}
	selector, err := metav1.LabelSelectorAsSelector(ig.Spec.Selector)
	if err != nil {
		err = fmt.Errorf("invalid selector: %v", err)
		return false
	}

	// If a ReplicaSet with a nil or empty selector creeps in, it should match nothing, not everything.
	if selector.Empty() || !selector.Matches(labels.Set(instance.Labels)) {
		return false
	}
	return true
}

// When a instance is created, enqueue the replica set that manages it and update it's expectations.
func (rsc *InstanceGroupController) addInstance(obj interface{}) {
	instance := obj.(*cluster.Instance)
	glog.V(4).Infof("Instance %s created: %#v.", instance.Name, instance)

	ig := rsc.getInstanceGroup(instance)
	if ig == nil {
		return
	}
	igKey, err := controller.KeyFunc(ig)
	if err != nil {
		utilruntime.HandleError(fmt.Errorf("Couldn't get key for instance group %#v: %v", ig, err))
		return
	}
	if instance.DeletionTimestamp != nil {
		// on a restart of the controller manager, it's possible a new instance shows up in a state that
		// is already pending deletion. Prevent the instance from being a creation observation.
		rsc.deleteInstance(instance)
		return
	}
	rsc.expectations.CreationObserved(igKey)
	rsc.enqueueInstanceGroup(ig)
}

// When a instance is updated, figure out what replica set/s manage it and wake them
// up. If the labels of the instance have changed we need to awaken both the old
// and new replica set. old and cur must be *cluster.Instance types.
func (rsc *InstanceGroupController) updateInstance(old, cur interface{}) {
	curInstance := cur.(*cluster.Instance)
	oldInstance := old.(*cluster.Instance)
	if curInstance.ResourceVersion == oldInstance.ResourceVersion {
		// Periodic resync will send update events for all known instances.
		// Two different versions of the same instance will always have different RVs.
		return
	}
	glog.V(4).Infof("Instance %s updated, objectMeta %+v -> %+v.", curInstance.Name, oldInstance.ObjectMeta, curInstance.ObjectMeta)
	labelChanged := !reflect.DeepEqual(curInstance.Labels, oldInstance.Labels)
	if curInstance.DeletionTimestamp != nil {
		// when a instance is deleted gracefully it's deletion timestamp is first modified to reflect a grace period,
		// and after such time has passed, the kubelet actually deletes it from the store. We receive an update
		// for modification of the deletion timestamp and expect an rs to create more replicas asap, not wait
		// until the kubelet actually deletes the instance. This is different from the Phase of a instance changing, because
		// an rs never initiates a phase change, and so is never asleep waiting for the same.
		rsc.deleteInstance(curInstance)
		if labelChanged {
			// we don't need to check the oldInstance.DeletionTimestamp because DeletionTimestamp cannot be unset.
			rsc.deleteInstance(oldInstance)
		}
		return
	}

	// Enqueue the oldRC before the cur to give curRC a chance to adopt the oldInstance.
	if labelChanged {
		// If the old and new InstanceGroup are the same, the first one that syncs
		// will set expectations preventing any damage from the second.
		if oldIG := rsc.getInstanceGroup(oldInstance); oldIG != nil {
			rsc.enqueueInstanceGroup(oldIG)
		}
	}

	if curIG := rsc.getInstanceGroup(curInstance); curIG != nil {
		rsc.enqueueInstanceGroup(curIG)
	}
}

// When a instance is deleted, enqueue the replica set that manages the instance and update its expectations.
// obj could be an *cluster.Instance, or a DeletionFinalStateUnknown marker item.
func (rsc *InstanceGroupController) deleteInstance(obj interface{}) {
	instance, ok := obj.(*cluster.Instance)

	// When a delete is dropped, the relist will notice a instance in the store not
	// in the list, leading to the insertion of a tombstone object which contains
	// the deleted key/value. Note that this value might be stale. If the instance
	// changed labels the new InstanceGroup will not be woken up till the periodic resync.
	if !ok {
		tombstone, ok := obj.(cache.DeletedFinalStateUnknown)
		if !ok {
			utilruntime.HandleError(fmt.Errorf("Couldn't get object from tombstone %+v", obj))
			return
		}
		instance, ok = tombstone.Obj.(*cluster.Instance)
		if !ok {
			utilruntime.HandleError(fmt.Errorf("Tombstone contained object that is not a instance %#v", obj))
			return
		}
	}
	glog.V(4).Infof("Instance %s/%s deleted through %v, timestamp %+v: %#v.", instance.Namespace, instance.Name, utilruntime.GetCaller(), instance.DeletionTimestamp, instance)
	if ig := rsc.getInstanceGroup(instance); ig != nil {
		igKey, err := controller.KeyFunc(ig)
		if err != nil {
			utilruntime.HandleError(fmt.Errorf("Couldn't get key for InstanceGroup %#v: %v", ig, err))
			return
		}
		rsc.expectations.DeletionObserved(igKey, archoncontroller.InstanceKey(instance))
		rsc.enqueueInstanceGroup(ig)
	}
}

// obj could be an *extensions.InstanceGroup, or a DeletionFinalStateUnknown marker item.
func (rsc *InstanceGroupController) enqueueInstanceGroup(obj interface{}) {
	key, err := controller.KeyFunc(obj)
	if err != nil {
		utilruntime.HandleError(fmt.Errorf("Couldn't get key for object %+v: %v", obj, err))
		return
	}

	// TODO: Handle overlapping replica sets better. Either disallow them at admission time or
	// deterministically avoid syncing replica sets that fight over instances. Currently, we only
	// ensure that the same replica set is synced for a given instance. When we periodically relist
	// all replica sets there will still be some replica instability. One way to handle this is
	// by querying the store for all replica sets that this replica set overlaps, as well as all
	// replica sets that overlap this ReplicaSet, and sorting them.
	rsc.queue.Add(key)
}

// worker runs a worker thread that just dequeues items, processes them, and marks them done.
// It enforces that the syncHandler is never invoked concurrently with the same key.
func (rsc *InstanceGroupController) worker() {
	for rsc.processNextWorkItem() {
	}
}

func (rsc *InstanceGroupController) processNextWorkItem() bool {
	key, quit := rsc.queue.Get()
	if quit {
		return false
	}
	defer rsc.queue.Done(key)

	err := rsc.syncHandler(key.(string))
	if err == nil {
		rsc.queue.Forget(key)
		return true
	}

	utilruntime.HandleError(fmt.Errorf("Sync %q failed with %v", key, err))
	// TODO: Delay ig sync when running out of reserved instances
	rsc.queue.AddRateLimited(key)

	return true
}

// manageReplicas checks and updates replicas for the given ReplicaSet.
// Does NOT modify <filteredinstances>.
// It will requeue the replica set in case of an error while creating/deleting instances.
func (rsc *InstanceGroupController) manageReplicas(filteredInstances []*cluster.Instance, ig *cluster.InstanceGroup) error {
	diff := len(filteredInstances) - int(ig.Spec.Replicas)
	igKey, err := controller.KeyFunc(ig)
	if err != nil {
		utilruntime.HandleError(fmt.Errorf("Couldn't get key for InstanceGroup %#v: %v", ig, err))
		return nil
	}
	var errCh chan error
	if diff < 0 {
		diff *= -1
		errCh = make(chan error, diff)
		if diff > rsc.burstReplicas {
			diff = rsc.burstReplicas
		}
		// TODO: Track UIDs of creates just like deletes. The problem currently
		// is we'd need to wait on the result of a create to record the instance's
		// UID, which would require locking *across* the create, which will turn
		// into a performance bottleneck. We should generate a UID for the instance
		// beforehand and store it via ExpectCreations.
		rsc.expectations.ExpectCreations(igKey, diff)
		var wg sync.WaitGroup
		wg.Add(diff)
		glog.V(2).Infof("Too few %q/%q replicas, need %d, creating %d", ig.Namespace, ig.Name, ig.Spec.Replicas, diff)
		for i := 0; i < diff; i++ {
			go func() {
				defer wg.Done()
				var err error

				if rsc.garbageCollectorEnabled {
					var trueVar = true
					controllerRef := &metav1.OwnerReference{
						APIVersion: getIGKind().GroupVersion().String(),
						Kind:       getIGKind().Kind,
						Name:       ig.Name,
						UID:        ig.UID,
						Controller: &trueVar,
					}
					err = rsc.instanceControl.CreateInstancesWithControllerRef(ig.Namespace, &ig.Spec.Template, ig, controllerRef)
				} else {
					err = rsc.instanceControl.CreateInstances(ig.Namespace, &ig.Spec.Template, ig)
				}
				if err != nil {
					// Decrement the expected number of creates because the informer won't observe this instance
					glog.V(2).Infof("Failed creation, decrementing expectations for replica set %q/%q", ig.Namespace, ig.Name)
					rsc.expectations.CreationObserved(igKey)
					errCh <- err
				}
			}()
		}
		wg.Wait()
	} else if diff > 0 {
		if diff > rsc.burstReplicas {
			diff = rsc.burstReplicas
		}
		errCh = make(chan error, diff)
		glog.V(2).Infof("Too many %q/%q replicas, need %d, deleting %d", ig.Namespace, ig.Name, ig.Spec.Replicas, diff)
		// No need to sort instances if we are about to delete all of them
		if ig.Spec.Replicas != 0 {
			// Sort the instances in the order such that not-ready < ready, unscheduled
			// < scheduled, and pending < running. This ensures that we delete instances
			// in the earlier stages whenever possible.
			sort.Sort(archoncontroller.ActiveInstances(filteredInstances))
		}
		// Snapshot the UIDs (ns/name) of the instances we're expecting to see
		// deleted, so we know to record their expectations exactly once either
		// when we see it as an update of the deletion timestamp, or as a delete.
		// Note that if the labels on a instance/rs change in a way that the instance gets
		// orphaned, the rs will only wake up after the expectations have
		// expired even if other instances are deleted.
		deletedInstanceKeys := []string{}
		for i := 0; i < diff; i++ {
			deletedInstanceKeys = append(deletedInstanceKeys, archoncontroller.InstanceKey(filteredInstances[i]))
		}
		rsc.expectations.ExpectDeletions(igKey, deletedInstanceKeys)
		var wg sync.WaitGroup
		wg.Add(diff)
		for i := 0; i < diff; i++ {
			go func(ix int) {
				defer wg.Done()
				if err := rsc.instanceControl.DeleteInstance(ig.Namespace, filteredInstances[ix].Name, ig); err != nil {
					// Decrement the expected number of deletes because the informer won't observe this deletion
					instanceKey := archoncontroller.InstanceKey(filteredInstances[ix])
					glog.V(2).Infof("Failed to delete %v, decrementing expectations for controller %q/%q", instanceKey, ig.Namespace, ig.Name)
					rsc.expectations.DeletionObserved(igKey, instanceKey)
					errCh <- err
				}
			}(i)
		}
		wg.Wait()
	}

	select {
	case err := <-errCh:
		// all errors have been reported before and they're likely to be the same, so we'll only return the first one we hit.
		if err != nil {
			return err
		}
	default:
	}
	return nil
}

// syncReplicaSet will sync the ReplicaSet with the given key if it has had its expectations fulfilled,
// meaning it did not expect to see any more of its instances created or deleted. This function is not meant to be
// invoked concurrently with the same key.
func (rsc *InstanceGroupController) syncInstanceGroup(key string) error {
	startTime := time.Now()
	defer func() {
		glog.V(4).Infof("Finished syncing replica set %q (%v)", key, time.Now().Sub(startTime))
	}()

	obj, exists, err := rsc.instanceGroupStore.Indexer.GetByKey(key)
	if !exists {
		glog.V(4).Infof("InstanceGroup has been deleted %v", key)
		rsc.expectations.DeleteExpectations(key)
		return nil
	}
	if err != nil {
		return err
	}
	ig := *obj.(*cluster.InstanceGroup)

	// Check the expectations of the ReplicaSet before counting active instances, otherwise a new instance can sneak
	// in and update the expectations after we've retrieved active instances from the store. If a new instance enters
	// the store after we've checked the expectation, the ReplicaSet sync is just deferred till the next
	// relist.
	if err != nil {
		utilruntime.HandleError(fmt.Errorf("Couldn't get key for InstanceGroup %#v: %v", ig, err))
		// Explicitly return nil to avoid re-enqueue bad key
		return nil
	}
	igNeedsSync := rsc.expectations.SatisfiedExpectations(key)
	selector, err := metav1.LabelSelectorAsSelector(ig.Spec.Selector)
	if err != nil {
		utilruntime.HandleError(fmt.Errorf("Error converting instance selector to selector: %v", err))
		return nil
	}

	// NOTE: filteredinstances are pointing to objects from cache - if you need to
	// modify them, you need to copy it first.
	// TODO: Do the List and Filter in a single pass, or use an index.
	var filteredInstances []*cluster.Instance
	if rsc.garbageCollectorEnabled {
		// list all instances to include the instances that don't match the rs`s selector
		// anymore but has the stale controller ref.
		instances, err := rsc.instanceStore.Instances(ig.Namespace).List(labels.Everything())
		if err != nil {
			return err
		}
		cm := archoncontroller.NewInstanceControllerRefManager(rsc.instanceControl, ig.ObjectMeta, selector, getIGKind())
		matchesAndControlled, matchesNeedsController, controlledDoesNotMatch := cm.Classify(instances)
		for _, instance := range matchesNeedsController {
			err := cm.AdoptInstance(instance)
			// continue to next instance if adoption fails.
			if err != nil {
				// If the instance no longer exists, don't even log the error.
				if !errors.IsNotFound(err) {
					utilruntime.HandleError(err)
				}
			} else {
				matchesAndControlled = append(matchesAndControlled, instance)
			}
		}
		filteredInstances = matchesAndControlled
		// remove the controllerRef for the instances that no longer have matching labels
		var errlist []error
		for _, instance := range controlledDoesNotMatch {
			err := cm.ReleaseInstance(instance)
			if err != nil {
				errlist = append(errlist, err)
			}
		}
		if len(errlist) != 0 {
			aggregate := utilerrors.NewAggregate(errlist)
			// push the IG into work queue again. We need to try to free the
			// instances again otherwise they will stuck with the stale
			// controllerRef.
			return aggregate
		}
	} else {
		instances, err := rsc.instanceStore.Instances(ig.Namespace).List(selector)
		if err != nil {
			return err
		}
		filteredInstances = archoncontroller.FilterActiveInstances(instances)
	}

	var manageReplicasErr error
	if igNeedsSync && ig.DeletionTimestamp == nil {
		manageReplicasErr = rsc.manageReplicas(filteredInstances, &ig)
	}

	newStatus := calculateStatus(ig, filteredInstances, manageReplicasErr)

	// Always updates status as instances come up or die.
	if err := updateInstanceGroupStatus(rsc.kubeClient.Archon().InstanceGroups(ig.Namespace), ig, newStatus); err != nil {
		// Multiple things could lead to this update failing. Requeuing the replica set ensures
		// Returning an error causes a requeue without forcing a hotloop
		return err
	}
	return manageReplicasErr
}
