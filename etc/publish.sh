#!/bin/sh

orders=data/*.json

for f in $orders
do
    go run pub.go orders "$(cat $f)"
done
