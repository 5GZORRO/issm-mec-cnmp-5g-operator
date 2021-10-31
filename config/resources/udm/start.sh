#!/bin/sh

count=`ps aux | grep "[u]dm" | wc -l`
if [ $count -eq 0 ];then
  echo "Starting UDM..."
  /free5gc/udm/udm -udmcfg /free5gc/config/udm.yaml > /dev/null 2> /dev/null &
  echo "...done"
fi
