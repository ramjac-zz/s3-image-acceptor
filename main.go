package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"

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
		log.Println(errors.New("only accept POST requests"))
		return events.APIGatewayProxyResponse{StatusCode: http.StatusMethodNotAllowed}, nil
	}
	if len(request.Body) < 1 {
		log.Println(errors.New("no content was provided in the HTTP body"))
		return events.APIGatewayProxyResponse{StatusCode: http.StatusBadRequest}, nil
	}
	if request.Headers["Content-Type"] != "application/json" {
		log.Println(errors.New("only accepts application/json"))
		return events.APIGatewayProxyResponse{StatusCode: http.StatusBadRequest}, nil
	}

	fmt.Println(request.Body)

	// get the image data from the POST
	var image *ImageData

	err := json.Unmarshal([]byte(request.Body), &image)

	if err != nil {
		log.Println(errors.New("json could not be parsed"))
		return events.APIGatewayProxyResponse{StatusCode: http.StatusBadRequest}, nil
	}

	if len(strings.TrimSpace(image.Name)) == 0 || len(image.Image) == 0 {
		log.Println(errors.New("only accepts application/json"))
		return events.APIGatewayProxyResponse{StatusCode: http.StatusBadRequest}, nil
	}

	log.Println(image)
	log.Printf("Uploading image: %x", image)

	// put the image in S3
	// The session the S3 Uploader will use
	sess := session.Must(session.NewSession())

	// Create an uploader with the session and default options
	uploader := s3manager.NewUploader(sess)

	upload := &s3manager.UploadInput{
		Bucket: aws.String("gritty-image-uploads"),
		Key:    aws.String(image.Name),
		Body:   bytes.NewReader(image.Image),
	}

	// Upload the file to S3.
	result, err := uploader.Upload(upload)
	if err != nil {
		log.Println(errors.New("could not upload the file to s3"))
		return events.APIGatewayProxyResponse{StatusCode: http.StatusInternalServerError}, nil
	}
	log.Printf("file uploaded to, %s\n", aws.StringValue(&result.Location))

	// return the image's URL
	return events.APIGatewayProxyResponse{
		Body:       `{"url":"` + result.Location + "\"}",
		Headers:    map[string]string{"Content-Type": "application/json"},
		StatusCode: 200,
	}, nil
}

func main() {
	lambda.Start(Handler)
}
