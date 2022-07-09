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

package subscriber

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.ibm.com/Steve-Glover/5GOperator/api/v1alpha1"
	fivegv1alpha1 "github.ibm.com/Steve-Glover/5GOperator/api/v1alpha1"
	"github.ibm.com/Steve-Glover/5GOperator/pkg/common"
	"k8s.io/apimachinery/pkg/types"
	"net/http"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"strings"
)

var log = logf.Log.WithName("subscriber")

// Subscriber is a type to manage k8s objects for Subscriber 5G
type Subscriber struct {
	*common.EntityHandlerImpl
	operationRunner common.OperationRunner
	k8sUtils        *common.K8sUtils
	client          client.Client
}

func NewSubscriber(operationRunner common.OperationRunner, k8sUtils *common.K8sUtils, client client.Client) *Subscriber {
	subscriber := Subscriber{
		EntityHandlerImpl: common.NewEntityHandler(),
		operationRunner:   operationRunner,
		k8sUtils:          k8sUtils,
		client:            client,
	}

	subscriber.EntityHandlerImpl.AddTransitionFunction("start", subscriber.Start)

	return &subscriber
}

// IsReady check if ready
func (a Subscriber) IsReady(instance interface{}) bool {
	cr, ok := instance.(*v1alpha1.Subscriber)
	if ok {
		return cr.Status.IsReady()
	} else {
		return false
	}
}

// UpdateStatus called to update CR status
func (a Subscriber) UpdateStatus(opResult controllerutil.OperationResult, instance common.Entity, reconcileErr error) {
	cr := instance.(*v1alpha1.Subscriber)

	if reconcileErr != nil {
		cr.Status.SetError(reconcileErr)
	} else {
		obj, err := a.k8sUtils.GetPod(
			types.NamespacedName{Name: cr.Name, Namespace: cr.Namespace},
		)

		// TODO check err
		if err == nil {
			podReadyCondition := common.GetPodReadyCondition(obj.Status)
			isReconciled := false
			if podReadyCondition != nil && obj.Status.PodIP != "" {
				// the pod is ready, update the Mongo CR
				cr.Status.Outputs.IpAddr = obj.Status.PodIP
				cr.Status.Outputs.PodName = cr.Name
				isReconciled = true
			}

			cr.Status.Update(cr.ObjectMeta, isReconciled, opResult)
		}
	}
}

// Finalize method for Subscriber.  Executed on uninstall
func (a Subscriber) Finalize(request reconcile.Request, instance interface{}) {
	// Nothing to do
}

