package v13

import (
	"bytes"
	"github.com/5GZORRO/issm-mec-cnmp-5g-operator/api/v1alpha1"
	"text/template"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const SmfConfig = `info:
  version: 1.0.2
  description: SMF initial local configuration

configuration:
  smfName: SMF # the name of this SMF
  sbi: # Service-based interface information
    scheme: http # the protocol for sbi (http or https)
    registerIPv4: {{ .Name }}-sbi # IP used to register to NRF
    bindingIPv4: 0.0.0.0  # IP used to bind the service
    port: 8000 # Port used to bind the service
    tls: # the local path of TLS key
      key: ./config/TLS/smf.key # SMF TLS Certificate
      pem: ./config/TLS/smf.pem # SMF TLS Private key
  serviceNameList: # the SBI services provided by this SMF, refer to TS 29.502
    - nsmf-pdusession # Nsmf_PDUSession service
    - nsmf-event-exposure # Nsmf_EventExposure service
    - nsmf-oam # OAM service
  snssaiInfos: # the S-NSSAI (Single Network Slice Selection Assistance Information) list supported by this AMF
    - sNssai: # S-NSSAI (Single Network Slice Selection Assistance Information)
        sst: 1 # Slice/Service Type (uinteger, range: 0~255)
        sd: 010203 # Slice Differentiator (3 bytes hex string, range: 000000~FFFFFF)
      dnnInfos: # DNN information list
        - dnn: internet # Data Network Name
          dns: # the IP address of DNS
            ipv4: 8.8.8.8
    - sNssai: # S-NSSAI (Single Network Slice Selection Assistance Information)
        sst: 1 # Slice/Service Type (uinteger, range: 0~255)
        sd: 112233 # Slice Differentiator (3 bytes hex string, range: 000000~FFFFFF)
      dnnInfos: # DNN information list
        - dnn: internet # Data Network Name
          dns: # the IP address of DNS
            ipv4: 8.8.8.8
  plmnList:
    - mcc: "208"
      mnc: "93"
  pfcp: # the IP address of N4 interface on this SMF (PFCP)
    addr: {{ .Name }}-sbi
  ulcl: true
{{- if .Nodes }}
  userplaneInformation: # list of userplane information
    upNodes: # information of userplane node (AN or UPF)
{{- range .Nodes }}
      {{ .Name }}:
        type: {{ .Type }}
{{- if .AnIp }}
        anIp: {{ .AnIp }}
{{- end}}
{{- if .NrCellId }}
        nrCellId: {{ .NrCellId }}
{{- end}}
{{- if eq .Type "UPF" }}
        nodeId: {{ .NodeIdSbi }} # the IP/FQDN of N4 interface on this UPF (PFCP)
        sNssaiUpfInfos: # S-NSSAI information list for this UPF
          - sNssai: # S-NSSAI (Single Network Slice Selection Assistance Information)
              sst: {{ .Sst }} # Slice/Service Type (uinteger, range: 0~255)
              sd: {{ .Sd }} # Slice Differentiator (3 bytes hex string, range: 000000~FFFFFF)
            dnnUpfInfoList: # DNN information list for this S-NSSAI
              - dnn: internet
{{- if .Pool }}
                pools:
                  - cidr: {{ .Pool }}
{{- end}}
        interfaces: # Interface list for this UPF
          - interfaceType: N3 # the type of the interface (N3 or N9)
            endpoints: # the IP address of this N3/N9 interface on this UPF
              - {{ .NodeIdUp }}
            networkInstance: internet # Data Network Name (DNN)
          - interfaceType: N9 # the type of the interface (N3 or N9)
            endpoints: # the IP address of this N3/N9 interface on this UPF
              - {{ .NodeIdUp }}
            networkInstance: internet # Data Network Name (DNN)
{{- end}}
{{- end}}
{{- end}}
    links: null
  dnn:
    internet:
      dns:
        ipv4: 8.8.8.8
        ipv6: 2001:4860:4860::8888
    internet2:
      dns:
        ipv4: 8.8.4.4
        ipv6: 2001:4860:4860::8844
  nrfUri: http://{{ .NrfIPAddress }}:{{ .NrfPort }} # a valid URI of NRF
  locality: area1
  smfExtUri: http://127.0.0.1:8080

# the kind of log output
  # debugLevel: how detailed to output, value: trace, debug, info, warn, error, fatal, panic
  # ReportCaller: enable the caller report or not, value: true or false
logger:
  SMF:
    debugLevel: trace
    ReportCaller: false
  NAS:
    debugLevel: trace
    ReportCaller: false
  NGAP:
    debugLevel: trace
    ReportCaller: false
  Aper:
    debugLevel: info
    ReportCaller: false
  PathUtil:
    debugLevel: debug
    ReportCaller: false
  OpenApi:
    debugLevel: debug
    ReportCaller: false
  PFCP:
    debugLevel: trace
    ReportCaller: false
  PduSess:
    debugLevel: trace
    ReportCaller: false
  CTX:
    debugLevel: trace
    ReportCaller: false
`

const EuRoutingConfig = `info:
  version: 1.0.1
  description: Routing information for UE
`

func smfConfig(cfg *v1alpha1.SmfConfig) (string, error) {
	t, err := template.New("smf").Parse(SmfConfig)
	if err != nil {
		return "", err
	}

	var b bytes.Buffer
	err = t.Execute(&b, cfg)
	if err != nil {
		return "", err
	}

	return b.String(), nil
}

func euRoutingConfig(cfg *v1alpha1.SmfConfig) (string, error) {
	t, err := template.New("smf-eurouting").Parse(EuRoutingConfig)
	if err != nil {
		return "", err
	}

	var b bytes.Buffer
	err = t.Execute(&b, cfg)
	if err != nil {
		return "", err
	}

	return b.String(), nil
}

func SmfConfigMap(cr *v1alpha1.Smf) (*corev1.ConfigMap, error) {
	smfCfg, err := smfConfig(cr.Spec.Config)
	if err != nil {
		return nil, err
	}

	euRoutingCfg, err := euRoutingConfig(cr.Spec.Config)
	if err != nil {
		return nil, err
	}

	return &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cr.Name,
			Namespace: cr.Namespace,
			Labels: map[string]string{
				"app": "smf",
			},
		},
		Data: map[string]string{
			"smf.yaml":       smfCfg,
			"uerouting.yaml": euRoutingCfg,
		},
	}, nil
}
