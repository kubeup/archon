package controller

import (
	"fmt"
	"github.com/golang/glog"
	"k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/api/meta"
	"k8s.io/kubernetes/pkg/api/unversioned"
	"k8s.io/kubernetes/pkg/api/validation"
	"k8s.io/kubernetes/pkg/client/record"
	"k8s.io/kubernetes/pkg/labels"
	"k8s.io/kubernetes/pkg/runtime"
	"kubeup.com/archon/pkg/clientset"
	"kubeup.com/archon/pkg/cluster"
)

// Reasons for instance events
const (
	// FailedCreateInstanceReason is added in an event and in a replica set condition
	// when a instance for a replica set is failed to be created.
	FailedCreateInstanceReason = "FailedCreate"
	// SuccessfulCreateInstanceReason is added in an event when a instance for a replica set
	// is successfully created.
	SuccessfulCreateInstanceReason = "SuccessfulCreate"
	// FailedDeleteInstanceReason is added in an event and in a replica set condition
	// when a instance for a replica set is failed to be deleted.
	FailedDeleteInstanceReason = "FailedDelete"
	// SuccessfulDeleteInstanceReason is added in an event when a instance for a replica set
	// is successfully deleted.
	SuccessfulDeleteInstanceReason = "SuccessfulDelete"
)

func InstanceKey(instance *cluster.Instance) string {
	return fmt.Sprintf("%v/%v", instance.Namespace, instance.Name)
}

// FilterActiveInstances returns instances that have not terminated.
func FilterActiveInstances(instances []*cluster.Instance) []*cluster.Instance {
	var result []*cluster.Instance
	for _, p := range instances {
		if IsInstanceActive(p) {
			result = append(result, p)
		} else {
			glog.V(4).Infof("Ignoring inactive instance %v/%v in state %v, deletion time %v",
				p.Namespace, p.Name, p.Status.Phase, p.DeletionTimestamp)
		}
	}
	return result
}

func IsInstanceActive(p *cluster.Instance) bool {
	return cluster.InstanceFailed != p.Status.Phase &&
		p.DeletionTimestamp == nil
}

// ActiveInstances type allows custom sorting of instances so a controller can pick the best ones to delete.
type ActiveInstances []*cluster.Instance

func (s ActiveInstances) Len() int      { return len(s) }
func (s ActiveInstances) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

func (s ActiveInstances) Less(i, j int) bool {
	// 2. InstancePending < InstanceUnknown < InstanceRunning
	m := map[cluster.InstancePhase]int{cluster.InstancePending: 0, cluster.InstanceUnknown: 1, cluster.InstanceRunning: 2}
	if m[s[i].Status.Phase] != m[s[j].Status.Phase] {
		return m[s[i].Status.Phase] < m[s[j].Status.Phase]
	}
	// 3. Not ready < ready
	// If only one of the instances is not ready, the not ready one is smaller
	if cluster.IsInstanceReady(s[i]) != cluster.IsInstanceReady(s[j]) {
		return !cluster.IsInstanceReady(s[i])
	}
	// 6. Empty creation time instances < newer instances < older instances
	if !s[i].CreationTimestamp.Equal(s[j].CreationTimestamp) {
		return afterOrZero(s[i].CreationTimestamp, s[j].CreationTimestamp)
	}
	return false
}

// afterOrZero checks if time t1 is after time t2; if one of them
// is zero, the zero time is seen as after non-zero time.
func afterOrZero(t1, t2 unversioned.Time) bool {
	if t1.Time.IsZero() || t2.Time.IsZero() {
		return t1.Time.IsZero()
	}
	return t1.After(t2.Time)
}

// InstanceControlInterface is an interface that knows how to add or delete instances
// created as an interface to allow testing.
type InstanceControlInterface interface {
	// CreateInstances creates new instances according to the spec.
	CreateInstances(namespace string, template *cluster.InstanceTemplateSpec, object runtime.Object) error
	// CreateInstancesOnNode creates a new instance according to the spec on the specified node.
	CreateInstancesOnNode(nodeName, namespace string, template *cluster.InstanceTemplateSpec, object runtime.Object) error
	// CreateInstancesWithControllerRef creates new instances according to the spec, and sets object as the instance's controller.
	CreateInstancesWithControllerRef(namespace string, template *cluster.InstanceTemplateSpec, object runtime.Object, controllerRef *api.OwnerReference) error
	// DeleteInstance deletes the instance identified by instanceID.
	DeleteInstance(namespace string, instanceID string, object runtime.Object) error
	// PatchInstance patches the instance.
	PatchInstance(namespace, name string, data []byte) error
}

// RealInstanceControl is the default implementation of InstanceControlInterface.
type RealInstanceControl struct {
	KubeClient clientset.Interface
	Recorder   record.EventRecorder
}

var _ InstanceControlInterface = &RealInstanceControl{}

func getInstancesLabelSet(template *cluster.InstanceTemplateSpec) labels.Set {
	desiredLabels := make(labels.Set)
	for k, v := range template.Labels {
		desiredLabels[k] = v
	}
	return desiredLabels
}

