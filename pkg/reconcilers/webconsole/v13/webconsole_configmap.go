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

const WebconsoleConfig = `info:
  version: 1.0.0
  description: WebUI initial local configuration

configuration:
  mongodb: # the mongodb connected by this webui
    name: free5gc # name of the mongodb
    url: mongodb://{{ .MongoIPAddr }}:27017 # a valid URL of the mongodb
`

func WebconsoleConfigMap(cr *v1alpha1.Webconsole) (*corev1.ConfigMap, error) {
	t, err := template.New("webconsole").Parse(WebconsoleConfig)
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
				"app": "webconsole",
			},
		},
		Data: map[string]string{
			"webuicfg.yaml": b.String(),
		},
	}, nil
}
