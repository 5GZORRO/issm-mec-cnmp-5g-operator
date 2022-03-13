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
    {"A": "gNB1", "B": "UPF-R1"},
    {"A": "UPF-R1", "B": "UPF-T1"},
    {"A": "UPF-T1", "B": "UPF-C1"}
  ],
  "specificPath": [
    {"dest": "60.61.0.0/16", "path": ["UPF-R1", "UPF-T1", "UPF-C1"]}
  ]
}'

curl -X POST \
  $URL/ue-routes/blue/topology \
  -H 'content-type: application/json' \
  -d '{
  "topology": [
    {"A": "gNB1", "B": "UPF-R1"},
    {"A": "UPF-R1", "B": "UPF-T2"},
    {"A": "UPF-T2", "B": "UPF-C4"}
  ],
  "specificPath": [
    {"dest": "60.64.0.0/16", "path": ["UPF-R1", "UPF-T2", "UPF-C4"]}
  ]
}'


curl -X GET \
  $URL/ue-routes \
  -H 'content-type: application/json'
