#! /usr/bin/env bash

NAMESPACE=${1:-5g-core}

kubectl create namespace $NAMESPACE

export BETWEEN=5

kubectl apply -f ./datanetwork.yaml -n $NAMESPACE

kubectl apply -f ./5g_v1alpha1_upf2.yaml -n $NAMESPACE

sleep $BETWEEN

kubectl apply -f ./smf-sample-addnode-upf2.yaml -n $NAMESPACE

kubectl apply -f ./smf-sample-addlink-upf2.yaml -n $NAMESPACE
