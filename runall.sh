#!/bin/bash
echo "------- Downloading raw data -------"
go run getrawdata.go && \
    cd exp && \
    for d in 20*; do
        echo "------- Running experiment $d -------"
        cd $d && bash run.sh && cd ..;
    done;
