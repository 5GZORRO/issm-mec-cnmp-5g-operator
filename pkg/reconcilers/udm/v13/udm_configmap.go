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
	v1alpha1 "github.com/5GZORRO/issm-mec-cnmp-5g-operator/api/v1alpha1"
	"text/template"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const UdmConfig = `info:
  version: 1.0.2
  description: UDM initial local configuration

configuration:
  serviceNameList: # the SBI services provided by this UDM, refer to TS 29.503
    - nudm-sdm # Nudm_SubscriberDataManagement service
    - nudm-uecm # Nudm_UEContextManagement service
    - nudm-ueau # Nudm_UEAuthenticationManagement service
    - nudm-ee # Nudm_EventExposureManagement service
    - nudm-pp # Nudm_ParameterProvisionDataManagement service
  sbi: # Service-based interface information
    scheme: http # the protocol for sbi (http or https)
    registerIPv4: {{ .Name }}-sbi # IP used to register to NRF
    bindingIPv4: {{ .Name }}-sbi # IP used to bind the service
    port: 8000 # Port used to bind the service
    tls: # the local path of TLS key
      pem: ./config/TLS/udm.pem # UDM TLS Certificate
      key: ./config/TLS/udm.key # UDM TLS Private key
  nrfUri: http://{{ .NrfIPAddress }}:8000 # a valid URI of NRF

  # test data set from TS33501-f60 Annex C.4
  SuciProfile:
    - ProtectionScheme: 1 # Protect Scheme: Profile A
      PrivateKey: c53c22208b61860b06c62e5406a7b330c2b577aa5558981510d128247d38bd1d
      PublicKey: 5a8d38864820197c3394b92613b20b91633cbd897119273bf8e4a6f4eec0a650
    - ProtectionScheme: 2 # Protect Scheme: Profile B
      PrivateKey: F1AB1074477EBCC7F554EA1C5FC368B1616730155E0041AC447D6301975FECDA
      PublicKey: 0472DA71976234CE833A6907425867B82E074D44EF907DFB4B3E21C1C2256EBCD15A7DED52FCBB097A4ED250E036C7B9C8C7004C4EEDC4F068CD7BF8D3F900E3B4

# the kind of log output
  # debugLevel: how detailed to output, value: trace, debug, info, warn, error, fatal, panic
  # ReportCaller: enable the caller report or not, value: true or false
logger:
  UDM:
    debugLevel: debug
    ReportCaller: false
  OpenApi:
    debugLevel: info
    ReportCaller: false
  PathUtil:
    debugLevel: info
    ReportCaller: false
`

func UdmConfigMap(cr *v1alpha1.Udm) (*corev1.ConfigMap, error) {
	t, err := template.New("udm").Parse(UdmConfig)
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
			// TODO use cm.Name here
			Name:      cr.Name,
			Namespace: cr.Namespace,
			Labels: map[string]string{
				"app": "udm",
			},
		},
		Data: map[string]string{
			"udm.yaml": b.String(),
		},
	}, nil
}
