#!/bin/sh
cd account/service && go test -v -cover . && cd .. && cd ..
cd jwt/service && SECRETKEY=bwbwchen go test -v -cover . && cd .. && cd ..
cd upload/service && go test -v -cover . && cd .. && cd ..