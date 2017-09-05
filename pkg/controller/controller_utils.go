package controller

import (
	"fmt"
	"github.com/golang/glog"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/record"
	"k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/api/v1"
	"k8s.io/kubernetes/pkg/api/v1/ref"
	"k8s.io/kubernetes/pkg/api/validation"
	"kubeup.com/archon/pkg/clientset"
	"kubeup.com/archon/pkg/cluster"
	"kubeup.com/archon/pkg/initializer"
	"kubeup.com/archon/pkg/util"
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

var (
	CSRToken = cluster.AnnotationPrefix + "csr"
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
	removeMark1 := cluster.IsInstanceToBeRemoved(s[i])
	removeMark2 := cluster.IsInstanceToBeRemoved(s[j])
	if removeMark1 != removeMark2 {
		// The one to be removed is less
		return removeMark1
	}

	// InstancePending < InstanceUnknown < InstanceRunning
	m := map[cluster.InstancePhase]int{cluster.InstancePending: 0, cluster.InstanceUnknown: 1, cluster.InstanceRunning: 2}
	if m[s[i].Status.Phase] != m[s[j].Status.Phase] {
		return m[s[i].Status.Phase] < m[s[j].Status.Phase]
	}
	// Not ready < ready
	// If only one of the instances is not ready, the not ready one is smaller
	if cluster.IsInstanceReady(s[i]) != cluster.IsInstanceReady(s[j]) {
		return !cluster.IsInstanceReady(s[i])
	}
	// Empty creation time instances < newer instances < older instances
	if !s[i].CreationTimestamp.Equal(s[j].CreationTimestamp) {
		return afterOrZero(s[i].CreationTimestamp, s[j].CreationTimestamp)
	}
	return false
}

// afterOrZero checks if time t1 is after time t2; if one of them
// is zero, the zero time is seen as after non-zero time.
func afterOrZero(t1, t2 metav1.Time) bool {
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
	CreateInstancesWithControllerRef(namespace string, template *cluster.InstanceTemplateSpec, object runtime.Object, controllerRef *metav1.OwnerReference) error
	// DeleteInstance deletes the instance identified by instanceID.
	DeleteInstance(namespace string, instanceID string, object runtime.Object) error
	// PatchInstance patches the instance.
	PatchInstance(namespace, name string, data []byte) error
	// Unbind
	UnbindInstanceWithReservedInstance(instance *cluster.Instance) error
	BindReservedInstance(ri *cluster.ReservedInstance, instance *cluster.Instance, unbind bool) (updatedRI *cluster.ReservedInstance, err error)
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

// WIP: not used.
func getInstancesInitializers(template *cluster.InstanceTemplateSpec) *metav1.Initializers {
	ret := &metav1.Initializers{}
	if template.Initializers != nil && len(template.Initializers.Pending) > 0 {
		ret.Pending = make([]metav1.Initializer, len(template.Initializers.Pending))
		copy(ret.Pending, template.Initializers.Pending)
	} else {
		ret.Result = &metav1.Status{
			Status: metav1.StatusSuccess,
		}
	}
	return ret
}

func getInstancesAnnotationSet(template *cluster.InstanceTemplateSpec, object runtime.Object) (labels.Set, error) {
	desiredAnnotations := make(labels.Set)
	for k, v := range template.Annotations {
		desiredAnnotations[k] = v
	}
	createdByRef, err := ref.GetReference(api.Scheme, object)
	if err != nil {
		return desiredAnnotations, fmt.Errorf("unable to get controller reference: %v", err)
	}

	// TODO: this code was not safe previously - as soon as new code came along that switched to v2, old clients
	//   would be broken upon reading it. This is explicitly hardcoded to v1 to guarantee predictable deployment.
	//   We need to consistently handle this case of annotation versioning.
	codec := api.Codecs.LegacyCodec(schema.GroupVersion{Group: v1.GroupName, Version: "v1"})

	createdByRefJson, err := runtime.Encode(codec, &v1.SerializedReference{
		Reference: *createdByRef,
	})
	if err != nil {
		return desiredAnnotations, fmt.Errorf("unable to serialize controller reference: %v", err)
	}
	desiredAnnotations[v1.CreatedByAnnotation] = string(createdByRefJson)
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

func (r RealInstanceControl) CreateInstancesWithControllerRef(namespace string, template *cluster.InstanceTemplateSpec, controllerObject runtime.Object, controllerRef *metav1.OwnerReference) error {
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
	_, err := r.KubeClient.Archon().Instances(namespace).Patch(name, types.StrategicMergePatchType, data)
	return err
}

func GetInstanceFromTemplate(template *cluster.InstanceTemplateSpec, parentObject runtime.Object, controllerRef *metav1.OwnerReference) (*cluster.Instance, error) {
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
		TypeMeta: metav1.TypeMeta{
			APIVersion: "archon.kubeup.com/v1",
			Kind:       "Instance",
		},
		ObjectMeta: metav1.ObjectMeta{
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

func GetSecretsFromTemplate(template *cluster.InstanceTemplateSpec, parentObject runtime.Object, controllerRef *metav1.OwnerReference) ([]v1.Secret, error) {
	secrets := make([]v1.Secret, 0)
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

		secret.Annotations["archon.kubeup.com/instance"] = accessor.GetName()

		if controllerRef != nil {
			secret.OwnerReferences = append(secret.OwnerReferences, *controllerRef)
		}
		secrets = append(secrets, secret)
	}

	return secrets, nil
}

func (r RealInstanceControl) createSecrets(namespace string, template *cluster.InstanceTemplateSpec, object runtime.Object, controllerRef *metav1.OwnerReference) ([]*v1.Secret, error) {
	secretTemplates, err := GetSecretsFromTemplate(template, object, controllerRef)
	if err != nil {
		return nil, err
	}
	secrets := make([]*v1.Secret, 0)
	for _, secret := range secretTemplates {
		if newSecret, err := r.KubeClient.Core().Secrets(namespace).Create(&secret); err != nil {
			r.Recorder.Eventf(object, api.EventTypeWarning, FailedCreateInstanceReason, "Error creating: %v", err)
			// Still return secrets so we can revert creation
			return secrets, fmt.Errorf("unable to create secrets: %v", err)
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

func (r RealInstanceControl) createInstances(nodeName, namespace string, template *cluster.InstanceTemplateSpec, object runtime.Object, controllerRef *metav1.OwnerReference) (err error) {
	var (
		newInstance *cluster.Instance
		secrets     []*v1.Secret
	)
	defer func() {
		// Undo in case of error
		if err != nil {
			if newInstance != nil {
				err2 := r.UnbindInstanceWithReservedInstance(newInstance)
				if err2 != nil {
					glog.Errorf("Unable to undo bindReservedInstance: %v", err2)
				}
				err2 = r.KubeClient.Archon().Instances(namespace).Delete(newInstance.Name)
				if err2 != nil {
					glog.Errorf("Unable to undo createInstances: %v", err2)
				}
			}
			if len(secrets) > 0 {
				for _, s := range secrets {
					err2 := r.KubeClient.Core().Secrets(namespace).Delete(s.Name, &metav1.DeleteOptions{})
					if err2 != nil {
						glog.Errorf("Unable to undo creating secrets: %v", err2)
					}
				}
			}
		}
	}()
	instance, err := GetInstanceFromTemplate(template, object, controllerRef)
	if err != nil {
		return err
	}
	if labels.Set(instance.Labels).AsSelectorPreValidated().Empty() {
		return fmt.Errorf("unable to create instances, no labels")
	}

	if len(template.Secrets) > 0 {
		// Automatically add CSR initializer so secrets can be processed
		initializer.AddInitializer(instance, CSRToken)
	}

	// Until all required secrets are created and reserved instance is bound, we
	// want InstanceController to ignore this instance. So we put it in Initializing phase.
	instance.Status.Phase = cluster.InstanceInitializing

	if newInstance, err = r.KubeClient.Archon().Instances(namespace).Create(instance); err != nil {
		r.Recorder.Eventf(object, api.EventTypeWarning, FailedCreateInstanceReason, "Error creating: %v", err)
		return fmt.Errorf("unable to create instances: %v", err)
	}

	if err = r.handleProvisionPolicy(newInstance, object); err != nil {
		return fmt.Errorf("Unable to provision instance: %v", err)
	}
	glog.Infof("Handle provision policy: %v", newInstance.Annotations)

	accessor, err := meta.Accessor(object)
	if err != nil {
		glog.Errorf("parentObject does not have ObjectMeta, %v", err)
	} else {
		glog.V(4).Infof("Controller %v created instance %v", accessor.GetName(), newInstance.Name)
		r.Recorder.Eventf(object, api.EventTypeNormal, SuccessfulCreateInstanceReason, "Created instance: %v", newInstance.Name)
	}

	secrets, err = r.createSecrets(namespace, template, newInstance, nil)
	if err != nil {
		return fmt.Errorf("unable to create secrets for instance %s", newInstance.Name)
	}
	if len(secrets) > 0 {
		for _, s := range secrets {
			newInstance.Spec.Secrets = append(newInstance.Spec.Secrets, cluster.LocalObjectReference{Name: s.Name})
		}
	}

	newInstance.Status.Phase = cluster.InstancePending
	_, err = r.KubeClient.Archon().Instances(namespace).Update(newInstance)
	if err != nil {
		return fmt.Errorf("unable to save instance %s: %+v", newInstance.Name, err)
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

func (r RealInstanceControl) UnbindInstanceWithReservedInstance(instance *cluster.Instance) (err error) {
	_, err = r.BindReservedInstance(nil, instance, true)
	return
}

// Update both reserved instance and instance. Get related reserved instance if not provided.
// Update only instance if no reserved instance provided or found
func (r RealInstanceControl) BindReservedInstance(ri *cluster.ReservedInstance, instance *cluster.Instance, unbind bool) (updatedRI *cluster.ReservedInstance, err error) {
	client := r.KubeClient.Archon().ReservedInstances(instance.Namespace)
	ref := instance.Spec.ReservedInstanceRef

	if ri == nil && ref != nil && ref.Name != "" {
		ri, err = client.Get(ref.Name)
		if err != nil {
			return
		}
	}

	if ri == nil && unbind == false {
		err = fmt.Errorf("Unable to bind a nil reversed instance")
		return
	}

	options := cluster.InstanceOptions{}
	err = util.MapToStruct(instance.Labels, &options, cluster.AnnotationPrefix)
	if err != nil {
		err = fmt.Errorf("Can't get instance options: %s", err.Error())
		return
	}

	useInstanceID := ""
	ref = nil
	reclaimPolicy := cluster.InstanceReclaimDelete

	// Update reserved instance
	if ri != nil {
		if unbind == true {
			ri.Status.InstanceName = ""
			ri.Status.Phase = cluster.ReservedInstanceAvailable
		} else {
			ri.Status.InstanceName = instance.Name
			ri.Status.Phase = cluster.ReservedInstanceBound
			reclaimPolicy = cluster.InstanceReclaimRecycle
			ref = &cluster.LocalObjectReference{
				Name: ri.Name,
			}
			useInstanceID = ri.Spec.InstanceID
			reclaimPolicy = cluster.InstanceReclaimRecycle
		}

		updatedRI, err = client.Update(ri)
		if err != nil {
			return
		}
	}

	// Update instance
	instance.Spec.ReservedInstanceRef = ref
	instance.Spec.ReclaimPolicy = reclaimPolicy
	options.UseInstanceID = useInstanceID

	if instance.Labels == nil {
		instance.Labels = make(map[string]string)
	}

	err = util.StructToMap(&options, instance.Labels, cluster.AnnotationPrefix)
	if err != nil {
		err = fmt.Errorf("Unable to update instance labels: %v", err)
		return
	}
	return
}

func (r RealInstanceControl) handleProvisionPolicy(instance *cluster.Instance, object runtime.Object) (err error) {
	var boundRI *cluster.ReservedInstance
	client := r.KubeClient.Archon().ReservedInstances(instance.Namespace)

	defer func() {
		if boundRI != nil && err != nil {
			_, err2 := r.BindReservedInstance(boundRI, instance, true)
			if err2 != nil {
				glog.Errorf("Unabel to revert locked reserved instance: %s, %v", boundRI.Name, err2)
			}
		}
	}()

	ig, ok := object.(*cluster.InstanceGroup)
	if !ok {
		return fmt.Errorf("Failed to get provision policy. Instance's controller object is not an InstanceGroup: %+v", object)
	}

	switch ig.Spec.ProvisionPolicy {
	case cluster.InstanceGroupProvisionReservedFirst, cluster.InstanceGroupProvisionReservedOnly:
		// Get an available reserved instance
		labelSelector := ""
		if ig.Spec.ReservedInstanceSelector != nil {
			labelSelector = ig.Spec.ReservedInstanceSelector.String()
		}
		ril, err := client.List(metav1.ListOptions{
			LabelSelector: labelSelector,
		})
		if err != nil {
			return fmt.Errorf("Unable to list reserved instances: %v", err)
		}
		for _, ri := range ril.Items {
			if ri.Status.Phase != cluster.ReservedInstanceAvailable {
				continue
			}

			// Try binding it
			boundRI, err = r.BindReservedInstance(&ri, instance, false)
			if err == nil {
				cluster.ReservedInstanceToInstance(boundRI, instance)
				glog.V(2).Infof("Bind instance %v with reserved %v", instance.Name, boundRI.Name)
				return err
			}

			// Otherwise we try the next one
		}

		// No match
		if ig.Spec.ProvisionPolicy == cluster.InstanceGroupProvisionReservedOnly {
			return fmt.Errorf("No available reserved instance for instance group: %v", ig.Name)
		}
		glog.V(2).Infof("Instance %v doesn't have a matching reserved instance. Will dynamically provision", instance.Name)
	case cluster.InstanceGroupProvisionDynamicOnly:
		return
	default:
		err = fmt.Errorf("Unsupported provisionPolicy: %v", ig.Spec.ProvisionPolicy)
	}

	return
}
