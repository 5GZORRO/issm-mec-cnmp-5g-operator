DOCKER_IMAGE_OWNER = 'free5gc'
DOCKER_IMAGE_NAME = 'base'
DOCKER_IMAGE_TAG = 'latest'

.PHONY: base
base:
	docker build -t ${DOCKER_IMAGE_OWNER}/${DOCKER_IMAGE_NAME}:${DOCKER_IMAGE_TAG} ./base
	docker build -t ${DOCKER_IMAGE_OWNER}/nf_amf:${DOCKER_IMAGE_TAG} ./nf_amf
	docker build -t ${DOCKER_IMAGE_OWNER}/nf_ausf:${DOCKER_IMAGE_TAG} ./nf_ausf
	docker build -t ${DOCKER_IMAGE_OWNER}/nf_n3iwf:${DOCKER_IMAGE_TAG} ./nf_n3iwf
	docker build -t ${DOCKER_IMAGE_OWNER}/nf_nrf:${DOCKER_IMAGE_TAG} ./nf_nrf
	docker build -t ${DOCKER_IMAGE_OWNER}/nf_nssf:${DOCKER_IMAGE_TAG} ./nf_nssf
	docker build -t ${DOCKER_IMAGE_OWNER}/nf_pcf:${DOCKER_IMAGE_TAG} ./nf_pcf
	docker build -t ${DOCKER_IMAGE_OWNER}/nf_smf:${DOCKER_IMAGE_TAG} ./nf_smf
	docker build -t ${DOCKER_IMAGE_OWNER}/nf_udm:${DOCKER_IMAGE_TAG} ./nf_udm
	docker build -t ${DOCKER_IMAGE_OWNER}/nf_udr:${DOCKER_IMAGE_TAG} ./nf_udr
	docker build -t ${DOCKER_IMAGE_OWNER}/nf_upf:${DOCKER_IMAGE_TAG} ./nf_upf
	docker build -t ${DOCKER_IMAGE_OWNER}/webui:${DOCKER_IMAGE_TAG} ./webui
	docker build -t ${DOCKER_IMAGE_OWNER}/nettools:${DOCKER_IMAGE_TAG} ./nettools
	docker build -t ${DOCKER_IMAGE_OWNER}/grafana:${DOCKER_IMAGE_TAG} ./grafana
	docker build -t ${DOCKER_IMAGE_OWNER}/gnb-gw:${DOCKER_IMAGE_TAG} ./gnb-gw
	docker build -t ${DOCKER_IMAGE_OWNER}/prometheus:${DOCKER_IMAGE_TAG} ./prometheus
	docker pull mongo
	docker pull nginxdemos/hello
