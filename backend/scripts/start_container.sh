#!/bin/bash

docker pull 010629071893.dkr.ecr.us-east-1.amazonaws.com/test:latest
docker run --name server -p 80:80 -p 443:443 -v /certs/:/certs/ --rm -d 010629071893.dkr.ecr.us-east-1.amazonaws.com/test:latest 
