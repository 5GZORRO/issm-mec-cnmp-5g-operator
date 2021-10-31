#!/bin/sh

count=`ps aux | grep "[a]mf" | wc -l`
echo "AMF count: $count"
if [ $count -eq 0 ]; then
  echo "Starting AMF..."
  /free5gc/amf/amf -amfcfg /free5gc/config/amf.yaml > /dev/null 2> /dev/null &
  echo "...done"
fi
