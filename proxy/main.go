package main

import (
	"bytes"
	"context"
	"encoding/base64"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

// Response is of type APIGatewayProxyResponse since we're leveraging the
// AWS Lambda Proxy Request functionality (default behavior)
//
// https://serverless.com/framework/docs/providers/aws/events/apigateway/#lambda-proxy-integration
type Response events.APIGatewayProxyResponse

// Handler is our lambda handler invoked by the `lambda.Start` function call
func Handler(ctx context.Context) (Response, error) {
	res, err := http.Get("http://arxiv.org/pdf/nucl-th/9911047.pdf")
	if err != nil {
		log.Panic(err)
	}
	defer res.Body.Close()

	ct := res.Header.Get("Content-Type")
	log.Printf("content type is " + ct)

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Panic(err)
	}

	var buf bytes.Buffer
	encoder := base64.NewEncoder(base64.StdEncoding, &buf)
	encoder.Write(data)
	encoder.Close()

	resp := Response{
		StatusCode:      200,
		IsBase64Encoded: true,
		Body:            buf.String(),
		Headers: map[string]string{
			"Content-Type": "application/pdf",
		},
	}

	return resp, nil
}

func main() {
	lambda.Start(Handler)
}
