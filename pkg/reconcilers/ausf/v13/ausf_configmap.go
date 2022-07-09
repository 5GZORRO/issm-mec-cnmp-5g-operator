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

const AusfConfig = `info:
  version: 1.0.0
  description: AUSF initial local configuration

configuration:
  sbi: # Service-based interface information
    scheme: http # the protocol for sbi (http or https)
    registerIPv4: {{ .Name }}-sbi # IP used to register to NRF
    bindingIPv4: {{ .Name }}-sbi  # IP used to bind the service
    port: 8000 # Port used to bind the service
  serviceNameList: # the SBI services provided by this AUSF, refer to TS 29.509
    - nausf-auth # Nausf_UEAuthentication service
  nrfUri: http://{{ .NrfIPAddress }}:{{ .NrfPort }} # a valid URI of NRF
  plmnSupportList: # the PLMNs (Public Land Mobile Network) list supported by this AUSF
    - mcc: 208 # Mobile Country Code (3 digits string, digit: 0~9)
      mnc: 93  # Mobile Network Code (2 or 3 digits string, digit: 0~9)
    - mcc: 123 # Mobile Country Code (3 digits string, digit: 0~9)
      mnc: 45  # Mobile Network Code (2 or 3 digits string, digit: 0~9)
  groupId: ausfGroup001 # ID for the group of the AUSF

# the kind of log output
  # debugLevel: how detailed to output, value: trace, debug, info, warn, error, fatal, panic
  # ReportCaller: enable the caller report or not, value: true or false
logger:
  AUSF:
    debugLevel: debug
    ReportCaller: false
  PathUtil:
    debugLevel: info
    ReportCaller: false
  OpenApi:
    debugLevel: info
    ReportCaller: false
`

func AusfConfigMap(cr *v1alpha1.Ausf) (*corev1.ConfigMap, error) {
	t, err := template.New("ausf").Parse(AusfConfig)
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
				"app": "ausf",
			},
		},
		Data: map[string]string{
			"ausf.yaml": b.String(),
		},
	}, nil
}
