package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/endpoints"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
	"github.com/sirupsen/logrus"
	"log"
	"net/http"
)

const (
	workerEndpoint   = "http://worker:4566"
	analyzerEndpoint = "http://analyzer:8081"
)

func main() {
	r := gin.Default()
	r.GET("/publish/msg", func(c *gin.Context) {
		// send message to worker_sqs queue that will trigger lambda
		sess, err := session.NewSession(&aws.Config{
			Credentials: credentials.NewStaticCredentials("some", "secret", ""),
			Region:      aws.String(endpoints.UsEast1RegionID),
			Endpoint:    aws.String(workerEndpoint)},
		)
		if err != nil {
			log.Printf("NewSession Error %v ", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		svc := sqs.New(sess)

		qURL := fmt.Sprintf("%s/queue/worker_sqs", workerEndpoint)
		mresult, err := svc.SendMessage(&sqs.SendMessageInput{
			DelaySeconds:      aws.Int64(10),
			MessageAttributes: map[string]*sqs.MessageAttributeValue{},
			MessageBody:       aws.String("Some Information"),
			QueueUrl:          &qURL,
		})

		if err != nil {
			log.Printf("SendMessage Error %v ", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		log.Printf("MessageId %v", *mresult.MessageId)
		c.JSON(http.StatusOK, gin.H{"result": mresult.GoString()})
	})

	r.GET("/publish/stats", func(c *gin.Context) {
		client := resty.New().
			SetHostURL(analyzerEndpoint).
			SetDebug(true).
			SetLogger(logrus.StandardLogger())

		req := client.R()
		res, err := req.Get("/analyze/stats")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"StatsStatusCode": res.StatusCode(), "statRes": res.String()})
	})

	if err := r.Run(":8085"); err != nil {
		return
	}
}
