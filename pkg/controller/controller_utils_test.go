package controller

import (
	"fmt"
	"math/rand"
	"reflect"
	"sort"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/kubernetes/pkg/api"
	"kubeup.com/archon/pkg/cluster"
	"testing"
)

func TestSortingActiveInstances(t *testing.T) {
	now := metav1.Now()
	then := metav1.Time{Time: now.AddDate(0, -1, 0)}
	instanceParams := []struct {
		Name              string
		Annotation        string
		Phase             cluster.InstancePhase
		Condition         *cluster.InstanceCondition
		CreationTimestamp *metav1.Time
	}{
		{
			Annotation: cluster.InstanceToBeRemovedAnnotation,
			Phase:      cluster.InstanceRunning,
		},
		{
			Phase: cluster.InstancePending,
		},
		{
			Phase: cluster.InstanceUnknown,
		},
		{
			Phase: cluster.InstanceRunning,
		},
		{
			Condition: &cluster.InstanceCondition{Type: cluster.InstanceReady, Status: api.ConditionTrue},
			Phase:     cluster.InstanceRunning,
		},
		{
			Condition:         &cluster.InstanceCondition{Type: cluster.InstanceReady, Status: api.ConditionTrue},
			Phase:             cluster.InstanceRunning,
			CreationTimestamp: &now,
		},
		{
			Condition:         &cluster.InstanceCondition{Type: cluster.InstanceReady, Status: api.ConditionTrue},
			Phase:             cluster.InstanceRunning,
			CreationTimestamp: &then,
		},
	}

	instances := make([]*cluster.Instance, len(instanceParams))
	for i := range instanceParams {
		name := fmt.Sprintf("instance%d", i)
		instanceParams[i].Name = name
		params := instanceParams[i]
		instances[i] = &cluster.Instance{
			ObjectMeta: metav1.ObjectMeta{
				Name:      name,
				Namespace: "default",
				Annotations: map[string]string{
					params.Annotation: "true",
				},
			},
			Status: cluster.InstanceStatus{
				Phase: params.Phase,
			},
		}
		if params.Condition != nil {
			instances[i].Status.Conditions = []cluster.InstanceCondition{*params.Condition}
		}
		if params.CreationTimestamp != nil {
			instances[i].CreationTimestamp = *params.CreationTimestamp
		}
	}

	getOrder := func(instances []*cluster.Instance) []string {
		names := make([]string, len(instances))
		for i := range instances {
			names[i] = instances[i].Name
		}
		return names
	}

	expected := getOrder(instances)
	total := len(instances)

	for i := 0; i < 20; i++ {
		idx := rand.Perm(total)
		randomized := make([]*cluster.Instance, total)
		for j := 0; j < total; j++ {
			randomized[j] = instances[idx[j]]
		}
		sort.Sort(ActiveInstances(randomized))
		actual := getOrder(randomized)

		if !reflect.DeepEqual(actual, expected) {
			t.Errorf("expected %v, got %v", expected, actual)
		}
	}
}
