#!/bin/sh

count=`ps aux | grep "[w]ebconsole" | wc -l`
if [ $count -eq 0 ];then
  echo "Starting Webconsole..."
  /free5gc/webconsole/webconsole -webuicfg /free5gc/config/webui.yaml > /dev/null 2> /dev/null &
  echo "...done"
fi
