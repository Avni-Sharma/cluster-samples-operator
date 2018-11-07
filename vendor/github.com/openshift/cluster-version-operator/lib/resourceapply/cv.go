package resourceapply

import (
	"github.com/openshift/cluster-version-operator/lib/resourcemerge"
	cvv1 "github.com/openshift/cluster-version-operator/pkg/apis/config.openshift.io/v1"
	cvclientv1 "github.com/openshift/cluster-version-operator/pkg/generated/clientset/versioned/typed/config.openshift.io/v1"
	cvlistersv1 "github.com/openshift/cluster-version-operator/pkg/generated/listers/config.openshift.io/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/pointer"
)

func ApplyClusterVersion(client cvclientv1.ClusterVersionsGetter, required *cvv1.ClusterVersion) (*cvv1.ClusterVersion, bool, error) {
	existing, err := client.ClusterVersions().Get(required.Name, metav1.GetOptions{})
	if errors.IsNotFound(err) {
		actual, err := client.ClusterVersions().Create(required)
		return actual, true, err
	}
	if err != nil {
		return nil, false, err
	}

	modified := pointer.BoolPtr(false)
	resourcemerge.EnsureClusterVersion(modified, existing, *required)
	if !*modified {
		return existing, false, nil
	}

	actual, err := client.ClusterVersions().Update(existing)
	return actual, true, err
}

func ApplyClusterVersionFromCache(lister cvlistersv1.ClusterVersionLister, client cvclientv1.ClusterVersionsGetter, required *cvv1.ClusterVersion) (*cvv1.ClusterVersion, bool, error) {
	obj, err := lister.Get(required.Name)
	if errors.IsNotFound(err) {
		actual, err := client.ClusterVersions().Create(required)
		return actual, true, err
	}
	if err != nil {
		return nil, false, err
	}

	// Don't want to mutate cache.
	existing := obj.DeepCopy()
	modified := pointer.BoolPtr(false)
	resourcemerge.EnsureClusterVersion(modified, existing, *required)
	if !*modified {
		return existing, false, nil
	}

	actual, err := client.ClusterVersions().Update(existing)
	return actual, true, err
}
