package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/endpoints"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/go-resty/resty/v2"
	"github.com/sirupsen/logrus"
	"log"
	"strings"
)

const (
	localstackQueueSvsEndpoint = "http://localhost/queue"
	localstackS3SvsEndpoint    = "http://localhost"
	qName                      = "result_sqs"
)

var qURL = fmt.Sprintf("http://localhost/queue/%s", qName)

type Result struct {
	Name    string `json:"name"`
	Details string `json:"details"`
}

func HandleRequest(ctx context.Context, sqsEvent events.SQSEvent) (Result, error) {
	details := "===== details "
	for _, message := range sqsEvent.Records {
		log.Printf("The message %v for event source %v = %v \n", message.MessageId, message.EventSource, message.Body)
		details += fmt.Sprintf(" msg %v", message.Body)
	}
	// send message to result_sqs queue that analyzer polls from
	sess, err := session.NewSession(&aws.Config{
		Credentials:      credentials.NewStaticCredentials("some", "secret", ""),
		Region:           aws.String(endpoints.UsEast1RegionID),
		Endpoint:         aws.String(localstackQueueSvsEndpoint),
		S3ForcePathStyle: aws.Bool(true)},
	)
	if err != nil {
		log.Printf("NewSession Error %v ", err)
	}
	svc := sqs.New(sess)

	_, err = svc.SendMessage(&sqs.SendMessageInput{
		DelaySeconds:      aws.Int64(10),
		MessageAttributes: map[string]*sqs.MessageAttributeValue{},
		MessageBody:       aws.String("Some Information"),
		QueueUrl:          &qURL,
	})

	if err != nil {
		details += " failed to send msg "
		log.Printf("SendMessage Error %v ", err)
	}

	// set different endpoint for s3, need to change only for localstack
	sess.Config.Endpoint = aws.String(localstackS3SvsEndpoint)
	// put file to s3 bucket. it should trigger sending message to another queue s3_sqs that analyzer polls from
	svcS3 := s3.New(sess)
	r, err := svcS3.ListBuckets(nil)
	if err != nil {
		details += " list b err"
	}

	details += " Buckets: "

	for _, b := range r.Buckets {
		details += aws.StringValue(b.Name)
	}

	uploader := s3manager.NewUploader(sess)
	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String("test"),
		Key:    aws.String("testresultfile"),
		Body:   strings.NewReader("success!!"),
	})
	if err != nil {
		details += " upload error"
	} else {
		details += " uploaded!"
	}

	client := NewRestyClient("http://localhost", true)
	analizarRquestResult := GetRequestLog(client, details)

	log.Printf("analizarRquestResult: %v", analizarRquestResult)
	result := Result{Details: details, Name: details}
	return result, nil
}

func NewRestyClient(s string, b bool) *resty.Client {
	client := resty.New().
		SetHostURL(s).
		SetDebug(b).
		SetLogger(logrus.StandardLogger())
	return client
}

func GetRequestLog(c *resty.Client, msg string) string {
	req := c.R()
	url := fmt.Sprintf("/analyze/msg/%s", msg)
	res, err := req.Get(url)
	return fmt.Sprintf("\n err: %v, result: %v \n", err, res)
}

func main() {
	lambda.Start(HandleRequest)
}
