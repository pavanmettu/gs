#!/bin/bash

for (( c=1; c<=$1; c++ ))
do
	echo ./gossipserver --port $2 --nodeid $c --numnodes $1
	./gossipserver -port $2 -nodeid $c -numnodes $1 &
done
