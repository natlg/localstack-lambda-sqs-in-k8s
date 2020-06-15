package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/endpoints"
	"github.com/aws/aws-sdk-go/aws/session"
)

const (
	workerEndpoint = "http://worker:4566"
)

var details string
var qResultURL = fmt.Sprintf("%s/queue/result_sqs", workerEndpoint)
var qS3URL = fmt.Sprintf("%s/queue/s3_sqs", workerEndpoint)
var directMsgReceived int
var sqsMsgReceived int
var s3MsgReceived int

func main() {
	sess, err := session.NewSession(&aws.Config{
		Credentials:      credentials.NewStaticCredentials("some", "secret", ""),
		Region:           aws.String(endpoints.UsEast1RegionID),
		Endpoint:         aws.String(workerEndpoint),
		S3ForcePathStyle: aws.Bool(true)},
	)
	if err != nil {
		panic(err)
	}
	qc := sqs.New(sess)
	qp := &QueuePoller{*qc}
	go qp.pollQueue(qResultURL)
	go qp.pollQueue(qS3URL)

	r := gin.Default()

	r.GET("/analyze/msg/:details", func(c *gin.Context) {
		directMsgReceived++
		details += c.Params.ByName("txt")
		log.Printf("directMsgReceived: %v, sqsMsgReceived %v s3MsgReceived %v", directMsgReceived, sqsMsgReceived, s3MsgReceived)
		c.JSON(http.StatusOK, gin.H{"directMsgReceived": directMsgReceived, "sqsMsgReceived": sqsMsgReceived, "s3MsgReceived": s3MsgReceived, "details": details})
	})

	r.GET("/analyze/stats", func(c *gin.Context) {
		log.Printf("stats, directMsgReceived: %v, sqsMsgReceived %v s3MsgReceived %v", directMsgReceived, sqsMsgReceived, s3MsgReceived)

		filesDetails, err := listBucketsAndFiles(sess)
		if err != nil {
			log.Printf("listBucketsAndFiles err %v", err)
		}

		c.JSON(http.StatusOK, gin.H{"directMsgReceived": directMsgReceived, "sqsMsgReceived": sqsMsgReceived, "s3MsgReceived": s3MsgReceived, "details": details, "files": filesDetails})
	})
	if err := r.Run(":8081"); err != nil {
		return
	}
}

func listBucketsAndFiles(sess *session.Session) (string, error) {

	svc := s3.New(sess)
	r, err := svc.ListBuckets(nil)
	if err != nil {
		log.Printf("ListBuckets err %v", err)
		return "ListBuckets failed", err
	}
	res := " Buckets: "

	for _, b := range r.Buckets {
		res += aws.StringValue(b.Name)

		resp, err := svc.ListObjectsV2(&s3.ListObjectsV2Input{Bucket: b.Name})
		if err != nil {
			log.Printf("ListObjectsV2 err %v", err)
			res += " ListObjects error"
			return res, err
		}
		log.Printf(" total files %v ", len(resp.Contents))
		res += " Files: "
		for _, item := range resp.Contents {
			log.Printf(" file: %v, size: %v StorageClass %v", *item.Key, *item.Size, *item.StorageClass)
			res += *item.Key
		}
	}
	return res, nil
}

type QueuePoller struct {
	client sqs.SQS
}

func (mp *QueuePoller) pollQueue(qURL string) {
	for {
		log.Printf("polling %v", qURL)

		msgInpput := &sqs.ReceiveMessageInput{
			MaxNumberOfMessages: aws.Int64(10),
			QueueUrl:            &qURL,
			VisibilityTimeout:   aws.Int64(10),
			WaitTimeSeconds:     aws.Int64(10),
		}
		msgOutput, err := mp.client.ReceiveMessage(msgInpput)
		if err != nil {
			log.Printf("ReceiveMessage err %v", err)
		}

		if len(msgOutput.Messages) < 1 {
			log.Printf("no messages in %v", qURL)
			time.Sleep(20 * time.Second)
		}

		for _, message := range msgOutput.Messages {
			mp.processMessage(message, qURL)
		}
	}
}

func (mp *QueuePoller) processMessage(message *sqs.Message, qURL string) {
	log.Printf("=== processMessage q %v, msg %v ", qURL, *message.Body)
	if qURL == qResultURL {
		sqsMsgReceived++
	} else if qURL == qS3URL {
		s3MsgReceived++
	} else {
		log.Printf("unknown q")
	}
	log.Printf("msg err %v", *message.Body)

	dmr := &sqs.DeleteMessageInput{
		QueueUrl:      &qURL,
		ReceiptHandle: message.ReceiptHandle,
	}
	_, err := mp.client.DeleteMessage(dmr)
	if err != nil {
		log.Printf("DeleteMessage err %v", err)
	}
}
