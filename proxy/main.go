package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func setCORSHeader(origin string) string {
	var allowedOrigins = [3]string{"https://scinapse.io", "https://dev.scinapse.io", "http://localhost:3000"}

	allowOrigin := "*"
	for i := 0; i < len(allowedOrigins); i++ {
		if allowedOrigins[i] == origin {
			allowOrigin = origin
		}
	}

	return allowOrigin
}

// Response is of type APIGatewayProxyResponse since we're leveraging the
// AWS Lambda Proxy Request functionality (default behavior)
//
// https://serverless.com/framework/docs/providers/aws/events/apigateway/#lambda-proxy-integration
type Response events.APIGatewayProxyResponse

// Handler is our lambda handler invoked by the `lambda.Start` function call
func Handler(req events.APIGatewayProxyRequest) (Response, error) {
	pdfURL := req.QueryStringParameters["pdf_url"]
	title := req.QueryStringParameters["title"]
	forceDownload := req.QueryStringParameters["download"]
	corsOrigin := setCORSHeader(req.Headers["origin"])

	if corsOrigin == "*" {
		resp := Response{
			StatusCode:      412,
			IsBase64Encoded: false,
			Body:            "Precondition Failed",
			Headers: map[string]string{
				"Access-Control-Allow-Origin": corsOrigin,
				"Content-Type":                "text/html",
			},
		}

		return resp, nil
	}

	var resType string
	if forceDownload != "" {
		resType = "inline"
	} else {
		resType = "attachment"
	}

	if len(pdfURL) == 0 {
		log.Fatal("Not Valid PDF URL")
	}

	res, err := http.Get(pdfURL)
	if err != nil {
		log.Panic(err)
	}
	defer res.Body.Close()

	ct := res.Header.Get("Content-Type")

	if ct != "application/pdf" {
		log.Panic(err)
	}
	log.Printf("content type is " + ct)

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Panic(err)
	}

	var buf bytes.Buffer
	encoder := base64.NewEncoder(base64.StdEncoding, &buf)
	encoder.Write(data)
	encoder.Close()

	cd := fmt.Sprintf("%s; filename=\"%s\"", resType, title)

	resp := Response{
		StatusCode:      200,
		IsBase64Encoded: true,
		Body:            buf.String(),
		Headers: map[string]string{
			"Content-Type":                "application/pdf",
			"Access-Control-Allow-Origin": corsOrigin,
			"Content-Disposition":         cd,
		},
	}

	return resp, nil
}

func main() {
	lambda.Start(Handler)
}
