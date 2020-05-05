#!/usr/bin/env bash

set -x
set -e

function configureKind() {
kind create cluster --name test-localstack --config=cluster/kind.yaml

kubectl config use-context kind-test-localstack
kubectl create namespace test-localstack

kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/nginx-0.30.0/deploy/static/mandatory.yaml
kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/nginx-0.30.0/deploy/static/provider/baremetal/service-nodeport.yaml
kubectl patch deployments -n ingress-nginx nginx-ingress-controller -p '{"spec":{"template":{"spec":{"containers":[{"name":"nginx-ingress-controller","ports":[{"containerPort":80,"hostPort":80},{"containerPort":443,"hostPort":443}]}],"nodeSelector":{"ingress-ready":"true"},"tolerations":[{"key":"node-role.kubernetes.io/master","operator":"Equal","effect":"NoSchedule"}]}}}}'
}

function loadAllImages() {
kind load docker-image natlg/worker:latest --name test-localstack -v 1
kind load docker-image natlg/publisher:latest --name test-localstack -v 1
kind load docker-image natlg/analyzer:latest --name test-localstack -v 1
kind load docker-image natlg/provision-localstack:latest --name test-localstack -v 1
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
