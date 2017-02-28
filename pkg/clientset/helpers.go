package clientset

import (
	"fmt"
	"github.com/golang/glog"
	"k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/api/errors"
	"k8s.io/kubernetes/pkg/api/unversioned"
	"k8s.io/kubernetes/pkg/apis/extensions"
	"k8s.io/kubernetes/pkg/util/wait"
	"kubeup.com/archon/pkg/clientset/archon"
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
		"instance":       func() error { _, err = a.Instances(api.NamespaceAll).List(api.ListOptions{}); return err },
		"instance-group": func() error { _, err = a.InstanceGroups(api.NamespaceAll).List(api.ListOptions{}); return err },
		"network":        func() error { _, err = a.Networks(api.NamespaceAll).List(); return err },
		"user":           func() error { _, err = a.Users(api.NamespaceAll).List(); return err },
	}

	for name, check := range data {
		err = check()

		if errors.IsNotFound(err) {
			tpr := extensions.ThirdPartyResource{
				TypeMeta: unversioned.TypeMeta{
					Kind:       "ThirdPartyResource",
					APIVersion: "v1/betav1",
				},
				ObjectMeta: api.ObjectMeta{
					Name: fmt.Sprintf("%s.%s", name, archon.GroupName),
				},
				Versions: []extensions.APIVersion{
					extensions.APIVersion{
						Name: archon.GroupVersion,
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
