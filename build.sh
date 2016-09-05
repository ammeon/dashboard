#!/bin/bash
#
# Builds a custom K8S dashboard container and pushes it to the Docker hub registry
#

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

# Build K8S dashbaord container 
sudo ${DIR}/build/run-gulp-in-docker.sh docker-image:canary

# Tag the dashboard container
sudo docker tag \
    gcr.io/google_containers/kubernetes-dashboard-amd64:canary \
    ammeon/kubernetes-helm-dashboard-amd64:latest

# Push container to Docker hub
sudo docker push ammeon/kubernetes-helm-dashboard-amd64:latest
