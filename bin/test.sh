#!/bin/bash

if [ "#" -ne 2 ]; then
	echo "Usage: test.sh NumNodes PortID"
	exit 1
fi

for (( c=1; c<=$1; c++ ))
do
	echo ./gossipserver --port $2 --nodeid $c --numnodes $1
	./gossipserver -port $2 -nodeid $c -numnodes $1 &
done
