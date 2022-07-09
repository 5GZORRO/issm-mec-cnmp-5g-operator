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
	"fmt"
	"github.com/go-logr/logr"
	"github.com/5GZORRO/issm-mec-cnmp-5g-operator/pkg/reconcilers/amf"
	"github.com/5GZORRO/issm-mec-cnmp-5g-operator/pkg/reconcilers/ausf"
	"github.com/5GZORRO/issm-mec-cnmp-5g-operator/pkg/reconcilers/mongo"
	"github.com/5GZORRO/issm-mec-cnmp-5g-operator/pkg/reconcilers/n3iwf"
	"github.com/5GZORRO/issm-mec-cnmp-5g-operator/pkg/reconcilers/nrf"
	"github.com/5GZORRO/issm-mec-cnmp-5g-operator/pkg/reconcilers/nssf"
	"github.com/5GZORRO/issm-mec-cnmp-5g-operator/pkg/reconcilers/pcf"
	"github.com/5GZORRO/issm-mec-cnmp-5g-operator/pkg/reconcilers/smf"
	"github.com/5GZORRO/issm-mec-cnmp-5g-operator/pkg/reconcilers/transition"
	"github.com/5GZORRO/issm-mec-cnmp-5g-operator/pkg/reconcilers/udm"
	"github.com/5GZORRO/issm-mec-cnmp-5g-operator/pkg/reconcilers/udr"
	"github.com/5GZORRO/issm-mec-cnmp-5g-operator/pkg/reconcilers/upf"
	"github.com/5GZORRO/issm-mec-cnmp-5g-operator/pkg/reconcilers/webconsole"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"strings"
	"time"

	fivegv1alpha1 "github.com/5GZORRO/issm-mec-cnmp-5g-operator/api/v1alpha1"
	"github.com/5GZORRO/issm-mec-cnmp-5g-operator/pkg/common"
)

// TransitionReconciler reconciles a Transition object
type TransitionReconciler struct {
	RestConfig *rest.Config
	*kubernetes.Clientset
	client.Client
	Log      logr.Logger
	Scheme   *runtime.Scheme
	k8sUtils *common.K8sUtils
	crMap    map[string]common.EntityHandler
}

func NewTransitionReconciler(k8scfg *rest.Config, mgr manager.Manager, log logr.Logger, scheme *runtime.Scheme) (*TransitionReconciler, error) {
	clientSet, err := kubernetes.NewForConfig(k8scfg)
	if err != nil {
		return nil, err
	}

	k8sUtils := common.NewK8sUtils(k8scfg, clientSet, mgr.GetClient(), scheme)

	crMap := make(map[string]common.EntityHandler)
	crMap["amfs.v1alpha1.5g.ibm.com"] = amf.NewAmf(k8sUtils, k8sUtils, mgr.GetClient())
	crMap["ausfs.v1alpha1.5g.ibm.com"] = ausf.NewAusf(k8sUtils, k8sUtils, mgr.GetClient())
	crMap["mongoes.v1alpha1.5g.ibm.com"] = mongo.NewMongo(k8sUtils, k8sUtils, mgr.GetClient())
	crMap["n3iwfs.v1alpha1.5g.ibm.com"] = n3iwf.NewN3iwf(k8sUtils, k8sUtils, mgr.GetClient())
	crMap["nrfs.v1alpha1.5g.ibm.com"] = nrf.NewNrf(k8sUtils, k8sUtils, mgr.GetClient())
	crMap["nssfs.v1alpha1.5g.ibm.com"] = nssf.NewNssf(k8sUtils, k8sUtils, mgr.GetClient())
	crMap["pcfs.v1alpha1.5g.ibm.com"] = pcf.NewPcf(k8sUtils, k8sUtils, mgr.GetClient())
	crMap["smfs.v1alpha1.5g.ibm.com"] = smf.NewSmf(k8sUtils, k8sUtils, mgr.GetClient())
	crMap["udms.v1alpha1.5g.ibm.com"] = udm.NewUdm(k8sUtils, k8sUtils, mgr.GetClient())
	crMap["udrs.v1alpha1.5g.ibm.com"] = udr.NewUdr(k8sUtils, k8sUtils, mgr.GetClient())
	crMap["upfs.v1alpha1.5g.ibm.com"] = upf.NewUpf(k8sUtils, k8sUtils, mgr.GetClient())
	crMap["webconsoles.v1alpha1.5g.ibm.com"] = webconsole.NewWebconsole(k8sUtils, k8sUtils, mgr.GetClient())

	r := &TransitionReconciler{
		Clientset: clientSet,
		Client:    mgr.GetClient(),
		Log:       log,
		Scheme:    scheme,
		k8sUtils:  k8sUtils,
		crMap:     crMap,
	}

	err = r.SetupWithManager(mgr)
	if err != nil {
		return nil, err
	}

	return r, nil
}

