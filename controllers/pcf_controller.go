/*
Copyright 2021 IBM.

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

package controllers

import (
	"context"
	"github.ibm.com/Steve-Glover/5GOperator/pkg/reconcilers/pcf"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"strings"

	"github.com/go-logr/logr"
	corev1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	fivegv1alpha1 "github.ibm.com/Steve-Glover/5GOperator/api/v1alpha1"
	"github.ibm.com/Steve-Glover/5GOperator/pkg/common"
)

// PcfReconciler reconciles a Pcf object
type PcfReconciler struct {
	*kubernetes.Clientset
	client.Client
	Log      logr.Logger
	Scheme   *runtime.Scheme
	k8sUtils *common.K8sUtils
}

func NewPcfReconciler(k8scfg *rest.Config, mgr manager.Manager, log logr.Logger, scheme *runtime.Scheme) (*PcfReconciler, error) {
	clientSet, err := kubernetes.NewForConfig(k8scfg)
	if err != nil {
		return nil, err
	}

	r := &PcfReconciler{
		Clientset: clientSet,
		Client:    mgr.GetClient(),
		Log:       log,
		Scheme:    scheme,
		k8sUtils:  common.NewK8sUtils(k8scfg, clientSet, mgr.GetClient(), scheme),
	}

	err = r.SetupWithManager(mgr)
	if err != nil {
		return nil, err
	}

	return r, nil
}

//+kubebuilder:rbac:groups=5g.ibm.com,resources=pcfs,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=5g.ibm.com,resources=pcfs/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=5g.ibm.com,resources=pcfs/finalizers,verbs=update
//+kubebuilder:rbac:groups=core,resources=pods,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=core,resources=services,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=core,resources=configmaps,verbs=get;list;watch;create;update;patch;delete

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Pcf object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.7.2/pkg/reconcile
func (r *PcfReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = r.Log.WithValues("pcf", req.NamespacedName)

	// Fetch the Pcf instance
	instance := &fivegv1alpha1.Pcf{}
	err := r.Client.Get(context.TODO(), req.NamespacedName, instance)
	if err != nil {
		if k8serrors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			return ctrl.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return ctrl.Result{}, err
	}

	// update observedGeneration to match generation
	if instance.Status.ObservedGeneration < instance.ObjectMeta.Generation {
		instance.Status.ObservedGeneration = instance.ObjectMeta.Generation
		err := r.Client.Update(context.TODO(), instance)
		if err != nil {
			r.Log.Error(err, "unable to update instance", "instance", instance)
			return ctrl.Result{}, err
		}
	}

	// Check if CR has been fully initialised
	if ok := r.isInitialized(instance); !ok {
		err := r.Client.Update(context.TODO(), instance)
		if err != nil {
			r.Log.Error(err, "unable to update instance", "instance", instance)
			return ctrl.Result{}, err
		}
		return ctrl.Result{}, nil
	}

	pcfReconciler := pcf.PcfReconcilerFactory{K8sUtils: r.k8sUtils, Client: r.Client}.Reconciler()

	// If is CR deletion, run cleanup / finalization logic
	if common.IsBeingDeleted(instance) {
		if !common.HasFinalizer(instance, common.FinalizerName) {
			return ctrl.Result{}, nil
		}
		pcfReconciler.Finalize(req, instance)
		if err != nil {
			r.Log.Error(err, "unable to delete instance", "instance", instance)
			return ctrl.Result{}, err
		}
		common.RemoveFinalizer(instance, common.FinalizerName)
		err = r.Client.Update(context.TODO(), instance)
		if err != nil {
			r.Log.Error(err, "unable to update instance", "instance", instance)
			return ctrl.Result{}, err
		}
		return ctrl.Result{}, nil
	}

	// Reconcile Pcf
	opResult, result, reconcileErr := pcfReconciler.Reconcile(req, instance)

	err = r.updateStatus(req, pcfReconciler, opResult, instance, reconcileErr)
	if err != nil {
		if !strings.Contains(err.Error(), "the object has been modified") {
			r.Log.Error(err, "unable to update PCF status", "status", instance)
			return result, err
		}
	}

	return ctrl.Result{}, nil
}

func (r *PcfReconciler) updateStatus(request ctrl.Request, reconciler common.FiveGReconciler, opResult controllerutil.OperationResult,
	instance *fivegv1alpha1.Pcf, reconcileErr error) error {
	reconciler.UpdateStatus(opResult, instance, reconcileErr)

	err := r.Client.Status().Update(context.Background(), instance)
	if err != nil {
		return err
	}

	return nil
}

func (r *PcfReconciler) isInitialized(obj metav1.Object) bool {
	mycrd, ok := obj.(*fivegv1alpha1.Pcf)
	if !ok {
		return false
	}
	if common.HasFinalizer(mycrd, common.FinalizerName) {
		return true
	}
	common.AddFinalizer(mycrd, common.FinalizerName)
	return false
}

// SetupWithManager sets up the controller with the Manager.
func (r *PcfReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&fivegv1alpha1.Pcf{}).
		Owns(&corev1.Pod{}).
		Owns(&corev1.ConfigMap{}).
		Owns(&corev1.Service{}).
		Complete(r)
}
