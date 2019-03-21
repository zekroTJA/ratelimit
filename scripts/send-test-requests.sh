#!/bin/bash

for i in {1..4}; do
    printf "\n\n///////////////////////////\n"
    printf "// REQUEST #${i}\n\n"

    curl -v http://localhost:8080/api/test

    sleep 1
done