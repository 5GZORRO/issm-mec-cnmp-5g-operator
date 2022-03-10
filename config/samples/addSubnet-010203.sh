#! /usr/bin/env bash

NAMESPACE=${1:-5g-test}

kubectl create namespace $NAMESPACE

export BETWEEN=5

kubectl apply -f ./5g_v1alpha1_upf1.yaml -n $NAMESPACE

sleep $BETWEEN

kubectl apply -f ./smf-sample-addnode-upf1.yaml -n $NAMESPACE

kubectl apply -f ./smf-sample-addlink-upf1.yaml -n $NAMESPACE

