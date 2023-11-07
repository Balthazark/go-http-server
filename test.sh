#!/bin/bash
for i in {1..20}
do
    curl http://localhost:8080/html
    echo $i
done