#!/bin/bash

URL=${1:-http://127.0.0.1:38080}

curl -X PUT \
  $URL/links \
  -H 'content-type: application/json' \
  -d '{"A": "gNB1", "B": "upf-c1-sample"}'

curl -X PUT \
  $URL/links \
  -H 'content-type: application/json' \
  -d '{"A": "gNB2", "B": "upf-c2-sample"}'

curl -X GET \
  $URL/links \
  -H 'content-type: application/json'
