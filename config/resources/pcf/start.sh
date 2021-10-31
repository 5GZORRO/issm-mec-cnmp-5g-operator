#!/bin/sh

count=`ps aux | grep "[p]cf" | wc -l`
if [ $count -eq 0 ];then
  echo "Starting PCF..."
  /free5gc/pcf/pcf -pcfcfg /free5gc/config/pcf.yaml > /dev/null 2> /dev/null &
  echo "...done"
fi
