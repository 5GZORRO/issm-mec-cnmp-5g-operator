#! /usr/bin/env bash

NAMESPACE=${1:-5g-core}

kubectl delete -f ./5g_v1alpha1_upf1.yaml -n $NAMESPACE

kubectl apply -f ./smf-sample-rmnode-upf1.yaml -n $NAMESPACE

kubectl apply -f ./smf-sample-rmlink-upf1.yaml -n $NAMESPACE

