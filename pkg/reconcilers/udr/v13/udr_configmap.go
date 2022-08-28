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
	"github.com/5GZORRO/issm-mec-cnmp-5g-operator/api/v1alpha1"
	"text/template"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const UdrConfig = `info:
 version: 1.0.1
 description: UDR initial local configuration

configuration:
 sbi: # Service-based interface information
   scheme: http # the protocol for sbi (http or https)
   registerIPv4: {{ .Name }}-sbi # IP used to register to NRF
   bindingIPv4: {{ .Name }}-sbi  # IP used to bind the service
   port: 8000 # port used to bind the service
    tls: # the local path of TLS key
      pem: ./config/TLS/udr.pem # UDR TLS Certificate
      key: ./config/TLS/udr.key # UDR TLS Private key
 mongodb:
   name: free5gc # Database name in MongoDB
   url: mongodb://{{ .MongoIPAddress }}:27017 # URL of MongoDB
 nrfUri: http://{{ .NrfIPAddress }}:{{ .NrfPort }} # a valid URI of NRF

# the kind of log output
 # debugLevel: how detailed to output, value: trace, debug, info, warn, error, fatal, panic
 # ReportCaller: enable the caller report or not, value: true or false
logger:
 UDR:
   debugLevel: debug
   ReportCaller: false
 MongoDBLibrary:
   debugLevel: info
   ReportCaller: false
 PathUtil:
   debugLevel: info
   ReportCaller: false
 OpenApi:
   debugLevel: info
   ReportCaller: false
`

func UdrConfigMap(cr *v1alpha1.Udr) (*corev1.ConfigMap, error) {
	t, err := template.New("udr").Parse(UdrConfig)
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
				"app": "udr",
			},
		},
		Data: map[string]string{
			"udr.yaml": b.String(),
		},
	}, nil
}
