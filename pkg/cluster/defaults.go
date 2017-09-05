/*
Copyright 2017 The Archon Authors.
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
	"k8s.io/apimachinery/pkg/runtime"
)

func addDefaultingFuncs(scheme *runtime.Scheme) error {
	return RegisterDefaults(scheme)
}

func SetDefaults_InstanceGroup(obj *InstanceGroup) {
	if obj.Spec.ProvisionPolicy == "" {
		obj.Spec.ProvisionPolicy = InstanceGroupProvisionDynamicOnly
	}
}

func SetDefaults_Instance(obj *Instance) {
	if obj.Spec.ReclaimPolicy == "" {
		obj.Spec.ReclaimPolicy = InstanceReclaimDelete
	}

	if obj.Status.Phase == "" {
		obj.Status.Phase = InstancePending
	}
}

func SetDefaults_Network(obj *Network) {
	if obj.Status.Phase == "" {
		obj.Status.Phase = NetworkPending
	}
}

func SetDefaults_ReservedInstance(obj *ReservedInstance) {
	if obj.Status.Phase == "" {
		obj.Status.Phase = ReservedInstanceAvailable
	}
}
