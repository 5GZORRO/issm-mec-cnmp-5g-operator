#!/bin/sh

count=`ps aux | grep "[u]dm" | wc -l`
if [ $count -gt 0 ];then
  pkill -f udm
fi

exit 0