func getInstancesFinalizers(template *cluster.InstanceTemplateSpec) []string {
	desiredFinalizers := make([]string, len(template.Finalizers))
	copy(desiredFinalizers, template.Finalizers)
	return desiredFinalizers
}

func getInstancesAnnotationSet(template *cluster.InstanceTemplateSpec, object runtime.Object) (labels.Set, error) {
	desiredAnnotations := make(labels.Set)
	for k, v := range template.Annotations {
		desiredAnnotations[k] = v
	}
	createdByRef, err := api.GetReference(object)
	if err != nil {
		return desiredAnnotations, fmt.Errorf("unable to get controller reference: %v", err)
	}

	// TODO: this code was not safe previously - as soon as new code came along that switched to v2, old clients
	//   would be broken upon reading it. This is explicitly hardcoded to v1 to guarantee predictable deployment.
	//   We need to consistently handle this case of annotation versioning.
	codec := api.Codecs.LegacyCodec(unversioned.GroupVersion{Group: api.GroupName, Version: "v1"})

	createdByRefJson, err := runtime.Encode(codec, &api.SerializedReference{
		Reference: *createdByRef,
	})
	if err != nil {
		return desiredAnnotations, fmt.Errorf("unable to serialize controller reference: %v", err)
	}
	desiredAnnotations[api.CreatedByAnnotation] = string(createdByRefJson)
	return desiredAnnotations, nil
}

func getSecretName(controllerName, alias string) string {
	// use the dash (if the name isn't too long) to make the instance name a bit prettier
	name := fmt.Sprintf("%s-%s", controllerName, alias)
	if len(validation.NameIsDNSSubdomain(name, false)) != 0 {
		name = controllerName
	}
	return name
}

func getInstancesPrefix(controllerName string) string {
	// use the dash (if the name isn't too long) to make the instance name a bit prettier
	prefix := fmt.Sprintf("%s-", controllerName)
	if len(validation.NameIsDNSSubdomain(prefix, true)) != 0 {
		prefix = controllerName
	}
	return prefix
}

func (r RealInstanceControl) CreateInstances(namespace string, template *cluster.InstanceTemplateSpec, object runtime.Object) error {
	return r.createInstances("", namespace, template, object, nil)
}

func (r RealInstanceControl) CreateInstancesWithControllerRef(namespace string, template *cluster.InstanceTemplateSpec, controllerObject runtime.Object, controllerRef *api.OwnerReference) error {
	if controllerRef == nil {
		return fmt.Errorf("controllerRef is nil")
	}
	if len(controllerRef.APIVersion) == 0 {
		return fmt.Errorf("controllerRef has empty APIVersion")
	}
	if len(controllerRef.Kind) == 0 {
		return fmt.Errorf("controllerRef has empty Kind")
	}
	if controllerRef.Controller == nil || *controllerRef.Controller != true {
		return fmt.Errorf("controllerRef.Controller is not set")
	}
	return r.createInstances("", namespace, template, controllerObject, controllerRef)
}

func (r RealInstanceControl) CreateInstancesOnNode(nodeName, namespace string, template *cluster.InstanceTemplateSpec, object runtime.Object) error {
	return r.createInstances(nodeName, namespace, template, object, nil)
}

func (r RealInstanceControl) PatchInstance(namespace, name string, data []byte) error {
	_, err := r.KubeClient.Archon().Instances(namespace).Patch(name, api.StrategicMergePatchType, data)
	return err
}

func GetInstanceFromTemplate(template *cluster.InstanceTemplateSpec, parentObject runtime.Object, controllerRef *api.OwnerReference) (*cluster.Instance, error) {
	desiredLabels := getInstancesLabelSet(template)
	desiredFinalizers := getInstancesFinalizers(template)
	desiredAnnotations, err := getInstancesAnnotationSet(template, parentObject)
	if err != nil {
		return nil, err
	}
	accessor, err := meta.Accessor(parentObject)
	if err != nil {
		return nil, fmt.Errorf("parentObject does not have ObjectMeta, %v", err)
	}
	prefix := getInstancesPrefix(accessor.GetName())

	instance := &cluster.Instance{
		TypeMeta: unversioned.TypeMeta{
			APIVersion: "archon.kubeup.com/v1",
			Kind:       "Instance",
		},
		ObjectMeta: api.ObjectMeta{
			Labels:       desiredLabels,
			Annotations:  desiredAnnotations,
			GenerateName: prefix,
			Finalizers:   desiredFinalizers,
		},
	}
	if controllerRef != nil {
		instance.OwnerReferences = append(instance.OwnerReferences, *controllerRef)
	}
	if err := api.Scheme.Convert(&template.Spec, &instance.Spec, nil); err != nil {
		return nil, fmt.Errorf("unable to convert instance template: %v", err)
	}
	return instance, nil
}

