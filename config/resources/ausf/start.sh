#!/bin/sh

count=`ps aux | grep "[a]usf" | wc -l`
if [ $count -eq 0 ];then
  echo "Starting AUSF..."
  /free5gc/ausf/ausf -ausfcfg /free5gc/config/ausf.yaml > /dev/null 2> /dev/null &
  echo "...done"
fi

