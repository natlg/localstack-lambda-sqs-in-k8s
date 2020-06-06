#!/usr/bin/env bash

set -x
set -e

scriptdir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
eval $(make --no-print-directory -C ${scriptdir} build.vars)

echo ${BUILD_REGISTRY}
echo ${DOCKER_REGISTRY}

WORKER_IMAGE="${BUILD_REGISTRY}/worker-amd64:latest"
KIND_WORKER_IMAGE="natlg/worker:latest"

PUBLISHER_IMAGE="${BUILD_REGISTRY}/publisher-amd64:latest"
KIND_PUBLISHER_IMAGE="natlg/publisher:latest"

ANALYZER_IMAGE="${BUILD_REGISTRY}/analyzer-amd64:latest"
KIND_ANALYZER_IMAGE="natlg/analyzer:latest"

PROVISION_IMAGE="${BUILD_REGISTRY}/provision-localstack-amd64:latest"
KIND_PROVISION_IMAGE="natlg/provision-localstack:latest"

function configureKind() {
kind create cluster --name test-localstack --config=cluster/kind.yaml

kubectl config use-context kind-test-localstack
kubectl create namespace test-localstack

kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/nginx-0.30.0/deploy/static/mandatory.yaml
kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/nginx-0.30.0/deploy/static/provider/baremetal/service-nodeport.yaml
kubectl patch deployments -n ingress-nginx nginx-ingress-controller -p '{"spec":{"template":{"spec":{"containers":[{"name":"nginx-ingress-controller","ports":[{"containerPort":80,"hostPort":80},{"containerPort":443,"hostPort":443}]}],"nodeSelector":{"ingress-ready":"true"},"tolerations":[{"key":"node-role.kubernetes.io/master","operator":"Equal","effect":"NoSchedule"}]}}}}'
}

function loadAllImages() {
  loadImage ${WORKER_IMAGE} ${KIND_WORKER_IMAGE}
  loadImage ${PUBLISHER_IMAGE} ${KIND_PUBLISHER_IMAGE}
  loadImage ${ANALYZER_IMAGE} ${KIND_ANALYZER_IMAGE}
  loadImage ${PROVISION_IMAGE} ${KIND_PROVISION_IMAGE}
}

function loadImage() {
  local build_image=$1
  local final_image=$2
  docker tag "${build_image}" "${final_image}"
  echo "Tagged image: ${build_image} - ${final_image}"
  if $ENABLE_KIND; then
    echo "loading $final_image in Kind"
    kind load docker-image $final_image --name test-localstack
  fi
}

function installAll() {
kubectl apply -f cluster/worker.yaml -n test-localstack
kubectl apply -f cluster/publisher.yaml -n test-localstack
kubectl apply -f cluster/analyzer.yaml -n test-localstack
kubectl apply -f cluster/ingress.yaml -n test-localstack

}

function deleteAll() {
kubectl delete -f  cluster/worker.yaml  -n test-localstack
kubectl delete -f  cluster/publisher.yaml  -n test-localstack
kubectl delete -f  cluster/analyzer.yaml  -n test-localstack
kubectl delete -f  cluster/ingress.yaml  -n test-localstack

}

case "${1:-}" in
  install)
    loadAllImages
    echo "loaded"

    installAll
    echo "installed"
    ;;

  delete)
    echo "removing "
    deleteAll
    echo "removed"
    ;;

  restart)
    echo "removing "
    deleteAll

    sleep 30

    echo "loading"
    loadAllImages

    echo "installing"
    installAll
    ;;

  init)
    configureKind
    echo "configured"
    ;;

  clean)
    kind delete cluster --name test-localstack
    echo "deleted kind cluster"
    ;;

  *)
    echo "usage:" >&2
    echo "  $0 install - load all images into cluster and deploy" >&2
    echo "  $0 delete all deployments" >&2
    echo "  $0 restart - reinstall" >&2
    echo "  $0 init - create and configure kind cluster and ingress" >&2
    echo "  $0 clean - delete kind cluster" >&2
esac
