package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-lambda-go/lambda"
	"log"

	"github.com/go-resty/resty/v2"
	"github.com/sirupsen/logrus"
)

type Message struct {
	Name string `json:"name"`
}

type Result struct {
	Name    string `json:"name"`
	Details string `json:"details"`
}

func HandleRequest(ctx context.Context, msg Message) (Result, error) {
	// test possible urls for a service
	// depending on network that lambda is running on different urls could work
	details := "attempt 1:  \n"

	client := NewRestyClient("http://analyzer.test-localstack.svc.cluster.local:8081", true)
	details += GetRequestLog(client)

	details += "attempt 2:  \n"
	client.SetHostURL("http://analyzer:8081")
	details += GetRequestLog(client)

	details += "attempt 3:  \n"
	client.SetHostURL("http://localhost:8081")
	details += GetRequestLog(client)

	details += "attempt 4:  \n"
	client.SetHostURL("http://localhost")
	details += GetRequestLog(client)

	log.Printf("lambda result: %v", details)
	result := Result{Details: details, Name: msg.Name}
	return result, nil
}

func NewRestyClient(s string, b bool) *resty.Client {
	client := resty.New().
		SetHostURL(s).
		SetDebug(b).
		SetLogger(logrus.StandardLogger())
	return client
}

func GetRequestLog(c *resty.Client) string {
	req := c.R()
	res, err := req.Get("/analyze/msg")
	return fmt.Sprintf("\n err: %v, result: %v \n", err, res)
}

func main() {
	lambda.Start(HandleRequest)
}
