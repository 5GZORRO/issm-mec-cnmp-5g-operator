#!/bin/sh

count=`ps aux | grep "[w]ebconsole" | wc -l`
if [ $count -gt 0 ];then
  pkill -f webconsole
fi

exit 0