//+kubebuilder:rbac:groups=5g.ibm.com,resources=transitions,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=5g.ibm.com,resources=transitions/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=5g.ibm.com,resources=transitions/finalizers,verbs=update
//+kubebuilder:rbac:groups=core,resources=pods,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=core,resources=pods/exec,verbs=create
//+kubebuilder:rbac:groups=core,resources=services,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=core,resources=configmaps,verbs=get;list;watch;create;update;patch;delete

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Transition object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.7.2/pkg/reconcile
func (r *TransitionReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = r.Log.WithValues("transition", req.NamespacedName)

	// Fetch the Transition instance
	transitionCR := &fivegv1alpha1.Transition{}
	err := r.Client.Get(context.TODO(), req.NamespacedName, transitionCR)
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
	if transitionCR.Status.ObservedGeneration < transitionCR.ObjectMeta.Generation {
		transitionCR.Status.ObservedGeneration = transitionCR.ObjectMeta.Generation
		err := r.Client.Update(context.TODO(), transitionCR)
		if err != nil {
			r.Log.Error(err, "unable to update instance", "instance", transitionCR)
			return ctrl.Result{}, err
		}
	}

	// Check if CR has been fully initialised
	if ok := r.isInitialized(transitionCR); !ok {
		err := r.Client.Update(context.TODO(), transitionCR)
		if err != nil {
			r.Log.Error(err, "unable to update instance", "transitionCR", transitionCR)
			return ctrl.Result{}, err
		}
		return ctrl.Result{}, nil
	}

	t := transition.NewTransition(r.k8sUtils)

	// If is CR deletion, run cleanup / finalization logic
	if common.IsBeingDeleted(transitionCR) {
		if !common.HasFinalizer(transitionCR, common.FinalizerName) {
			return ctrl.Result{}, nil
		}
		t.Finalize(req, transitionCR)
		if err != nil {
			r.Log.Error(err, "unable to delete instance", "instance", transitionCR)
			return ctrl.Result{}, err
		}
		common.RemoveFinalizer(transitionCR, common.FinalizerName)
		err = r.Client.Update(context.TODO(), transitionCR)
		if err != nil {
			r.Log.Error(err, "unable to update instance", "instance", transitionCR)
			return ctrl.Result{}, err
		}
		return ctrl.Result{}, nil
	}

	// Reconcile Transition

	entityHandler, err := r.getEntityHandler(transitionCR)
	if err != nil {
		r.Log.Error(err, "Unable to get transition status", "transitionCR", transitionCR)
		return ctrl.Result{}, err
	}

	entity, err := r.getEntity(transitionCR)
	if err != nil {
		r.Log.Error(err, "Unable to get entity for transition", "transitionCR", transitionCR)
		return ctrl.Result{}, err
	}

	crState, err := common.GetCRState(entityHandler, transitionCR)
	if err != nil {
		r.Log.Error(err, "Unable to get CR status", "instance", transitionCR)
		return ctrl.Result{}, err
	}

	switch strings.ToLower(transitionCR.Spec.Config.TransitionName) {
	case "create":
		// special handling for create
		switch crState {
		case fivegv1alpha1.CRNotFound:
			if transitionCR.Status.IsRunning() {
				// already running - but keep checking
				return ctrl.Result{RequeueAfter: time.Second, Requeue: true}, nil
			}

			transitionCR.Status.SetRunning(true)

			// create the entity
			_, err := entityHandler.Create(transitionCR)
			if err != nil {
				transitionCR.Status.SetError(err)
			} else {
				transitionCR.Status.SetReconciled(false)
			}

			err = r.updateTransition(transitionCR)
			if err != nil {
				return ctrl.Result{}, err
			}
			return ctrl.Result{RequeueAfter: time.Second, Requeue: true}, nil
		case fivegv1alpha1.CRFound:
			return ctrl.Result{RequeueAfter: time.Second, Requeue: true}, nil
		case fivegv1alpha1.CRReconciled:
			fallthrough
		case fivegv1alpha1.CRReady:
			//r.Log.Info("Create CRReady1", "transitionCR", transitionCR)
			transitionCR.Status.SetReady(true)
			transitionCR.Status.SetReconciled(true)
			transitionCR.Status.SetRunning(false)
			//r.Log.Info("Create CRReady2", "transitionCR", transitionCR)
			err = r.updateTransition(transitionCR)
			return ctrl.Result{}, err
		default:
			err = fmt.Errorf("Cannot handle CR state %s for transition %+v", crState, transitionCR)
			transitionCR.Status.SetError(err)
			err = r.updateTransition(transitionCR)
			return ctrl.Result{}, err
		}
	default:
		// all other transitions/operations are handled here
		switch crState {
		case fivegv1alpha1.CRNotFound:
			err = fmt.Errorf("CR not found for transition %s", transitionCR.Name)
			transitionCR.Status.SetError(err)
			err = r.updateTransition(transitionCR)
			return ctrl.Result{}, err
		case fivegv1alpha1.CRFound:
			r.Log.Info("CR %+v not ready for transition %+v", entity, transitionCR)
			return ctrl.Result{}, nil
		case fivegv1alpha1.CRReconciled:
			// CR is reconciled but not ready, can run transition

			r.Log.Info("Checking to run transition", "transitionCR", transitionCR, "isRunning", transitionCR.Status.IsRunning(),
				"isReady", transitionCR.Status.IsReady())

			if transitionCR.Status.IsRunning() {
				// already running - but keep checking
				return ctrl.Result{RequeueAfter: time.Second, Requeue: true}, nil
			}

			if transitionCR.Status.IsReady() {
				// done
				return ctrl.Result{}, nil
			}

			r.Log.Info("Running transition against CR", "transitionCR", transitionCR, "entity", entity)

			transitionCR.Status.SetRunning(true)
			err := r.updateTransition(transitionCR)
			if err != nil {
				return ctrl.Result{RequeueAfter: time.Second, Requeue: true}, err
			}

			err = entityHandler.RunTransition(transitionCR, entity, r.updateStatusForCompletedTransition)
			if err != nil {
				transitionCR.Status.SetError(err)
				err = r.updateTransition(transitionCR)
				if err != nil {
					return ctrl.Result{}, err
				}
				entity.SetError(err)
				err = r.updateResource(entity)
				return ctrl.Result{}, err
			}
		case fivegv1alpha1.CRReady:
			r.Log.Info("CR %+v ready, not running transition %+v", "entity", entity, "transitionCR", transitionCR)
			return ctrl.Result{}, err
		default:
			err = fmt.Errorf("Cannot handle CR state %s for transition %+v", crState, transitionCR)
			transitionCR.Status.SetError(err)
			err := r.updateTransition(transitionCR)
			if err != nil {
				return ctrl.Result{}, err
			}
			return ctrl.Result{}, err
		}
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *TransitionReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&fivegv1alpha1.Transition{}).
		Complete(r)
}

