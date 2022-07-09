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

package v13

import (
	"bytes"
	v1alpha1 "github.ibm.com/Steve-Glover/5GOperator/api/v1alpha1"
	"text/template"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const PcfConfig = `info:
  version: 1.0.0
  description: PCF initial local configuration

configuration:
  pcfName: PCF # the name of this PCF
  sbi: # Service-based interface information
    scheme: http # the protocol for sbi (http or https)
    registerIPv4: {{ .Name }}-sbi # IP used to register to NRF
    bindingIPv4: {{ .Name }}-sbi  # IP used to bind the service
    port: 8000              # port used to bind the service
  timeFormat: 2019-01-02 15:04:05 # time format of this PCF
  defaultBdtRefId: BdtPolicyId-   # BDT Reference ID, indicating transfer policies of background data transfer.
  nrfUri: http://{{ .NrfIPAddress }}:{{ .NrfPort }}  # a valid URI of NRF
  serviceList:   # the SBI services provided by this PCF, refer to TS 29.507
    - serviceName: npcf-am-policy-control # Npcf_AMPolicyControl service
    - serviceName: npcf-smpolicycontrol   # Npcf_SMPolicyControl service
      suppFeat: 3fff # the features supported by Npcf_SMPolicyControl, name defined in TS 29.512 5.8-1, value defined in TS 29.571 5.2.2
    - serviceName: npcf-bdtpolicycontrol    # Npcf_BDTPolicyControl service
    - serviceName: npcf-policyauthorization # Npcf_PolicyAuthorization service
      suppFeat: 3    # the features supported by Npcf_PolicyAuthorization, name defined in TS 29.514 5.8-1, value defined in TS 29.571 5.2.2
    - serviceName: npcf-eventexposure       # Npcf_EventExposure service
    - serviceName: npcf-ue-policy-control   # Npcf_UEPolicyControl service
  mongodb:       # the mongodb connected by this PCF
    name: free5gc                  # name of the mongodb
    url: mongodb://{{ .MongoIPAddress }}:27017 # a valid URL of the mongodb

# the kind of log output
  # debugLevel: how detailed to output, value: trace, debug, info, warn, error, fatal, panic
  # ReportCaller: enable the caller report or not, value: true or false
logger:
  PCF:
    debugLevel: info
    ReportCaller: false
  PathUtil:
    debugLevel: info
    ReportCaller: false
  OpenApi:
    debugLevel: info
    ReportCaller: false
`

func PcfConfigMap(cr *v1alpha1.Pcf) (*corev1.ConfigMap, error) {
	t, err := template.New("pcf").Parse(PcfConfig)
	if err != nil {
		return nil, err
	}

	var b bytes.Buffer
	err = t.Execute(&b, cr.Spec.Config)
	if err != nil {
		return nil, err
	}

	return &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cr.Name,
			Namespace: cr.Namespace,
			Labels: map[string]string{
				"app": "pcf",
			},
		},
		Data: map[string]string{
			"pcf.yaml": b.String(),
		},
	}, nil
}
