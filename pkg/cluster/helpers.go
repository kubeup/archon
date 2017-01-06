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

package cluster

import (
	"k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/api/unversioned"
	"time"
)

func NetworkStatusDeepCopy(s *NetworkStatus) *NetworkStatus {
	return &NetworkStatus{
		Phase: s.Phase,
	}
}

// IsInstanceAvailable returns true if a instance is available; false otherwise.
// Precondition for an available instance is that it must be ready. On top
// of that, there are two cases when a instance can be considered available:
// 1. minReadySeconds == 0, or
// 2. LastTransitionTime (is set) + minReadySeconds < current time
func IsInstanceAvailable(instance *Instance, minReadySeconds int32, now unversioned.Time) bool {
	if !IsInstanceReady(instance) {
		return false
	}

	c := GetInstanceReadyCondition(instance.Status)
	minReadySecondsDuration := time.Duration(minReadySeconds) * time.Second
	if minReadySeconds == 0 || !c.LastTransitionTime.IsZero() && c.LastTransitionTime.Add(minReadySecondsDuration).Before(now.Time) {
		return true
	}
	return false
}

// IsInstanceReady returns true if a instance is ready; false otherwise.
func IsInstanceReady(instance *Instance) bool {
	return IsInstanceReadyConditionTrue(instance.Status)
}

// IsInstanceReady retruns true if a instance is ready; false otherwise.
func IsInstanceReadyConditionTrue(status InstanceStatus) bool {
	condition := GetInstanceReadyCondition(status)
	return condition != nil && condition.Status == api.ConditionTrue
}

// Extracts the instance ready condition from the given status and returns that.
// Returns nil if the condition is not present.
func GetInstanceReadyCondition(status InstanceStatus) *InstanceCondition {
	_, condition := GetInstanceCondition(&status, InstanceReady)
	return condition
}

// GetInstanceCondition extracts the provided condition from the given status and returns that.
// Returns nil and -1 if the condition is not present, and the index of the located condition.
func GetInstanceCondition(status *InstanceStatus, conditionType InstanceConditionType) (int, *InstanceCondition) {
	if status == nil {
		return -1, nil
	}
	for i := range status.Conditions {
		if status.Conditions[i].Type == conditionType {
			return i, &status.Conditions[i]
		}
	}
	return -1, nil
}

func InstanceStatusEqual(l, r InstanceStatus) bool {
	return l.Phase == r.Phase && l.PrivateIP == r.PrivateIP && l.PublicIP == r.PublicIP && l.InstanceID == r.InstanceID
}

func InstanceStatusDeepCopy(s *InstanceStatus) *InstanceStatus {
	ret := &InstanceStatus{
		Phase:             s.Phase,
		PrivateIP:         s.PrivateIP,
		PublicIP:          s.PublicIP,
		InstanceID:        s.InstanceID,
		CreationTimestamp: s.CreationTimestamp,
	}
	for _, c := range s.Conditions {
		s.Conditions = append(s.Conditions, c)
	}

	return ret
}
