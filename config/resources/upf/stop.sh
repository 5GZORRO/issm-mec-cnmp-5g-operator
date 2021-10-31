#!/bin/sh

count=`ps aux | grep "[f]ree5gc-upfd" | wc -l`
if [ $count -gt 0 ];then
  pkill -f upf
fi

exit 0