#!/bin/sh

count=`ps aux | grep "[s]mf" | wc -l`
if [ $count -eq 0 ];then
  echo "Starting SMF..."
  nohup /free5gc/smf/smf -smfcfg /free5gc/config/smf.yaml -uerouting /free5gc/config/uerouting.yaml > /dev/null 2> /dev/null &
  echo "...done"
fi

