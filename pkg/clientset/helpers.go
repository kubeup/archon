package clientset

import (
	"fmt"
	"github.com/golang/glog"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/kubernetes/pkg/api"
	extensions "k8s.io/kubernetes/pkg/apis/extensions/v1beta1"
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
	data := map[string]checker{
		"instance":          func() error { _, err = a.Instances(api.NamespaceAll).List(metav1.ListOptions{}); return err },
		"instance-group":    func() error { _, err = a.InstanceGroups(api.NamespaceAll).List(metav1.ListOptions{}); return err },
		"network":           func() error { _, err = a.Networks(api.NamespaceAll).List(); return err },
		"user":              func() error { _, err = a.Users(api.NamespaceAll).List(); return err },
		"reserved-instance": func() error { _, err = a.ReservedInstances(api.NamespaceAll).List(metav1.ListOptions{}); return err },
	}

	for name, check := range data {
		err = check()

		if errors.IsNotFound(err) {
			tpr := extensions.ThirdPartyResource{
				TypeMeta: metav1.TypeMeta{
					Kind:       "ThirdPartyResource",
					APIVersion: "v1/betav1",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name: fmt.Sprintf("%s.%s", name, cluster.GroupName),
				},
				Versions: []extensions.APIVersion{
					extensions.APIVersion{
						Name: cluster.GroupVersion,
					},
				},
			}

			_, err = c.Extensions().ThirdPartyResources().Create(&tpr)
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
