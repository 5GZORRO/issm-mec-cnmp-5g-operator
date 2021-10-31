#!/bin/sh

count=`ps aux | grep "[n]rf" | wc -l`
if [ $count -eq 0 ];then
  echo "Starting NRF..."
  /free5gc/nrf/nrf -nrfcfg /free5gc/config/nrf.yaml > /dev/null 2> /dev/null &
  echo "...done"
fi
