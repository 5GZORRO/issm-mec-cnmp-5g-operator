#!/bin/bash

URL=${1:-http://127.0.0.1:38080}

curl -X POST \
  $URL/ue-routes/red

curl -X POST \
  $URL/ue-routes/blue

curl -X POST \
  $URL/ue-routes/red/members/imsi-208930000000001
curl -X POST \
  $URL/ue-routes/red/members/imsi-208930000000002
curl -X POST \
  $URL/ue-routes/red/members/imsi-208930000000003
curl -X POST \
  $URL/ue-routes/red/members/imsi-208930000000004
curl -X POST \
  $URL/ue-routes/red/members/imsi-208930000000005


curl -X POST \
  $URL/ue-routes/blue/members/imsi-208930000000006
curl -X POST \
  $URL/ue-routes/blue/members/imsi-208930000000007
curl -X POST \
  $URL/ue-routes/blue/members/imsi-208930000000008
curl -X POST \
  $URL/ue-routes/blue/members/imsi-208930000000009
curl -X POST \
  $URL/ue-routes/blue/members/imsi-208930000000010


curl -X POST \
  $URL/ue-routes/red/topology \
  -H 'content-type: application/json' \
  -d '{
  "topology": [
    {"A": "gNB", "B": "upf-r1-sample"},
    {"A": "upf-r1-sample", "B": "upf-t1-sample"},
    {"A": "upf-t1-sample", "B": "upf-c1-sample"}
  ],
  "specificPath": [
    {"dest": "60.61.0.0/16", "path": ["upf-r1-sample", "upf-t1-sample", "upf-c1-sample"]}
  ]
}'

curl -X POST \
  $URL/ue-routes/blue/topology \
  -H 'content-type: application/json' \
  -d '{
  "topology": [
    {"A": "gNB", "B": "upf-r1-sample"},
    {"A": "upf-r1-sample", "B": "upf-t2-sample"},
    {"A": "upf-t2-sample", "B": "upf-c4-sample"}
  ],
  "specificPath": [
    {"dest": "60.64.0.0/16", "path": ["upf-r1-sample", "upf-t2-sample", "upf-c4-sample"]}
  ]
}'


curl -X GET \
  $URL/ue-routes \
  -H 'content-type: application/json'
