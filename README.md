### localstack-lambda-sqs-in-k8s
Example of working with SQS and AWS Lambda in Go running in [Kind](https://kind.sigs.k8s.io/docs/user/quick-start/) k8s cluster. 

Uses [LocalStack](https://github.com/localstack/localstack) for emulating AWS cloud stack, 
[awslocal](https://github.com/localstack/awscli-local) as a wrapper for aws cli
 and [aws-sdk-go](github.com/aws/aws-sdk-go) for calling AWS from service in Go.

#### Build:
```
make
```

#### Create Kind cluster and configure Ingress:
```
./setup.sh init
```
Each Lambda function will run in a separate Docker container and docker socket needs to be mounted into the 
localstack container, so config for Kind has it mounted too

#### Install
```
./setup.sh install
```
It will run 2 services: `publisher` and `analyzer` and expose them locally (by Ingress). Publisher can send messages to SQS 
that invokes AWS Lambda (using event source mapping). Lambda function calls `analyzer` service that keeps tracks of 
how many times it was invoked, so that it can be checked later. Localstack is provisioned in init container for `analyzer` pod. 
It contains script for creating SQS, Lambda and event source mapping.

#### Call services:
- Check currently running pods and wait until analyzer is up:
```
watch kubectl get all -n test-localstack
```

- Open http://localhost/publish/msg in browser.<br>
This request calls `publisher` service that will send message to sqs. 
Sqs invokes Lambda and Lambda function calls `analyzer` service

- Get statistics, returns how many time analyzer was called: http://localhost/publish/stats
It can take few seconds to finish work and show updated statistics. Each Lambda call should increment this value.

Call analyzer manually:
http://localhost/analyze/msg

Get statistics directly from analyzer:
http://localhost/analyze/stats

 Will work on Ubuntu, but not likely on Mac because of network configuration. Docker container with Lambda 
 is running in `host` network that is not available on Mac. <br>
 Kind v0.8 uses separate network `kind` but setting it as a network for Lambda still doesn't work because it's 
 running as Docker outside of Docker.
 
 #### Clean up:
 
 `./setup.sh delete` - remove deployments 
 
 `./setup.sh clean` - delete kind cluster




