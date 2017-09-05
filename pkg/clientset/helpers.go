package clientset

import (
	"fmt"
	"github.com/golang/glog"
	apiextensions "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	apiextensionsclient "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/kubernetes/pkg/api"
	//extensions "k8s.io/kubernetes/pkg/apis/extensions/v1beta1"
	"kubeup.com/archon/pkg/cluster"
	"time"
)

type checker func() error

func (c *Clientset) EnsureResources(timeout time.Duration) (err error) {
	glog.Infof("Ensuring Archon resources")

	return wait.Poll(3*time.Second, timeout, c.ensureResources)
}

func (c *Clientset) ensureResources() (done bool, err error) {
	a := c.Archon()
	ac, err := apiextensionsclient.NewForConfig(c.config)
	if err != nil {
		return
	}

	data := map[string]struct {
		Checker    checker
		Plural     string
		Kind       string
		ShortNames []string
	}{
		"instance":         {Plural: "instances", Kind: "Instance", ShortNames: []string{"ins"}, Checker: func() error { _, err = a.Instances(api.NamespaceAll).List(metav1.ListOptions{}); return err }},
		"instancegroup":    {Plural: "instancegroups", Kind: "InstanceGroup", ShortNames: []string{"ig", "igs"}, Checker: func() error { _, err = a.InstanceGroups(api.NamespaceAll).List(metav1.ListOptions{}); return err }},
		"network":          {Plural: "networks", Kind: "Network", Checker: func() error { _, err = a.Networks(api.NamespaceAll).List(); return err }},
		"user":             {Plural: "users", Kind: "User", Checker: func() error { _, err = a.Users(api.NamespaceAll).List(); return err }},
		"reservedinstance": {Plural: "reservedinstances", Kind: "ReservedInstance", ShortNames: []string{"ri", "ris"}, Checker: func() error { _, err = a.ReservedInstances(api.NamespaceAll).List(metav1.ListOptions{}); return err }},
	}

	for name, v := range data {
		err = v.Checker()

		if errors.IsNotFound(err) {
			tpr := apiextensions.CustomResourceDefinition{
				TypeMeta: metav1.TypeMeta{
					Kind:       "CustomResourceDefinition",
					APIVersion: "apiextensions.k8s.io/v1beta",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name: fmt.Sprintf("%s.%s", v.Plural, cluster.GroupName),
				},
				Spec: apiextensions.CustomResourceDefinitionSpec{
					Group:   cluster.GroupName,
					Version: cluster.GroupVersion,
					Scope:   apiextensions.NamespaceScoped,
					Names: apiextensions.CustomResourceDefinitionNames{
						Plural:     v.Plural,
						Singular:   name,
						Kind:       v.Kind,
						ShortNames: v.ShortNames,
					},
				},
			}

			_, err = ac.ApiextensionsV1beta1().CustomResourceDefinitions().Create(&tpr)
			if err != nil {
				return
			}

			glog.Infof("Resource %s is created", name)
		} else if err != nil {
			glog.Errorf("%v", err)
			return
		}
	}

	glog.Infof("All Archon resource have been ensured")
	return true, nil
}