func (a Subscriber) Start(transition *v1alpha1.Transition) (string, error) {
	subscriberCR, err := a.getCRForTransition(transition)
	if err != nil {
		return "", err
	}

	postBody, _ := json.Marshal(map[string]interface{}{
		"plmnID": subscriberCR.Spec.Config.PlmnID,
		"ueId":   fmt.Sprintf("imsi-%s", subscriberCR.Spec.Config.IMSI),
		"AuthenticationSubscription": map[string]interface{}{
			"authenticationManagementField": "8000",
			"authenticationMethod":          "5G_AKA",
			"milenage": map[string]interface{}{
				"op": map[string]string{
					"encryptionAlgorithm": "0",
					"encryptionKey":       "0",
					"opValue":             "",
				},
			},
			"opc": map[string]interface{}{
				"encryptionAlgorithm": "0",
				"encryptionKey":       "0",
				"opcValue":            "8e27b6af0e692e750f32667a3b14605d",
			},
			"permanentKey": map[string]interface{}{
				"encryptionAlgorithm": "0",
				"encryptionKey":       "0",
				"permanentKeyValue":   "8baf473f2f8fd09487cccbd7097c6862",
			},
			"sequenceNumber": "16f3b3f70fc2",
		},
		"AccessAndMobilitySubscriptionData": map[string]interface{}{
			"gpsis": []string{"msisdn-0900000000"},
			"nssai": map[string]interface{}{
				"defaultSingleNssais": []map[string]interface{}{{
					"sst":       "1",
					"sd":        "010203",
					"isDefault": "true",
				},
					{
						"sst":       "1",
						"sd":        "112233",
						"isDefault": "true",
					}},
				"singleNssais": []map[string]interface{}{},
			},
			"subscribedUeAmbr": map[string]interface{}{
				"downlink": "2 Gbps",
				"uplink":   "1 Gbps",
			},
		},
		"SessionManagementSubscriptionData": []map[string]interface{}{
			{
				"singleNssai": map[string]interface{}{
					"sst": 1,
					"sd":  "010203",
				},
				"dnnConfigurations": map[string]interface{}{
					"internet": map[string]interface{}{
						"sscModes": map[string]interface{}{
							"defaultSscMode": "SSC_MODE_1",
							"allowedSscMode": []string{"SSC_MODE_2", "SSC_MODE_3"},
							"pduSessionTypes": map[string]interface{}{
								"defaultSessionType":  "IPV4",
								"allowedSessionTypes": []string{"IPV4"},
								"sessionAmbr":         map[string]string{"uplink": "200 Mbps", "downlink": "100 Mbps"},
								"5gQosProfile": map[string]interface{}{
									"5qi": 9,
									"arp": map[string]interface{}{
										"priorityLevel": 8,
									},
									"priorityLevel": 8,
								},
								"internet2": map[string]interface{}{
									"sscModes": map[string]interface{}{
										"defaultSscMode":  "SSC_MODE_1",
										"allowedSscModes": []string{"SSC_MODE_2", "SSC_MODE_3"},
										"pduSessionTypes": map[string]interface{}{
											"defaultSessionType":  "IPV4",
											"allowedSessionTypes": []string{"IPV4"},
											"sessionAmbr": map[string]interface{}{
												"uplink":   "200 Mbps",
												"downlink": "100 Mbps",
											},
											"5gQosProfile": map[string]interface{}{
												"5qi":           9,
												"arp":           map[string]interface{}{"priorityLevel": 8},
												"priorityLevel": 8,
											},
										},
									},
								},
							},
						},
					},
				},
			},
			{
				"singleNssai": map[string]interface{}{
					"sst": 1,
					"sd":  "112233",
				},
				"dnnConfigurations": map[string]interface{}{
					"{{ .Dnn }}": map[string]interface{}{
						"sscModes": map[string]interface{}{
							"defaultSscMode":  "SSC_MODE_1",
							"allowedSscModes": []string{"SSC_MODE_2", "SSC_MODE_3"},
						},
						"pduSessionTypes": map[string]interface{}{
							"defaultSessionType":  "IPV4",
							"allowedSessionTypes": []string{"IPV4"},
						},
						"sessionAmbr": map[string]interface{}{
							"uplink":   "200 Mbps",
							"downlink": "100 Mbps",
						},
						"5gQosProfile": map[string]interface{}{
							"5qi": 9,
							"arp": map[string]interface{}{
								"priorityLevel": 8,
							},
							"priorityLevel": 8,
						}},
					"internet2": map[string]interface{}{
						"sscModes": map[string]interface{}{
							"defaultSscMode":  "SSC_MODE_1",
							"allowedSscModes": []string{"SSC_MODE_2", "SSC_MODE_3"},
						},
						"pduSessionTypes": map[string]interface{}{
							"defaultSessionType":  "IPV4",
							"allowedSessionTypes": []string{"IPV4"},
						},
						"sessionAmbr": map[string]interface{}{
							"uplink":   "200 Mbps",
							"downlink": "100 Mbps",
						},
						"5gQosProfile": map[string]interface{}{
							"5qi": 9,
							"arp": map[string]interface{}{
								"priorityLevel": 8,
							},
							"priorityLevel": 8,
						},
					},
				},
			},
		},
		"SmfSelectionSubscriptionData": map[string]interface{}{
			"subscribedSnssaiInfos": map[string]interface{}{
				"01010203": map[string]interface{}{
					"dnnInfos": []map[string]interface{}{
						{"dnn": "{{ .Dnn }}"},
						{"dnn": "internet2"},
					},
				},
			},
		},
		"AmPolicyData": map[string]interface{}{
			"subscCats": []string{"free5gc"},
		},
		"SmPolicyData": map[string]interface{}{
			"smPolicySnssaiData": map[string]interface{}{
				"01010203": map[string]interface{}{
					"snssai": map[string]interface{}{
						"sst": 1,
						"sd":  "010203",
					},
					"smPolicyDnnData": map[string]interface{}{
						"{{ .Dnn }}": map[string]interface{}{
							"dnn": "{{ .Dnn }}",
						},
						"internet2": map[string]interface{}{
							"dnn": "internet2",
						},
					},
				},
				"01112233": map[string]interface{}{
					"snssai": map[string]interface{}{
						"sst": 1,
						"sd":  "112233",
					},
					"smPolicyDnnData": map[string]interface{}{
						"{{ .Dnn }}": map[string]interface{}{
							"dnn": "{{ .Dnn }}",
						},
						"internet2": map[string]interface{}{
							"dnn": "internet2",
						},
					},
				},
			},
		},
		"FlowRules": []interface{}{},
	})
	body := bytes.NewBuffer(postBody)
	resp, err := http.Post(fmt.Sprintf("http://%s:5000/api/subscriber/imsi-%s/%s", subscriberCR.Spec.Config.WebconsoleIPAddress,
		subscriberCR.Spec.Config.IMSI, subscriberCR.Spec.Config.PlmnID), "application/json", body)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return "", fmt.Errorf("Invalid subscribe %v", resp)
	}

	return "", nil
}

func (a Subscriber) CheckReady(transitionCR *fivegv1alpha1.Transition, entity common.Entity) (bool, error) {
	return strings.ToLower(transitionCR.Spec.Config.TransitionName) == "start", nil
}

func (a Subscriber) GetCR(transitionCR *fivegv1alpha1.Transition) (common.Entity, error) {
	return a.k8sUtils.GetSubscriber(types.NamespacedName{
		Namespace: transitionCR.Spec.Config.ResourceNamespace,
		Name:      transitionCR.Spec.Config.ResourceName,
	})
}

func (a Subscriber) IsReconciled(transitionCR *fivegv1alpha1.Transition) (bool, error) {
	cr, err := a.getCRForTransition(transitionCR)
	if err != nil {
		transitionCR.Status.SetError(err)
		return false, err
	}
	if cr == nil {
		err = fmt.Errorf("Cannot find cr")
		transitionCR.Status.SetError(err)
		return false, err
	}
	return cr.Status.IsReconciled(), nil
}

func (a Subscriber) getCRForTransition(transition *v1alpha1.Transition) (*v1alpha1.Subscriber, error) {
	cr, err := a.k8sUtils.GetSubscriber(types.NamespacedName{
		Namespace: transition.Spec.Config.ResourceNamespace,
		Name:      transition.Spec.Config.ResourceName,
	})
	if err != nil {
		return nil, err
	}

	return cr, nil
}
