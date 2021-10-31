#!/bin/sh

count=`ps aux | grep "[u]dr" | wc -l`
if [ $count -gt 0 ];then
  pkill -f udr
fi

exit 0