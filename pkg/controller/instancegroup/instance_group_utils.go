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

// If you make changes to this file, you should also make the corresponding change in InstanceGroupController.

package instancegroup

import (
	"fmt"
	"reflect"

	"github.com/golang/glog"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/kubernetes/pkg/api"
	"kubeup.com/archon/pkg/clientset/archon"
	"kubeup.com/archon/pkg/cluster"
)

// updateInstanceGroupStatus attempts to update the Status.Replicas of the given InstanceGroup, with a single GET/PUT retry.
func updateInstanceGroupStatus(c archon.InstanceGroupInterface, ig cluster.InstanceGroup, newStatus cluster.InstanceGroupStatus) (updateErr error) {
	// This is the steady state. It happens when the InstanceGroup doesn't have any expectations, since
	// we do a periodic relist every 30s. If the generations differ but the replicas are
	// the same, a caller might've resized to the same replica count.
	if ig.Status.Replicas == newStatus.Replicas &&
		ig.Status.FullyLabeledReplicas == newStatus.FullyLabeledReplicas &&
		ig.Status.ReadyReplicas == newStatus.ReadyReplicas &&
		ig.Status.AvailableReplicas == newStatus.AvailableReplicas &&
		ig.Generation == ig.Status.ObservedGeneration &&
		reflect.DeepEqual(ig.Status.Conditions, newStatus.Conditions) {
		return nil
	}

	// deep copy to avoid mutation now.
	// TODO this method need some work.  Retry on conflict probably, though I suspect this is stomping status to something it probably shouldn't
	copyObj, err := api.Scheme.DeepCopy(ig)
	if err != nil {
		return err
	}
	ig = copyObj.(cluster.InstanceGroup)

	// Save the generation number we acted on, otherwise we might wrongfully indicate
	// that we've seen a spec update when we retry.
	// TODO: This can clobber an update if we allow multiple agents to write to the
	// same status.
	newStatus.ObservedGeneration = ig.Generation

	var getErr error
	for i, ig := 0, &ig; ; i++ {
		glog.V(4).Infof(fmt.Sprintf("Updating replica count for InstanceGroup: %s/%s, ", ig.Namespace, ig.Name) +
			fmt.Sprintf("replicas %d->%d (need %d), ", ig.Status.Replicas, newStatus.Replicas, ig.Spec.Replicas) +
			fmt.Sprintf("fullyLabeledReplicas %d->%d, ", ig.Status.FullyLabeledReplicas, newStatus.FullyLabeledReplicas) +
			fmt.Sprintf("readyReplicas %d->%d, ", ig.Status.ReadyReplicas, newStatus.ReadyReplicas) +
			fmt.Sprintf("availableReplicas %d->%d, ", ig.Status.AvailableReplicas, newStatus.AvailableReplicas) +
			fmt.Sprintf("sequence No: %v->%v", ig.Status.ObservedGeneration, newStatus.ObservedGeneration))

		ig.Status = newStatus
		_, updateErr = c.UpdateStatus(ig)
		if updateErr == nil || i >= statusUpdateRetries {
			return updateErr
		}
		// Update the InstanceGroup with the latest resource version for the next poll
		if ig, getErr = c.Get(ig.Name); getErr != nil {
			// If the GET fails we can't trust status.Replicas anymore. This error
			// is bound to be more interesting than the update failure.
			return getErr
		}
	}
}

// overlappingInstanceGroups sorts a list of InstanceGroups by creation timestamp, using their names as a tie breaker.
type overlappingInstanceGroups []*cluster.InstanceGroup

func (o overlappingInstanceGroups) Len() int      { return len(o) }
func (o overlappingInstanceGroups) Swap(i, j int) { o[i], o[j] = o[j], o[i] }

func (o overlappingInstanceGroups) Less(i, j int) bool {
	if o[i].CreationTimestamp.Equal(o[j].CreationTimestamp) {
		return o[i].Name < o[j].Name
	}
	return o[i].CreationTimestamp.Before(o[j].CreationTimestamp)
}

