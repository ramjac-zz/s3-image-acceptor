package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

type ImageData struct {
	Image      []byte `json:"image"`
	Name       string `json:"name,omitempty"`
	Expiration int    `json:"expiration,omitempty"`
}

// Handler is your Lambda function handler
// It uses Amazon API Gateway request/responses provided by the aws-lambda-go/events package,
// However you could use other event sources (S3, Kinesis etc), or JSON-decoded primitive types such as 'string'.
func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	// stdout and stderr are sent to AWS CloudWatch Logs
	log.Printf("Processing Lambda request %s\n", request.RequestContext.RequestID)

	// Do some basic validation
	if request.HTTPMethod != http.MethodPost {
		return events.APIGatewayProxyResponse{StatusCode: http.StatusMethodNotAllowed},
			errors.New("only accept POST requests")
	}
	if len(request.Body) < 1 {
		return events.APIGatewayProxyResponse{StatusCode: http.StatusBadRequest},
			errors.New("no content was provided in the HTTP body")
	}
	if request.Headers["Content-Type"] != "application/json" {
		return events.APIGatewayProxyResponse{StatusCode: http.StatusBadRequest},
			errors.New("only accepts application/json")
	}

	fmt.Println(request.Body)

	// get the image data from the POST
	var image *ImageData

	err := json.Unmarshal([]byte(request.Body), &image)

	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: http.StatusBadRequest},
			errors.New("json could not be parsed")
	}

	// put the image in S3
	// The session the S3 Uploader will use
	sess := session.Must(session.NewSession())

	// Create an uploader with the session and default options
	uploader := s3manager.NewUploader(sess)

	// Upload the file to S3.
	result, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String("gritty-image-uploads"),
		Key:    aws.String(image.Name),
		Body:   bytes.NewReader(image.Image),
	})
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: http.StatusInternalServerError},
			errors.New("could not upload the file to s3")
	}
	fmt.Printf("file uploaded to, %s\n", aws.StringValue(&result.Location))

	// return the image's URL
	return events.APIGatewayProxyResponse{
		Body:       image.Name,
		StatusCode: 200,
	}, nil
}

func main() {
	lambda.Start(Handler)
}
