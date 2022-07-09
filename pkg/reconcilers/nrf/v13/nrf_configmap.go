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

const NrfConfig = `info:
  version: 1.0.0
  description: NRF initial local configuration

configuration:
  MongoDBName: free5gc # database name in MongoDB
  MongoDBUrl: mongodb://{{ .MongoIPAddr }}:27017 # a valid URL of the mongodb
  sbi: # Service-based interface information
    registerIPv4: {{ .Name }}-sbi # IP used to serve NFs or register to another NRF
    bindingIPv4: {{ .Name }}-sbi # IP used to bind the service
    scheme: http # the protocol for sbi (http or https)
    port: 8000 # port used to bind the service
  DefaultPlmnId:
    mcc: {{ .Mcc }} # Mobile Country Code (3 digits string, digit: 0~9)
    mnc: {{ .Mnc }} # Mobile Network Code (2 or 3 digits string, digit: 0~9)
  serviceNameList: # the SBI services provided by this NRF, refer to TS 29.510
    - nnrf-nfm # Nnrf_NFManagement service
    - nnrf-disc # Nnrf_NFDiscovery service

# the kind of log output
  # debugLevel: how detailed to output, value: trace, debug, info, warn, error, fatal, panic
  # ReportCaller: enable the caller report or not, value: true or false
logger:
  NRF:
    debugLevel: debug
    ReportCaller: false
  PathUtil:
    debugLevel: debug
    ReportCaller: false
  OpenApi:
    debugLevel: debug
    ReportCaller: false
  MongoDBLibrary:
    debugLevel: info
    ReportCaller: false
`

func NrfConfigMap(cr *v1alpha1.Nrf) (*corev1.ConfigMap, error) {
	t, err := template.New("nrf").Parse(NrfConfig)
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
				"app": "nrf",
			},
		},
		Data: map[string]string{
			"nrf.yaml": b.String(),
		},
	}, nil
}
