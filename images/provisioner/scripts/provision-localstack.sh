#!/bin/bash
set -x
echo "Start provisioning"

until awslocal sqs list-queues; do
   echo "SQS is unavailable, wait"
   sleep 10
done

awslocal sqs create-queue --queue-name worker_sqs
awslocal lambda create-function --function-name worker  --runtime go1.x --role arn:aws:iam::000000000000:role/lambda-worker-executor --handler worker --zip-file fileb://worker.zip
awslocal lambda create-event-source-mapping --event-source-arn arn:aws:sqs:us-east-1:000000000000:worker_sqs --function-name worker

echo "done "