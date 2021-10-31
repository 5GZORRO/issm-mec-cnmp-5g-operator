package common

import (
	routev1 "github.com/openshift/api/route/v1"

	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/kubernetes"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
)

// VerifyAPI will verify that the given group/version is present in the cluster.
func verifyAPI(group string, version string) (bool, error) {
	cfg, err := config.GetConfig()
	if err != nil {
		// Unable to get KBS Config
		return false, err
	}

	k8s, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		// Unable to create K8s client
		return false, err
	}

	gv := schema.GroupVersion{
		Group:   group,
		Version: version,
	}

	if err = discovery.ServerSupportsVersion(k8s, gv); err != nil {
		// error, API not available
		return false, nil
	}

	// API Exists
	return true, nil
}

func VerifyRouteAPI() (bool, error) {
	found, err := verifyAPI(routev1.GroupName, routev1.SchemeGroupVersion.Version)
	if err != nil {
		return false, err
	}

	return found, nil
}
