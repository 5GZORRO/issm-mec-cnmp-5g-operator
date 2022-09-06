#! /usr/bin/env bash

NAMESPACE=${1:-5g-test}

export BETWEEN=10

kubectl create namespace $NAMESPACE

kubectl apply -f 5g-core-admin.yaml -n $NAMESPACE

# 
# It is important to deploy NFs in the following order.
#
# DB > NRF > UDR > UDM > AUSF > NSSF > AMF > PCF > SMF > N3IWF
#

kubectl apply -f ./networks.yaml -n $NAMESPACE

kubectl apply -f ./5g_v1alpha1_mongo.yaml -n $NAMESPACE

sleep $BETWEEN

kubectl apply -f ./5g_v1alpha1_nrf.yaml -n $NAMESPACE

sleep $BETWEEN

kubectl apply -f ./5g_v1alpha1_udr.yaml -n $NAMESPACE

sleep $BETWEEN

kubectl apply -f ./5g_v1alpha1_udm.yaml -n $NAMESPACE

sleep $BETWEEN

kubectl apply -f ./5g_v1alpha1_ausf.yaml -n $NAMESPACE

sleep $BETWEEN

kubectl apply -f ./5g_v1alpha1_nssf.yaml -n $NAMESPACE

sleep $BETWEEN

kubectl apply -f ./5g_v1alpha1_amf.yaml -n $NAMESPACE

sleep $BETWEEN

kubectl apply -f ./5g_v1alpha1_pcf.yaml -n $NAMESPACE

sleep $BETWEEN

kubectl apply -f ./5g_v1alpha1_smf.yaml -n $NAMESPACE

sleep $BETWEEN

kubectl apply -f ./5g_v1alpha1_webconsole.yaml -n $NAMESPACE

sleep $BETWEEN

curl -X POST http://127.0.0.1:30050/api/subscriber/imsi-208930000000003/20893 -H "Token: admin" -d @sub.json
