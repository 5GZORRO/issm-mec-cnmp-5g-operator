#!/bin/sh

count=`ps aux | grep "[n]ssf" | wc -l`
if [ $count -eq 0 ];then
  echo "Starting NSSF..."
  /free5gc/nssf/nssf -nssfcfg > /dev/null 2> /dev/null &
  echo "...done"
fi
