#! /usr/bin/env bash

NAMESPACE=${1:-5g-test}

export BETWEEN=3

echo ""
echo ""
echo "-=-=-=-=-=-= TRACE -=-=-=-=-=-=-=-=-=-"
echo "Stop NFs.. & networks.."
echo "-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=--=-=-=-"

kubectl  delete vCache  --all -n $NAMESPACE

kubectl delete -f . -n $NAMESPACE

kubectl delete UPF -n $NAMESPACE --all