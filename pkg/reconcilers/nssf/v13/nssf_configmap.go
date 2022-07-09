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

const NssfConfig = `info:
  version: 1.0.0
  description: NSSF initial local configuration

configuration:
  nssfName: NSSF # the name of this NSSF
  sbi: # Service-based interface information
    scheme: http # the protocol for sbi (http or https)
    registerIPv4: {{ .Name }}-sbi # IP used to register to NRF
    bindingIPv4: {{ .Name }}-sbi  # IP used to bind the service
    port: 8000 # Port used to bind the service
  serviceNameList: # the SBI services provided by this SMF, refer to TS 29.531
    - nnssf-nsselection # Nnssf_NSSelection service
    - nnssf-nssaiavailability # Nnssf_NSSAIAvailability service
  nrfUri: http://{{ .NrfIPAddress }}:{{ .NrfPort }} # a valid URI of NRF
  supportedPlmnList: # the PLMNs (Public land mobile network) list supported by this NSSF
    - mcc: {{ .Mcc }} # Mobile Country Code (3 digits string, digit: 0~9)
      mnc: {{ .Mnc }} # Mobile Network Code (2 or 3 digits string, digit: 0~9)
  supportedNssaiInPlmnList: # Supported S-NSSAI List for each PLMN
    - plmnId: # Public Land Mobile Network ID, <PLMN ID> = <MCC><MNC>
        mcc: {{ .Mcc }} # Mobile Country Code (3 digits string, digit: 0~9)
        mnc: {{ .Mnc }} # Mobile Network Code (2 or 3 digits string, digit: 0~9)
      supportedSnssaiList: # Supported S-NSSAIs of the PLMN
        - sst: 1 # Slice/Service Type (uinteger, range: 0~255)
          sd: 010203 # Slice Differentiator (3 bytes hex string, range: 000000~FFFFFF)
        - sst: 1 # Slice/Service Type (uinteger, range: 0~255)
          sd: 112233 # Slice Differentiator (3 bytes hex string, range: 000000~FFFFFF)
  nsiList: # List of available Network Slice Instance (NSI)
    - snssai: # S-NSSAI of this NSI
        sst: 1 # Slice/Service Type (uinteger, range: 0~255)
      nsiInformationList: # Information list of this NSI
        # the NRF to be used to select the NFs/services within the selected NSI, and an optonal ID
        - nrfId: http://{{ .NrfIPAddress }}:{{ .NrfPort }}/nnrf-nfm/v1/nf-instances
          nsiId: 10
    - snssai: # S-NSSAI of this NSI
        sst: 1 # Slice/Service Type (uinteger, range: 0~255)
        sd: 010203 # Slice Differentiator (3 bytes hex string, range: 000000~FFFFFF)
      nsiInformationList: # Information list of this NSI
        # the NRF to be used to select the NFs/services within the selected NSI, and an optonal ID
        - nrfId: http://{{ .NrfIPAddress }}:{{ .NrfPort }}/nnrf-nfm/v1/nf-instances
          nsiId: 22
    - snssai: # S-NSSAI of this NSI
        sst: 1 # Slice/Service Type (uinteger, range: 0~255)
        sd: 112233 # Slice Differentiator (3 bytes hex string, range: 000000~FFFFFF)
      nsiInformationList: # Information list of this NSI
        # the NRF to be used to select the NFs/services within the selected NSI, and an optonal ID
        - nrfId: http://{{ .NrfIPAddress }}:{{ .NrfPort }}/nnrf-nfm/v1/nf-instances
          nsiId: 23

# the kind of log output
  # debugLevel: how detailed to output, value: trace, debug, info, warn, error, fatal, panic
  # ReportCaller: enable the caller report or not, value: true or false
logger:
  NSSF:
    debugLevel: info
    ReportCaller: false
  PathUtil:
    debugLevel: info
    ReportCaller: false
  OpenApi:
    debugLevel: info
    ReportCaller: false
`

func NssfConfigMap(cr *v1alpha1.Nssf) (*corev1.ConfigMap, error) {
	t, err := template.New("nssf").Parse(NssfConfig)
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
				"app": "nssf",
			},
		},
		Data: map[string]string{
			"nssf.yaml": b.String(),
		},
	}, nil
}
