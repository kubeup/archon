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

package controller

import (
	"fmt"
	"strings"

	"kubeup.com/archon/pkg/cluster"

	"github.com/golang/glog"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

type InstanceControllerRefManager struct {
	instanceControl    InstanceControlInterface
	controllerObject   metav1.ObjectMeta
	controllerSelector labels.Selector
	controllerKind     schema.GroupVersionKind
}

// NewInstanceControllerRefManager returns a InstanceControllerRefManager that exposes
// methods to manage the controllerRef of instances.
func NewInstanceControllerRefManager(
	instanceControl InstanceControlInterface,
	controllerObject metav1.ObjectMeta,
	controllerSelector labels.Selector,
	controllerKind schema.GroupVersionKind,
) *InstanceControllerRefManager {
	return &InstanceControllerRefManager{instanceControl, controllerObject, controllerSelector, controllerKind}
}

// Classify first filters out inactive instances, then it classify the remaining instances
// into three categories: 1. matchesAndControlled are the instances whose labels
// match the selector of the RC, and have a controllerRef pointing to the
// controller 2. matchesNeedsController are the instances whose labels match the RC,
// but don't have a controllerRef. (Instances with matching labels but with a
// controllerRef pointing to other object are ignored) 3. controlledDoesNotMatch
// are the instances that have a controllerRef pointing to the controller, but their
// labels no longer match the selector.
func (m *InstanceControllerRefManager) Classify(instances []*cluster.Instance) (
	matchesAndControlled []*cluster.Instance,
	matchesNeedsController []*cluster.Instance,
	controlledDoesNotMatch []*cluster.Instance) {
	for i := range instances {
		instance := instances[i]
		if !IsInstanceActive(instance) {
			glog.V(4).Infof("Ignoring inactive instance %v/%v in state %v, deletion time %v",
				instance.Namespace, instance.Name, instance.Status.Phase, instance.DeletionTimestamp)
			continue
		}
		controllerRef := getControllerOf(instance.ObjectMeta)
		if controllerRef != nil {
			if controllerRef.UID == m.controllerObject.UID {
				// already controlled
				if m.controllerSelector.Matches(labels.Set(instance.Labels)) {
					matchesAndControlled = append(matchesAndControlled, instance)
				} else {
					controlledDoesNotMatch = append(controlledDoesNotMatch, instance)
				}
			} else {
				// ignoring the instance controlled by other controller
				glog.V(4).Infof("Ignoring instance %v/%v, it's owned by [%s/%s, name: %s, uid: %s]",
					instance.Namespace, instance.Name, controllerRef.APIVersion, controllerRef.Kind, controllerRef.Name, controllerRef.UID)
				continue
			}
		} else {
			if !m.controllerSelector.Matches(labels.Set(instance.Labels)) {
				continue
			}
			matchesNeedsController = append(matchesNeedsController, instance)
		}
	}
	return matchesAndControlled, matchesNeedsController, controlledDoesNotMatch
}

// getControllerOf returns the controllerRef if controllee has a controller,
// otherwise returns nil.
func getControllerOf(controllee metav1.ObjectMeta) *metav1.OwnerReference {
	for _, owner := range controllee.OwnerReferences {
		// controlled by other controller
		if owner.Controller != nil && *owner.Controller == true {
			return &owner
		}
	}
	return nil
}

// AdoptInstance sends a patch to take control of the instance. It returns the error if
// the patching fails.
func (m *InstanceControllerRefManager) AdoptInstance(instance *cluster.Instance) error {
	// we should not adopt any instances if the controller is about to be deleted
	if m.controllerObject.DeletionTimestamp != nil {
		return fmt.Errorf("cancel the adopt attempt for instance %s because the controlller is being deleted",
			strings.Join([]string{instance.Namespace, instance.Name, string(instance.UID)}, "_"))
	}
	addControllerPatch := fmt.Sprintf(
		`{"metadata":{"ownerReferences":[{"apiVersion":"%s","kind":"%s","name":"%s","uid":"%s","controller":true}],"uid":"%s"}}`,
		m.controllerKind.GroupVersion(), m.controllerKind.Kind,
		m.controllerObject.Name, m.controllerObject.UID, instance.UID)
	return m.instanceControl.PatchInstance(instance.Namespace, instance.Name, []byte(addControllerPatch))
}

// ReleaseInstance sends a patch to free the instance from the control of the controller.
// It returns the error if the patching fails. 404 and 422 errors are ignored.
func (m *InstanceControllerRefManager) ReleaseInstance(instance *cluster.Instance) error {
	glog.V(2).Infof("patching instance %s_%s to remove its controllerRef to %s/%s:%s",
		instance.Namespace, instance.Name, m.controllerKind.GroupVersion(), m.controllerKind.Kind, m.controllerObject.Name)
	deleteOwnerRefPatch := fmt.Sprintf(`{"metadata":{"ownerReferences":[{"$patch":"delete","uid":"%s"}],"uid":"%s"}}`, m.controllerObject.UID, instance.UID)
	err := m.instanceControl.PatchInstance(instance.Namespace, instance.Name, []byte(deleteOwnerRefPatch))
	if err != nil {
		if errors.IsNotFound(err) {
			// If the instance no longer exists, ignore it.
			return nil
		}
		if errors.IsInvalid(err) {
			// Invalid error will be returned in two cases: 1. the instance
			// has no owner reference, 2. the uid of the instance doesn't
			// match, which means the instance is deleted and then recreated.
			// In both cases, the error can be ignored.

			// TODO: If the instance has owner references, but none of them
			// has the owner.UID, server will silently ignore the patch.
			// Investigate why.
			return nil
		}
	}
	return err
}
