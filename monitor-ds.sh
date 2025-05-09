#!/bin/bash

# Replace these with the actual DS IPs from dig
DS_IPS=("16.162.197.255" "18.163.17.131" "18.166.142.150")
PORT=8800

echo "Monitoring active connections to datastreamer IPs on port $PORT..."
echo "Press Ctrl+C to stop."

while true; do
    echo "------ $(date) ------"
    
    for ip in "${DS_IPS[@]}"; do
        count=$(netstat -atn | grep "$ip:$PORT" | wc -l)
        echo "Connections to $ip:$PORT => $count"
    done
    
    sleep 5
done