func GetSecretsFromTemplate(template *cluster.InstanceTemplateSpec, parentObject runtime.Object, controllerRef *api.OwnerReference) ([]*api.Secret, error) {
	secrets := make([]*api.Secret, 0)
	for _, secret := range template.Secrets {
		accessor, err := meta.Accessor(parentObject)
		if err != nil {
			return nil, fmt.Errorf("parentObject does not have ObjectMeta, %v", err)
		}

		if secret.Annotations == nil {
			secret.Annotations = make(map[string]string)
		}

		if secret.Name != "" {
			secret.Annotations["archon.kubeup.com/alias"] = secret.Name
		}

		if secret.Annotations == nil {
			secret.Annotations = make(map[string]string)
		}
		alias := secret.Annotations["archon.kubeup.com/alias"]
		if alias != "" {
			secret.ObjectMeta.Name = getSecretName(accessor.GetName(), alias)
		} else {
			return nil, fmt.Errorf("secret has no name or alias: %v", secret)
		}

		if controllerRef != nil {
			secret.OwnerReferences = append(secret.OwnerReferences, *controllerRef)
		}
		secrets = append(secrets, &secret)
	}

	return secrets, nil
}

func (r RealInstanceControl) createSecrets(namespace string, template *cluster.InstanceTemplateSpec, object runtime.Object, controllerRef *api.OwnerReference) ([]*api.Secret, error) {
	secretTemplates, err := GetSecretsFromTemplate(template, object, controllerRef)
	if err != nil {
		return nil, err
	}
	secrets := make([]*api.Secret, 0)
	for _, secret := range secretTemplates {
		if newSecret, err := r.KubeClient.Core().Secrets(namespace).Create(secret); err != nil {
			r.Recorder.Eventf(object, api.EventTypeWarning, FailedCreateInstanceReason, "Error creating: %v", err)
			return nil, fmt.Errorf("unable to create secrets: %v", err)
		} else {
			secrets = append(secrets, newSecret)
			accessor, err := meta.Accessor(object)
			if err != nil {
				glog.Errorf("parentObject does not have ObjectMeta, %v", err)
			} else {
				glog.V(4).Infof("Controller %v created secret %v", accessor.GetName(), newSecret.Name)
				r.Recorder.Eventf(object, api.EventTypeNormal, SuccessfulCreateInstanceReason, "Created secret: %v", newSecret.Name)
			}
		}
	}
	return secrets, nil
}

func (r RealInstanceControl) createInstances(nodeName, namespace string, template *cluster.InstanceTemplateSpec, object runtime.Object, controllerRef *api.OwnerReference) error {
	instance, err := GetInstanceFromTemplate(template, object, controllerRef)
	if err != nil {
		return err
	}
	if labels.Set(instance.Labels).AsSelectorPreValidated().Empty() {
		return fmt.Errorf("unable to create instances, no labels")
	}
	if len(template.Secrets) > 0 {
		instance.Status.Phase = cluster.InstanceInitializing
	}
	if newInstance, err := r.KubeClient.Archon().Instances(namespace).Create(instance); err != nil {
		r.Recorder.Eventf(object, api.EventTypeWarning, FailedCreateInstanceReason, "Error creating: %v", err)
		return fmt.Errorf("unable to create instances: %v", err)
	} else {
		accessor, err := meta.Accessor(object)
		if err != nil {
			glog.Errorf("parentObject does not have ObjectMeta, %v", err)
		} else {
			glog.V(4).Infof("Controller %v created instance %v", accessor.GetName(), newInstance.Name)
			r.Recorder.Eventf(object, api.EventTypeNormal, SuccessfulCreateInstanceReason, "Created instance: %v", newInstance.Name)
		}
		secrets, err := r.createSecrets(namespace, template, newInstance, nil)
		if err != nil {
			return fmt.Errorf("unable to create secrets for instance %s", newInstance.Name)
		}
		if len(secrets) > 0 {
			for _, s := range secrets {
				newInstance.Spec.Secrets = append(newInstance.Spec.Secrets, cluster.LocalObjectReference{Name: s.Name})
			}
			newInstance.Status.Phase = cluster.InstancePending
			_, err = r.KubeClient.Archon().Instances(namespace).Update(newInstance)
			if err != nil {
				return fmt.Errorf("unable to save secrets dependency for instance %s", newInstance.Name)
			}
		}
	}
	return nil
}

func (r RealInstanceControl) DeleteInstance(namespace string, instanceID string, object runtime.Object) error {
	accessor, err := meta.Accessor(object)
	if err != nil {
		return fmt.Errorf("object does not have ObjectMeta, %v", err)
	}
	glog.V(2).Infof("Controller %v deleting instance %v/%v", accessor.GetName(), namespace, instanceID)
	if err := r.KubeClient.Archon().Instances(namespace).Delete(instanceID); err != nil {
		r.Recorder.Eventf(object, api.EventTypeWarning, FailedDeleteInstanceReason, "Error deleting: %v", err)
		return fmt.Errorf("unable to delete instances: %v", err)
	} else {
		r.Recorder.Eventf(object, api.EventTypeNormal, SuccessfulDeleteInstanceReason, "Deleted instance: %v", instanceID)
	}
	return nil
}
