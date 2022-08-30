#!/usr/bin/env bash

SCRIPT_DIR="$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )"

NAMESPACE=${1:-5g-test}
URL=${2:-http://127.0.0.1:31527}

kubectl create namespace $NAMESPACE

export BETWEEN=30

# RAN
kubectl apply -f $SCRIPT_DIR/5g_v1alpha1_upf-r1.yaml -n $NAMESPACE

# Transport
kubectl apply -f $SCRIPT_DIR/5g_v1alpha1_upf-t1.yaml -n $NAMESPACE
kubectl apply -f $SCRIPT_DIR/5g_v1alpha1_upf-t2.yaml -n $NAMESPACE

# Core
kubectl apply -f $SCRIPT_DIR/5g_v1alpha1_upf-c1.yaml -n $NAMESPACE
kubectl apply -f $SCRIPT_DIR/5g_v1alpha1_upf-c2.yaml -n $NAMESPACE
kubectl apply -f $SCRIPT_DIR/5g_v1alpha1_upf-c3.yaml -n $NAMESPACE
kubectl apply -f $SCRIPT_DIR/5g_v1alpha1_upf-c4.yaml -n $NAMESPACE

sleep $BETWEEN

echo ""
echo ""
echo "-=-=-=-=-=-= TRACE -=-=-=-=-=-=-=-=-=-"
echo "Register upf-r1.."
echo "-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=--=-=-=-"
curl -H "Content-type: application/json" -X POST -d "@$SCRIPT_DIR/upf-r1.json" $URL/upi/v1/upf

echo ""
echo ""
echo "-=-=-=-=-=-= TRACE -=-=-=-=-=-=-=-=-=-"
echo "Register upf-t1.."
echo "-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=--=-=-=-"
curl -H "Content-type: application/json" -X POST -d "@$SCRIPT_DIR/upf-t1.json" $URL/upi/v1/upf

echo ""
echo ""
echo "-=-=-=-=-=-= TRACE -=-=-=-=-=-=-=-=-=-"
echo "Register upf-t2.."
echo "-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=--=-=-=-"
curl -H "Content-type: application/json" -X POST -d "@$SCRIPT_DIR/upf-t2.json" $URL/upi/v1/upf

echo ""
echo ""
echo "-=-=-=-=-=-= TRACE -=-=-=-=-=-=-=-=-=-"
echo "Register upf-c1.."
echo "-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=--=-=-=-"
curl -H "Content-type: application/json" -X POST -d "@$SCRIPT_DIR/upf-c1.json" $URL/upi/v1/upf

echo ""
echo ""
echo "-=-=-=-=-=-= TRACE -=-=-=-=-=-=-=-=-=-"
echo "Register upf-c2.."
echo "-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=--=-=-=-"
curl -H "Content-type: application/json" -X POST -d "@$SCRIPT_DIR/upf-c2.json" $URL/upi/v1/upf

echo ""
echo ""
echo "-=-=-=-=-=-= TRACE -=-=-=-=-=-=-=-=-=-"
echo "Register upf-c3.."
echo "-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=--=-=-=-"
curl -H "Content-type: application/json" -X POST -d "@$SCRIPT_DIR/upf-c3.json" $URL/upi/v1/upf

echo ""
echo ""
echo "-=-=-=-=-=-= TRACE -=-=-=-=-=-=-=-=-=-"
echo "Register upf-c41.."
echo "-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=--=-=-=-"
curl -H "Content-type: application/json" -X POST -d "@$SCRIPT_DIR/upf-c4.json" $URL/upi/v1/upf