func (r *TransitionReconciler) getEntityHandler(transitionCR *fivegv1alpha1.Transition) (common.EntityHandler, error) {
	resource, exists := r.crMap[strings.ToLower(transitionCR.Spec.Config.ResourceType)]
	if exists {
		return resource, nil
	} else {
		return nil, fmt.Errorf("Unknown resource type %s", transitionCR.Spec.Config.ResourceType)
	}
}

func (r *TransitionReconciler) getEntity(transitionCR *fivegv1alpha1.Transition) (common.Entity, error) {
	crHandler, exists := r.crMap[strings.ToLower(transitionCR.Spec.Config.ResourceType)]
	if exists {
		entity, err := crHandler.GetCR(transitionCR)
		if err != nil {
			if k8serrors.IsNotFound(err) {
				return nil, nil
			}
			return nil, err
		}
		return entity, nil
	} else {
		return nil, fmt.Errorf("Unknown resource type %s", transitionCR.Spec.Config.ResourceType)
	}
}

func (r *TransitionReconciler) updateStatusForCompletedTransition(entity common.Entity, transitionCR *fivegv1alpha1.Transition, reconcileErr error) error {
	entityHandler, err := r.getEntityHandler(transitionCR)
	if err != nil {
		r.Log.Error(err, "Unable to get CR for transition", "transitionCR", transitionCR)
		return err
	}

	transitionCR.Status.SetRunning(false)
	if reconcileErr != nil {
		transitionCR.Status.SetError(reconcileErr)
	} else {
		transitionCR.Status.SetReconciled(true)
		transitionCR.Status.SetReady(true)
	}

	err = r.updateTransition(transitionCR)
	if err != nil {
		return err
	}

	if reconcileErr != nil {
		entity.SetError(reconcileErr)
	} else {
		isReady, err := entityHandler.CheckReady(transitionCR, entity)
		if err != nil {
			return err
		}
		if isReady {
			entity.SetReady(true)
		}
	}

	err = r.updateResource(entity)
	if err != nil {
		return err
	}

	return nil
}

func (r *TransitionReconciler) updateTransition(transitionCR *fivegv1alpha1.Transition) error {
	err := r.Client.Status().Update(context.Background(), transitionCR)
	if err != nil {
		if !strings.Contains(err.Error(), "the object has been modified") {
			r.Log.Error(err, "unable to update Transition status", "transition", transitionCR)
			return err
		}
	}

	return nil
}

func (r *TransitionReconciler) updateResource(entity common.Entity) error {
	err := r.Client.Status().Update(context.Background(), entity.GetObject())
	if err != nil {
		if !strings.Contains(err.Error(), "the object has been modified") {
			r.Log.Error(err, "unable to update Transition status", "resource", entity)
			return err
		}
	}

	return nil
}

func (r *TransitionReconciler) isInitialized(obj metav1.Object) bool {
	mycrd, ok := obj.(*fivegv1alpha1.Transition)
	if !ok {
		return false
	}
	if common.HasFinalizer(mycrd, common.FinalizerName) {
		return true
	}
	common.AddFinalizer(mycrd, common.FinalizerName)
	return false
}
