#!/bin/bash
# 
# Serve dashboard locally (NB: first run "kubectl proxy --port=8080")
#

# Docker options
DOCKER_RUN_OPTS=${DOCKER_RUN_OPTS:-}
DASHBOARD_IMAGE_NAME="kubernetes-dashboard-build-image"
DEFAULT_COMMAND=${DEFAULT_COMMAND:-"node_modules/.bin/gulp"}
DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

# Serve dashboard locally
sudo docker run \
	-it \
	--rm \
	--net=host \
	-v /var/run/docker.sock:/var/run/docker.sock \
	-v ${DIR}/src/app/:/dashboard/src/app/ \
	${DOCKER_RUN_OPTS} \
	${DASHBOARD_IMAGE_NAME} \
	${DEFAULT_COMMAND} serve
