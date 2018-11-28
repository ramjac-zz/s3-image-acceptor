package main_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	main "github.com/ramjac/s3-image-acceptor"
	"github.com/stretchr/testify/assert"
)

// base64 encoding of a single pixel .png image
var imageToTest = "iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mNk+P+/HgAFhAJ/wlseKgAAAABJRU5ErkJggg=="

func TestHandler(t *testing.T) {
	tests := []struct {
		request events.APIGatewayProxyRequest
		expect  string
		err     error
	}{
		{
			// Test that the handler responds with the correct response
			// when a valid POST is made
			request: events.APIGatewayProxyRequest{Body: `{
				"image":"iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mNk+P+/HgAFhAJ/wlseKgAAAABJRU5ErkJggg==",
				"name":"test_image.png",
				"expiration": 1
			}`,
				HTTPMethod: "POST",
				Headers:    map[string]string{"Content-Type": "application/json"},
			},
			expect: "test_image.png",
			err:    nil,
		},
		{
			// invalid content type
			request: events.APIGatewayProxyRequest{Body: `{
				"image":"iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mNk+P+/HgAFhAJ/wlseKgAAAABJRU5ErkJggg==",
				"name":"test_image.png",
				"expiration": 1
			}`,
				HTTPMethod: "POST",
				//Headers:    map[string]string{"Content-Type": "application/json"},
			},
			expect: "",
			err:    errors.New("only accepts application/json"),
		},
		{
			// no posted content
			request: events.APIGatewayProxyRequest{Body: "",
				HTTPMethod: "POST",
				Headers:    map[string]string{"Content-Type": "application/json"},
			},
			expect: "",
			err:    errors.New("no content was provided in the HTTP body"),
		},
	}

	for _, test := range tests {
		response, err := main.Handler(test.request)
		assert.IsType(t, test.err, err)
		assert.Equal(t, test.expect, response.Body)
		fmt.Println(err)
	}
}
