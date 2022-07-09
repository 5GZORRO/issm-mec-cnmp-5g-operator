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

const N3iwfConfig = `info:
  version: 1.0.0
  description: N3IWF initial local configuration

configuration:
  N3IWFInformation:
    GlobalN3IWFID:
      PLMNID:
        MCC: {{ .Mcc }}
        MNC: {{ .Mnc }}
      N3IWFID: 135
    Name:
      free5GC_N3IWF
    SupportedTAList:
      - TAC: 1
        BroadcastPLMNList:
          - PLMNID:
              MCC: {{ .Mcc }}
              MNC: {{ .Mnc }}
            TAISliceSupportList:
              - SNSSAI:
                  SST: 1
                  SD: 010203
              - SNSSAI:
                  SST: 1
                  SD: 112233
  AMFSCTPAddresses:
    - IP: 
      - {{ .AmfIPAddress }}

  # Bind Interfaces
  # IKE interface
  # TODO will 0.0.0.0 work?
  IKEBindAddress: 0.0.0.0
  # IPSec virtual interface
  IPSecInterfaceAddress: {{ .IPSecAddress }}
  # IPSec virtual interface mark
  IPSecInterfaceMark: 5
  # GTP interface
  # TODO will 0.0.0.0 work?
  GTPBindAddress: 0.0.0.0

  # NAS TCP Listen Port
  NASTCPPort: 20000

  # N3IWF FQDN
  FQDN: n3iwf.free5gc.org

  # Security
  # Private Key File Path
  PrivateKey:
  # Certificate Authority (CA)
  CertificateAuthority:
  # Certificate
  Certificate:

  # IP address that will be allocated to UE in IPSec tunnel
  UEIPAddressRange: {{ .UECIDR }}
`

func N3iwfConfigMap(cr *v1alpha1.N3iwf) (*corev1.ConfigMap, error) {
	t, err := template.New("n3iwf").Parse(N3iwfConfig)
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
				"app": "n3iwf",
			},
		},
		Data: map[string]string{
			"free5GC.conf": b.String(),
		},
	}, nil
}

const Free5GCConfig = `#all logging levels
#panic
#fatal
#error
#warn
#info
#debug
#trace
logger:
# network function
  AMF:
    debugLevel: info
    ReportCaller: true
  SMF:
    debugLevel: info
    ReportCaller: true
  UDR:
    debugLevel: info
    ReportCaller: true
  UDM:
    debugLevel: info
    ReportCaller: true
  NRF:
    debugLevel: info
    ReportCaller: true
  PCF:
    debugLevel: info
    ReportCaller: true
  AUSF:
    debugLevel: info
    ReportCaller: true
  N3IWF:
    debugLevel: info
    ReportCaller: true
# library
  NAS:
    debugLevel: info
    ReportCaller: true
  FSM:
    debugLevel: info
    ReportCaller: true
  NGAP:
    debugLevel: info
    ReportCaller: true
  NamfComm:
    debugLevel: info
    ReportCaller: true
  NamfEventExposure:
    debugLevel: info
    ReportCaller: true
  NsmfPDUSession:
    debugLevel: info
    ReportCaller: true
  NudrDataRepository:
    debugLevel: info
    ReportCaller: true
  OpenApi:
    debugLevel: info
    ReportCaller: true
  Aper:
    debugLevel: info
    ReportCaller: true
  CommonConsumerTest:
    debugLevel: info
    ReportCaller: true
# webui
  WEBUI:
    debugLevel: info
    ReportCaller: true
`

func Free5gCConfigMap(cr *v1alpha1.N3iwf) (*corev1.ConfigMap, error) {
	t := template.Must(template.New("n3iwf-free5gC").Parse(Free5GCConfig))

	var b bytes.Buffer
	err := t.Execute(&b, cr.Spec.Config)
	if err != nil {
		return nil, err
	}

	return &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "n3iwf-free5gC",
			Namespace: cr.Namespace,
			Labels: map[string]string{
				"app": "n3iwf",
			},
		},
		Data: map[string]string{
			"free5GC.conf": b.String(),
		},
	}, nil
}