func calculateStatus(ig cluster.InstanceGroup, filteredInstances []*cluster.Instance, manageReplicasErr error) cluster.InstanceGroupStatus {
	newStatus := ig.Status
	// Count the number of instances that have labels matching the labels of the instance
	// template of the replica set, the matching instances may have more
	// labels than are in the template. Because the label of instanceTemplateSpec is
	// a superset of the selector of the replica set, so the possible
	// matching instances must be part of the filteredInstances.
	fullyLabeledReplicasCount := 0
	readyReplicasCount := 0
	availableReplicasCount := 0
	templateLabel := labels.Set(ig.Spec.Template.Labels).AsSelectorPreValidated()
	for _, instance := range filteredInstances {
		if templateLabel.Matches(labels.Set(instance.Labels)) {
			fullyLabeledReplicasCount++
		}
		if cluster.IsInstanceReady(instance) {
			readyReplicasCount++
			if cluster.IsInstanceAvailable(instance, ig.Spec.MinReadySeconds, metav1.Now()) {
				availableReplicasCount++
			}
		}
	}

	failureCond := GetCondition(ig.Status, cluster.InstanceGroupReplicaFailure)
	if manageReplicasErr != nil && failureCond == nil {
		var reason string
		if diff := len(filteredInstances) - int(ig.Spec.Replicas); diff < 0 {
			reason = "FailedCreate"
		} else if diff > 0 {
			reason = "FailedDelete"
		}
		cond := NewInstanceGroupCondition(cluster.InstanceGroupReplicaFailure, api.ConditionTrue, reason, manageReplicasErr.Error())
		SetCondition(&newStatus, cond)
	} else if manageReplicasErr == nil && failureCond != nil {
		RemoveCondition(&newStatus, cluster.InstanceGroupReplicaFailure)
	}

	newStatus.Replicas = int32(len(filteredInstances))
	newStatus.FullyLabeledReplicas = int32(fullyLabeledReplicasCount)
	newStatus.ReadyReplicas = int32(readyReplicasCount)
	newStatus.AvailableReplicas = int32(availableReplicasCount)
	return newStatus
}

// NewInstanceGroupCondition creates a new replica set condition.
func NewInstanceGroupCondition(condType cluster.InstanceGroupConditionType, status api.ConditionStatus, reason, msg string) cluster.InstanceGroupCondition {
	return cluster.InstanceGroupCondition{
		Type:               condType,
		Status:             status,
		LastTransitionTime: metav1.Now(),
		Reason:             reason,
		Message:            msg,
	}
}

// GetCondition returns a replica set condition with the provided type if it exists.
func GetCondition(status cluster.InstanceGroupStatus, condType cluster.InstanceGroupConditionType) *cluster.InstanceGroupCondition {
	for i := range status.Conditions {
		c := status.Conditions[i]
		if c.Type == condType {
			return &c
		}
	}
	return nil
}

// SetCondition adds/replaces the given condition in the replica set status. If the condition that we
// are about to add already exists and has the same status and reason then we are not going to update.
func SetCondition(status *cluster.InstanceGroupStatus, condition cluster.InstanceGroupCondition) {
	currentCond := GetCondition(*status, condition.Type)
	if currentCond != nil && currentCond.Status == condition.Status && currentCond.Reason == condition.Reason {
		return
	}
	newConditions := filterOutCondition(status.Conditions, condition.Type)
	status.Conditions = append(newConditions, condition)
}

// RemoveCondition removes the condition with the provided type from the replica set status.
func RemoveCondition(status *cluster.InstanceGroupStatus, condType cluster.InstanceGroupConditionType) {
	status.Conditions = filterOutCondition(status.Conditions, condType)
}

// filterOutCondition returns a new slice of replica set conditions without conditions with the provided type.
func filterOutCondition(conditions []cluster.InstanceGroupCondition, condType cluster.InstanceGroupConditionType) []cluster.InstanceGroupCondition {
	var newConditions []cluster.InstanceGroupCondition
	for _, c := range conditions {
		if c.Type == condType {
			continue
		}
		newConditions = append(newConditions, c)
	}
	return newConditions
}
