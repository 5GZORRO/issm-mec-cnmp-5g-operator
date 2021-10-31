package common

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	client "sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

// IsBeingDeleted check for deletingTimestamp on the object
func IsBeingDeleted(obj metav1.Object) bool {
	return !obj.GetDeletionTimestamp().IsZero()
}

// HasFinalizer check has the object got the named finalizer
func HasFinalizer(obj metav1.Object, finalizer string) bool {
	for _, fin := range obj.GetFinalizers() {
		if fin == finalizer {
			return true
		}
	}
	return false
}

// AddFinalizer Add a finalizer to an object
func AddFinalizer(obj client.Object, finalizer string) {
	controllerutil.AddFinalizer(obj, finalizer)
}

// RemoveFinalizer Remove finalizer from an object
func RemoveFinalizer(instance client.Object, finalizer string) {
	controllerutil.RemoveFinalizer(instance, finalizer)
}
