#!/bin/bash
set -x
echo "Start provisioning"

until awslocal sqs list-queues; do
   echo "SQS is unavailable, wait"
   sleep 10
done

# seems sqs dead letter queue not implemented yet - keep an eye on https://github.com/localstack/localstack/issues/2107
awslocal sqs create-queue --queue-name worker-dead-queue
awslocal sqs create-queue --queue-name worker_sqs  --attributes '{"RedrivePolicy": "{\"deadLetterTargetArn\":\"arn:aws:sqs:us-east-1:000000000000:worker-dead-queue\",\"maxReceiveCount\":5}"}'
awslocal sqs create-queue --queue-name result_sqs
awslocal sqs create-queue --queue-name s3_sqs

awslocal lambda create-function --function-name worker  --runtime go1.x --role arn:aws:iam::000000000000:role/lambda-worker-executor --handler worker --zip-file fileb://worker.zip
awslocal lambda create-event-source-mapping --event-source-arn arn:aws:sqs:us-east-1:000000000000:worker_sqs --function-name worker

echo "created queues "

echo "configuring s3" &&
if awslocal s3api head-bucket --bucket "test" 2>/dev/null;
then
  echo "Bucket already exists"
else
  awslocal s3api create-bucket --bucket "test"
fi

awslocal s3api put-bucket-notification \
--bucket test \
--notification-configuration '{
  "QueueConfiguration":
    {
      "Id": "s3-upload-notify-sqs",
      "Queue": "arn:aws:sqs:us-east-1:000000000000:s3_sqs",
      "Event": "s3:ObjectCreated:*"
    }
}'

# new format is not implemented yet in localstack
# when it's ready, use something like
#aws s3api put-bucket-notification-configuration \
#--bucket test \
#--notification-configuration '{
#  "QueueConfigurations": [
#    {
#      "Id": "s3-upload-notify-sqs",
#      "QueueArn": "arn:aws:sqs:us-east-1:000000000000:s3_sqs",
#      "Events": [
#        "s3:ObjectCreated:*"
#      ]
#    }
#  ]
#}'

echo "configured s3"