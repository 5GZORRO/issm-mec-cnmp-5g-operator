#!/bin/sh

count=`ps aux | grep "[u]dr" | wc -l`
if [ $count -eq 0 ];then
  echo "Starting UDR..."
  /free5gc/udr/udr -udrcfg /free5gc/config/udr.yaml > /dev/null 2> /dev/null &
  echo "...done"
fi
