#!/bin/sh

count=`ps aux | grep "[f]ree5gc-upfd" | wc -l`
if [ $count -eq 0 ];then
  echo "Starting UPF..."
  /free5gc/free5gc-upfd/free5gc-upfd -f /free5gc/config/upf.yaml > /dev/null 2> /dev/null &
  echo "...done"
fi
