#! /usr/bin/env bash

NAMESPACE=${1:-5g-core}

kubectl delete -f ./5g_v1alpha1_upf2.yaml -n $NAMESPACE

kubectl apply -f ./smf-sample-rmnode-upf2.yaml -n $NAMESPACE

kubectl apply -f ./smf-sample-rmlink-upf2.yaml -n $NAMESPACE

